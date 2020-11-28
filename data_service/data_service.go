package data_service

import (
    "log"
    "fmt"
    "golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) UploadChunks(ctx context.Context, message *Book) (*Book, error) {

	res := make([]byte, 1) //prueba de respuesta

    log.Printf("Received a new message body from client: %s", message.Chunks)
    fmt.Println(message.Data)
    return &Book{Chunks: message.Chunks, Data: res}, nil
}
