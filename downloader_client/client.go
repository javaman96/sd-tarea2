package main

import (
  "log"

  "golang.org/x/net/context"
  "google.golang.org/grpc"
  //"github.com/dcordova/sd_tarea2/data_service"
  "github.com/dcordova/sd_tarea2/name_service"

  //"strings"
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

  c := name_service.NewNameServiceClient(conn)

  // Hello world
  message := name_service.Message{
    Body: "Conectandose desde downloader_client!",
  }

  response, err := c.SayHello(context.Background(), &message)
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
          fmt.Printf("No implementado")
        }
      }

      if option == "2" {
        fmt.Printf("No implementado")
      }
    }
}
