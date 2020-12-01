package main

import (
    "log"
    "net"
    "fmt"
    "google.golang.org/grpc"
    "github.com/dcordova/sd_tarea2/data_service"
)


func main() {

    // Comentar y descomentar al pasar a las maquinas del lab!!!
    // Segun corresponda cambiar las ip en vez de puertos

    lis1, err := net.Listen("tcp", ":9001")
    if err != nil {
        log.Fatalf("Failed to listen on port 9001: %v", err)
    }

    fmt.Println("Escuchando en el puerto 9001...")

    //-----------------Comentar en maquinas reales---------------------//

    lis2, err := net.Listen("tcp", ":9002")
    if err != nil {
        log.Fatalf("Failed to listen on port 9002: %v", err)
    }

    fmt.Println("Escuchando en el puerto 9002...")

    //-----------------Comentar en maquinas reales---------------------//

    lis3, err := net.Listen("tcp", ":9003")
    if err != nil {
        log.Fatalf("Failed to listen on port 9003: %v", err)
    }

    fmt.Println("Escuchando en el puerto 9003...")

    //------------------------------------------------------//

    s := data_service.Server{}
    grpcServer := grpc.NewServer()
    data_service.RegisterDataServiceServer(grpcServer, &s)

    go func() {
        if err := grpcServer.Serve(lis2); err != nil {
            log.Fatalf("Failed to serve gRPC Truck server on port 9002: %v", err)
        }
    }()

    go func() {
        if err := grpcServer.Serve(lis3); err != nil {
            log.Fatalf("Failed to serve gRPC Truck server on port 9003: %v", err)
        }
    }()
    

    ////// Servicio de clientes ///////
    if err := grpcServer.Serve(lis1); err != nil {
        log.Fatalf("Failed to serve gRPC server on port 9001: %v", err)
    }
    fmt.Println("Server corriendo...")
}
