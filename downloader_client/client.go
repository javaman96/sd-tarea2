package main

import (
  "log"

  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "github.com/dcordova/sd_tarea2/data_service"
  "github.com/dcordova/sd_tarea2/name_service"
  "encoding/hex"
  "strings"
  "archive/zip"
  "path/filepath"
  //"bufio"
  //"encoding/csv"
  "fmt"
  "io"
  "os"
  //"strconv" // Conversion de strings a int y viceversa
)

func unzip(zipfile string){

    reader, err := zip.OpenReader(zipfile)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer reader.Close()

    for _, f := range reader.Reader.File {

        zipped, err := f.Open()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        defer zipped.Close()

        // get the individual file name and extract the current directory
        path := filepath.Join("./", f.Name)

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
            fmt.Println("Creating directory", path)
        } else {
            writer, err := os.OpenFile(f.Name, os.O_WRONLY|os.O_CREATE, f.Mode())

            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            defer writer.Close()

            if _, err = io.Copy(writer, zipped); err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            fmt.Println("Decompressing : ", path)
        }
    }
}

// Convierte nombre del libro a hex
func encodeString(titulo string) string {
    src := []byte(titulo)
    encodedStr := hex.EncodeToString(src)
    return encodedStr
}

//func writeLibro(chunks []bait)

// Recibe nombre del libro, retorna id del chunk
func chunk_id(name_n string) (string) {
  sub := strings.Split(name_n, "_")
  return sub[len(sub)-1]
}

func main() {  
    
  //------------------------------------------------------
  //////// Conectarse como cliente al NameService ////////
  //------------------------------------------------------
  var conn *grpc.ClientConn
  conn, err := grpc.Dial(":9009", grpc.WithInsecure())
  if err != nil {
    log.Fatalf("Could not connect: %s", err)
  }
  defer conn.Close()

  s := name_service.NewNameServiceClient(conn)

  // Hello world
  message := name_service.Message{
    Body: "Conectandose desde downloader_client!",
  }

  response, err := s.SayHello(context.Background(), &message)
  if err != nil {
    log.Fatalf("Error when calling SayHello: %s", err)
  }

  log.Printf("Response from Server: %s", response.Body)  



  //------------------------------------------------------
  ///////////////////// CICLO PRINCIPAL //////////////////
  //------------------------------------------------------

  for true {
    //// OPCIONES A ELEGIR ////
    var input string
    fmt.Printf("\n 1) Mostrar lista de libros \n 2) Solicitar libro (Ver lista primero) \n Ingrese opcion: ")
    fmt.Scanln(&input)
    option := strings.TrimSpace(input)
    //option := strconv.ParseInt(input, 10, 64)

    /// OPCION 1: ENVIAR TODOS LOS ENCARGOS AL SERVER
    if option == "1" {
        message := name_service.Message{
          Body: "solicitando lista de libros",
        }
        response, err := s.PedirNombresLibros(context.Background(), &message)
        if err != nil {
          log.Fatalf("Error al llamar a PedirNombresLibros: %s", err)
        }

        for i, libro := range response.Libros {
          fmt.Printf("%d - Nombre: %s\n", i, libro.Nombre)
        }
    }

    /// OPCION 1: PEDIR UBICACIONES DE CHUNKS DEL LIBRO
    if option == "2" {
      fmt.Printf("Ingrese nombre del libro: ")
      fmt.Scanln(&input)
      
      message := name_service.LibroInfo{
        Nombre: input,
      }
      response, err := s.PedirChunksLibro(context.Background(), &message)
      if err != nil {
        log.Fatalf("Error al llamar a PedirChunksLibro: %s", err)
      }  
      
      newFileName := "downloader_client/downloads/" + strings.Split(input,".pdf")[0]+".zip"

      // create file
      _, err = os.Create(newFileName)

      if err != nil {
          fmt.Println(err)
          os.Exit(1)
      }     
      
      // set the newFileName file to APPEND MODE!!
      // open files r and w

      file, err := os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)      
      if err != nil {
        log.Fatal(err)
        os.Exit(1)
      }    

      /// DESCARGAR CHUNKS
      for _, chunk := range response.Chunks {
        fmt.Printf("Nombre_chunk: %s - Ip_maquina: %s\n", chunk.Nombrechunk, chunk.Ipmaquina)
        
        in_hex := encodeString(input) + "_" + chunk_id(chunk.Nombrechunk)

        /// CONECTARSE A LA MAQUINA
        var conn2 *grpc.ClientConn
        conn2, err := grpc.Dial(chunk.Ipmaquina, grpc.WithInsecure())
        if err != nil {
          log.Fatalf("Could not connect to "+chunk.Ipmaquina+": %s", err)
        }
        defer conn2.Close()

        c2 := data_service.NewDataServiceClient(conn2)

        
        currentChunkFileName, err := c2.RecuperarChunks(context.Background(), &data_service.Message{Body:in_hex})
        if err != nil {
          log.Fatalf("Error when calling Server "+chunk.Ipmaquina+": %s", err)
        }   

        //----------------------------------------//
        ///////////// RECOMBINAR CHUNKS ////////////
        //----------------------------------------//        

        // DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
        // write/save buffer to disk
        //ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)

        _, err = file.Write(currentChunkFileName.Data)

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        file.Sync() //flush to disk

        // free up the buffer for next cycle
        // should not be a problem if the chunk size is small, but
        // can be resource hogging if the chunk size is huge.
        // also a good practice to clean up your own plate after eating

        currentChunkFileName = nil // reset or empty our buffer        
      } 

      // Descomprimir
      unzip(newFileName)

      // Delete remaining zip files
      chunk_dir, err := os.Open("downloader_client/downloads/")
      if err != nil {
        log.Fatalf("failed opening directory: %s", err)
      }

      chunk_list,_ := chunk_dir.Readdirnames(0) // 0 to read all files and folders

      for _, name := range chunk_list {
        
        if strings.Contains(name, ".zip") {
          fmt.Println("BORRANDO:", name)
          os.Remove(name)
        }
      }
      defer chunk_dir.Close()             
    }
  }
}
