package main

import (
    "log"
    "net"
    "fmt"
    "google.golang.org/grpc"
    "github.com/dcordova/sd_tarea2/data_service"
)


func main() {

    lis, err := net.Listen("tcp", ":9000")
    if err != nil {
        log.Fatalf("Failed to listen on port 9000: %v", err)
    }

    fmt.Println("Escuchando en el puerto 9000...")

    s := data_service.Server{}
    grpcServer := grpc.NewServer()
    data_service.RegisterDataServiceServer(grpcServer, &s)

    ////// Servicio de clientes ///////
    if err = grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve gRPC server on port 9000: %v", err)
    }
    fmt.Println("Server corriendo...")
}