package main

import (
    "log"
    "net"
    "fmt"
    "google.golang.org/grpc"    
    "github.com/dcordova/sd_tarea2/data_service"
)


func main() {

    // Borrar
    var input string
    fmt.Printf("\n Ingrese numero de maquina (1,2 o 3 owo): ")
    fmt.Scanln(&input)

    // Comentar y descomentar al pasar a las maquinas del lab!!!
    // Segun corresponda cambiar las ip en vez de puertos

    s := data_service.Server{}
    grpcServer := grpc.NewServer()
    data_service.RegisterDataServiceServer(grpcServer, &s)


    lis, err := net.Listen("tcp", ":900" + input)
    if err != nil {
        log.Fatalf("Failed to listen on port 900"+input+": %v", err)
    }

    fmt.Println("Escuchando en el puerto 900"+input+"...")

    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve gRPC server on port 900"+input+": %v", err)
    }        
    
    fmt.Println("Server corriendo...")
}
