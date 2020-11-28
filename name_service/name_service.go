package name_service

import (
    "log"
    //"fmt"
    "golang.org/x/net/context"
)

type Libro struct {
  id, nombre string
}

type Server struct {
  ListadoLibros []Libro
}

//////   Esta función era del tutorial pero la dejamos    ///////
//////   para ratificar la conexion con el servidor       ///////
func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
    log.Printf("Received a new message body from client: %s", message.Body)
    return &Message{Body: "Conectado desde name_service! "}, nil
}

func (s *Server) PedirNombresLibros(ctx context.Context, message *Message) (*ListadoLibros, error) {
    return &s.Info_libros, nil
}

// Pedir lista de libros
