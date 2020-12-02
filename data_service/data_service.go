package data_service

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/dcordova/sd_tarea2/name_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const ip_datanode_1 = ":9001"
const ip_datanode_2 = ":9002"
const ip_datanode_3 = ":9003"
const ip_namenode = ":9009"

// Guardar los chunks en el disco de la machina
func (s *Server) save_chunks(chunk_name string, data []byte) {
	//fmt.Println(chunk_name)
	dir := "data_service/chunks/" + chunk_name + ".chunk"

	// create file
	if err := os.MkdirAll(filepath.Dir(dir), 0777); err != nil {
		log.Fatal(err)
	}

	// write/save buffer to disk
	ioutil.WriteFile(dir, data, os.ModeAppend)
}

type Server struct {
}

// Esta funcion se encarga de generar (o solicitar) una propeusta
// UNa vez validada se guardan los chunks en cada máquina

func (s *Server) GenerarPropuesta(chunks []*Chunk) {

	//"nombre/lugar; nombre/lugar; nombre/lugar; nombre/lugar; "

	// Extraer nombre del libro ya que viene en hexadecimal
	name_hex := strings.Split(chunks[0].Id, "_")[0]
	book_name_bytes, err := hex.DecodeString(name_hex)
	if err != nil {
		log.Fatal(err)
	}
	book_name := string(book_name_bytes)

	totalPartsNum := len(chunks)
	var centralizado = true

	// Generar propuesta
	propuesta := name_service.DistribucionChunks{}
	propuesta.Nombrelibro = book_name

	for i := 0; i < totalPartsNum; i++ {
		if i < int(math.Ceil(float64(totalPartsNum)/3)) {

			propuesta.Chunks = append(propuesta.Chunks, &name_service.ChunkIp{Nombrechunk: chunks[i].Id, Ipmaquina: ip_datanode_1})

		} else if i < int(2*math.Ceil(float64(totalPartsNum)/3)) {

			propuesta.Chunks = append(propuesta.Chunks, &name_service.ChunkIp{Nombrechunk: chunks[i].Id, Ipmaquina: ip_datanode_2})

		} else if i < int(3*math.Ceil(float64(totalPartsNum)/3)) {

			propuesta.Chunks = append(propuesta.Chunks, &name_service.ChunkIp{Nombrechunk: chunks[i].Id, Ipmaquina: ip_datanode_3})
		}
	}

	// Validar propuesta
	if centralizado {
		//Pedir propuesta al name_node
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(ip_namenode, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		s_name_node := name_service.NewNameServiceClient(conn)

		// Se pide la PROPUESTA y se ESPERA
		response, err := s_name_node.SolicitarPropuesta(context.Background(), &propuesta)
		if err != nil {
			log.Fatalf("Error al pedir propuesta a namenode: %s", err)
		}
		//fmt.Println("respuesta propuesta: ", response)
		propuesta = name_service.DistribucionChunks{Nombrelibro: response.Nombrelibro,
			Chunks: response.Chunks}
	} else {
		fmt.Println("Distribuido aun no implementado")
	}

	// Una vez validada la propuesta se guardan los chunks en los nodos
	for i := 0; i < len(propuesta.Chunks); i++ {

		nombre_chunk := propuesta.Chunks[i].Nombrechunk
		ip_chunk := propuesta.Chunks[i].Ipmaquina

		fmt.Println("Nombre chunk: ", nombre_chunk, "  Ip:", ip_chunk)

		if ip_chunk == ip_datanode_1 {
			s.save_chunks(nombre_chunk, chunks[i].Data)
		} else {
			var conn2 *grpc.ClientConn
			conn2, err := grpc.Dial(ip_chunk, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("Could not connect to %s: %s", ip_chunk, err)
			}
			defer conn2.Close()

			c2 := NewDataServiceClient(conn2)

			_, err = c2.SendChunks(context.Background(), chunks[i])
			if err != nil {
				log.Fatalf("Error when calling Server 2: %s", err)
			}
		}

		fmt.Println("Se mando el chunk: ", nombre_chunk)
		fmt.Println(" a la maquina: ", ip_chunk)

	}
}

func (s *Server) UploadChunks(stream DataService_UploadChunksServer) error {

	chunkArr := []*Chunk{}
	for {

		chunk, err := stream.Recv()

		if err == io.EOF {
			go s.GenerarPropuesta(chunkArr)
			return stream.SendAndClose(&Message{Body: "Received!"})
		}
		if err != nil {
			return err
		}
		//save_chunks(string(chunk.Id), chunk.Data)
		chunkArr = append(chunkArr, chunk)
		fmt.Println(chunk.Id)
	}
}

func (s *Server) SendChunks(ctx context.Context, chunk *Chunk) (*Message, error) {

	s.save_chunks(string(chunk.Id), chunk.Data)
	fmt.Println(chunk.Id)
	return &Message{Body: "Recibido!"}, nil
}
