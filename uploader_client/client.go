package main

import (
  "log"

  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "github.com/dcordova/sd_tarea2/data_service"
  "fmt"
  //"io/ioutil"
  "math"
  "os"
  "io"
  "strconv"
  "strings"
  "archive/zip"
  "encoding/hex"
)


func encodeString(titulo string) string {
    src := []byte(titulo)
    encodedStr := hex.EncodeToString(src)
    return encodedStr
}

func ZipFile(filename string) error {      

    newfile, err := os.Create(strings.Split(filename,".pdf")[0] + ".zip")
    if err != nil {
        return err
    }
    defer newfile.Close()

    zipit := zip.NewWriter(newfile)

    defer zipit.Close()

    zipfile, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer zipfile.Close()

    // get the file information
    info, err := zipfile.Stat()
    if err != nil {
        return err
    }

    header, err := zip.FileInfoHeader(info)
    if err != nil {
        return err
    }

    header.Method = zip.Deflate

    writer, err := zipit.CreateHeader(header)
    if err != nil {
        return err
    }
    _, err = io.Copy(writer, zipfile)
    return err
}


func main() {

  // Conectarse como cliente al servidor1 //
  var conn *grpc.ClientConn
  conn, err := grpc.Dial(":9001", grpc.WithInsecure())
  if err != nil {
    log.Fatalf("Could not connect to 9001: %s", err)
  }
  defer conn.Close()

  c := data_service.NewDataServiceClient(conn)

  //--------------------------------------------------//

  // Conectarse como cliente al servidor2 //
  var conn2 *grpc.ClientConn
  conn2, err = grpc.Dial(":9002", grpc.WithInsecure())
  if err != nil {
    log.Fatalf("Could not connect to 9002: %s", err)
  }
  defer conn2.Close()

  c2 := data_service.NewDataServiceClient(conn)

  //--------------------------------------------------//

  // Conectarse como cliente al servidor3 //
  var conn3 *grpc.ClientConn
  conn3, err = grpc.Dial(":9003", grpc.WithInsecure())
  if err != nil {
    log.Fatalf("Could not connect to 9003: %s", err)
  }
  defer conn3.Close()

  c3 := data_service.NewDataServiceClient(conn)

  //--------------------------------------------------//

  var input string
  fmt.Printf("\n Ingrese nombre completo del archivo (archivo.pdf): ")
  fmt.Scanln(&input)

  input = strings.TrimSpace(input)

  chunkname := encodeString(input) + "_"
  //fmt.Println(chunkname)

  err = ZipFile(input)
  if err != nil {
      fmt.Println(err)
      os.Exit(1)
  }

  // Modificar ruta de prueba
  fileToBeChunked := "./"+strings.Split(input,".pdf")[0] + ".zip"    

  file, err := os.Open(fileToBeChunked)

  if err != nil {
      fmt.Println(err)
      os.Exit(1)
  }

  defer file.Close()

  fileInfo, _ := file.Stat()

  var fileSize int64 = fileInfo.Size()

  //fmt.Println(fileSize)

  const fileChunk = 250 * (1 << 10) // 250KB

  // calculate total number of parts the file will be chunked into

  totalPartsNum := int(math.Ceil(float64(fileSize) / float64(fileChunk)))

  fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)  

  for i := 0; i < totalPartsNum; i++ {

    partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
    partBuffer := make([]byte, partSize)

    file.Read(partBuffer)
    
    IdLibroyChunk := chunkname + strconv.Itoa(i)      
                
    chunk := data_service.Book{
      Chunks: IdLibroyChunk,
      Data: partBuffer,
    }

    // Por ahora el criterio es enviar a cada server una porcion

    if i < int(math.Ceil(float64(totalPartsNum)/3)){    

      _, err := c.UploadChunks(context.Background(), &chunk)
      if err != nil {
        log.Fatalf("Error when calling Server: %s", err)
      }      
    }else if i < int(2*math.Ceil(float64(totalPartsNum)/3)){   
      
      _, err := c2.UploadChunks(context.Background(), &chunk)
      if err != nil {
        log.Fatalf("Error when calling Server: %s", err)
      }      
    }else if i < int(3*math.Ceil(float64(totalPartsNum)/3)){
      
      _, err := c3.UploadChunks(context.Background(), &chunk)
      if err != nil {
        log.Fatalf("Error when calling Server: %s", err)
      } 
    }
  }  



  // Delete remaining zip files
  _dir, err := os.Open(".")
  if err != nil {
      log.Fatalf("failed opening directory: %s", err)
  }
  defer _dir.Close()

  file_list,_ := _dir.Readdirnames(0) // 0 to read all files and folders

  for _, name := range file_list {
      if strings.Contains(name, ".zip") {
          os.Remove(name)
      }
  }  
}
