package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	//"github.com/dcordova/sd_tarea2/data_service"
	"github.com/dcordova/sd_tarea2/name_service"
)

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	fmt.Println("Escuchando en el puerto 9000...")

	// Setear server
	s := name_service.Server{}

	grpcServer := grpc.NewServer()
	name_service.RegisterNameServiceServer(grpcServer, &s)

	////// Servicio de clientes ///////
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server on port 9000: %v", err)
	}
	fmt.Println("Server corriendo...")
}
