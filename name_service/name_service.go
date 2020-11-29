package name_service

import (
    "log"
    //"fmt"
    "golang.org/x/net/context"
    "os"
    "bufio"
    //"fmt"
    "strings"
    "strconv"
)

type Libro struct {
  id, nombre string
}

type Server struct {
}

//////   Esta función era del tutorial pero la dejamos    ///////
//////   para ratificar la conexion con el servidor       ///////
func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
    log.Printf("Received a new message body from client: %s", message.Body)
    return &Message{Body: "Conectado desde name_service! "}, nil
}

func (s *Server) PedirNombresLibros(ctx context.Context, message *Message) (*ListadoLibros, error) {

    retorno := []*LibroInfo{}

    logs_file, err:= os.Open("name_node/logs/logs.txt")
    if err != nil {
      log.Fatal(err)
      return nil, err
    }

    scanner := bufio.NewScanner(logs_file)
    scanner.Split(bufio.ScanLines)
    var file_lines []string
    for scanner.Scan() {
        file_lines = append(file_lines, scanner.Text())
    }

    logs_file.Close()

    for i := 0; i < len(file_lines); i++ {
      s := strings.Split(file_lines[i], " ")
      nombre_libro, cantidad_chunks := s[0],s[1]

      number, _ := strconv.Atoi(cantidad_chunks)
      i += int(number)
      retorno = append(retorno, &LibroInfo{Id:"Id", Nombre:nombre_libro})
    }

    return &ListadoLibros{Libros: retorno}, nil
}

// Pedir lista de libros
