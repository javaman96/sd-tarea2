package main

import (
    "log"
    "net"
    "fmt"
    "google.golang.org/grpc"
    //"github.com/dcordova/sd_tarea2/data_service"
    "github.com/dcordova/sd_tarea2/name_service"
)

type libro struct {
  id, nombre string
}

var lista_libros = []libro{ libro{id:"123", nombre:"Don Quijote"}, libro{id:"124", nombre:"La Divina Comedia"}}

func main() {

    lis, err := net.Listen("tcp", ":9000")
    if err != nil {
        log.Fatalf("Failed to listen on port 9000: %v", err)
    }

    fmt.Println("Escuchando en el puerto 9000...")

    // Setear server
    s := name_service.Server{}
    s.Info_libros.append(name_service.LibroInfo{id:"123", nombre: "Quijote"})

    grpcServer := grpc.NewServer()
    name_service.RegisterNameServiceServer(grpcServer, &s)

    ////// Servicio de clientes ///////
    if err = grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve gRPC server on port 9000: %v", err)
    }
    fmt.Println("Server corriendo...")
}
