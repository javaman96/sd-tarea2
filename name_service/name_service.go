package name_service

import (
	"log"

	//"fmt"
	"bufio"
	"os"

	"golang.org/x/net/context"

	//"fmt"
	"strconv"
	"strings"
)

type Libro struct {
	id, nombre string
}

type Server struct {
	TotalLibros int
}

//////   Esta función era del tutorial pero la dejamos    ///////
//////   para ratificar la conexion con el servidor       ///////
func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received a new message body from client: %s", message.Body)
	return &Message{Body: "Conectado desde name_service! "}, nil
}

func (s *Server) PedirNombresLibros(ctx context.Context, message *Message) (*ListadoLibros, error) {

	retorno := []*LibroInfo{}

	logs_file, err := os.Open("name_node/logs/logs.txt")
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
		nombre_libro, cantidad_chunks := s[0], s[1]

		number, _ := strconv.Atoi(cantidad_chunks)
		i += int(number)
		retorno = append(retorno, &LibroInfo{Nombre: nombre_libro})
	}

	return &ListadoLibros{Libros: retorno}, nil
}

func (s *Server) PedirChunksLibro(ctx context.Context, libro_solicitado *LibroInfo) (*DistribucionChunks, error) {

	retorno := []*ChunkIp{}

	logs_file, err := os.Open("name_node/logs/logs.txt")
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
		nombre_libro, cantidad_chunks := s[0], s[1]
		numero_chunks, _ := strconv.Atoi(cantidad_chunks)
		if nombre_libro == libro_solicitado.Nombre {
			var j = 0
			for j < int(numero_chunks) {
				i++
				j++
				s := strings.Split(file_lines[i], " ")
				nombre_chunk, ip_maquina := s[0], s[1]
				//log.Printf("nombre_chunk: %s\n", nombre_chunk)
				retorno = append(retorno, &ChunkIp{Nombrechunk: nombre_chunk, Ipmaquina: ip_maquina})
			}
			return &DistribucionChunks{Nombrelibro: libro_solicitado.Nombre, Chunks: retorno}, nil
		} else {
			i += int(numero_chunks)
		}

	}
	return &DistribucionChunks{}, nil
}

func (s *Server) SolicitarPropuesta(ctx context.Context, chunks_nombres *DistribucionChunks) (*DistribucionChunks, error) {

	//totalPartsNum := int(len(chunks_nombres.Chunks))

	//Generar nueva propuesta

	return chunks_nombres, nil
}

// Pedir lista de libros
