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

	totalPartsNum := int(len(chunks_nombres.Chunks))

	//////          Validar      //////////
	//Debido a que la generación de la propuesta tenderá a siempre ser correcta,
	// Se hará un random para rechazarla, 30% de prob
	log.Println("Validando propuesta.")

	///// Generar nueva propuesta ///////

	//Escriber propuesta en el final del archivo logs
	f, err := os.OpenFile("name_node/logs/logs.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	fi, err := os.Stat("name_node/logs/logs.txt")
	if fi.Size() != 0 {
		if _, err := f.WriteString("\n"); err != nil {
			log.Println(err)
		}
	}
	// Escribir nombre + cantidad de chunks
	if _, err := f.WriteString(chunks_nombres.Nombrelibro + " " + strconv.Itoa(totalPartsNum) + "\n"); err != nil {
		log.Println(err)
	}

	for i, chunk := range chunks_nombres.Chunks {

		// Escribir 1 chunk en logs
		if _, err := f.WriteString(chunk.Nombrechunk + " " + chunk.Ipmaquina); err != nil {
			log.Println(err)
		}
		if i != len(chunks_nombres.Chunks)-1 {
			// Escribir salto de linea si no es la ultima linea
			if _, err := f.WriteString("\n"); err != nil {
				log.Println(err)
			}
		}
	}

	log.Println("Se escribio en logs.")

	return chunks_nombres, nil
}

// Pedir lista de libros
