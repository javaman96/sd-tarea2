package data_service

import (
    "log"
    //"fmt"
    "golang.org/x/net/context"
)

type Server struct {
}

//////   Esta función era del tutorial pero la dejamos    ///////
//////   para ratificar la conexion con el servidor       ///////
func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
    log.Printf("Received a new message body from client: %s", message.Body)
    return &Message{Body: "Hello From the server! "}, nil
}
