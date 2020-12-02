package main

import (
  "log"

  "golang.org/x/net/context"
  "google.golang.org/grpc"
  //"github.com/dcordova/sd_tarea2/data_service"
  "github.com/dcordova/sd_tarea2/name_service"

  "strings"
  //"bufio"
  //"encoding/csv"
  "fmt"
  //"io"
  //"os"
  //"strconv" // Conversion de strings a int y viceversa
)






func main() {

  // Conectarse como cliente al servidor //
  var conn *grpc.ClientConn
  conn, err := grpc.Dial(":9000", grpc.WithInsecure())
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

  /////// CICLO PRINCIPAL ////////
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

      if option == "2" {
        fmt.Printf("Ingrese nombre del libro: ")
        fmt.Scanln(&input)

        // Solicitar chunks
        message := name_service.LibroInfo{
          Nombre: input,
        }
        response, err := s.PedirChunksLibro(context.Background(), &message)
        if err != nil {
          log.Fatalf("Error al llamar a PedirChunksLibro: %s", err)
        }
        for _, chunk := range response.Chunks {
          fmt.Printf("Nombre_chunk: %s - Ip_maquina: %s\n", chunk.Nombrechunk, chunk.Ipmaquina)
        }
      }
    }
}
