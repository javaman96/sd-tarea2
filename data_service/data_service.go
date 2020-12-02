package data_service

import (
    "log"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "strconv"
    "math"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
)




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


func (s *Server) read_chunks(currentChunkFileName string) {
    
    //read a chunk
    //currentChunkFileName := input + "_" + strconv.FormatUint(j, 10)

    newFileChunk, err := os.Open(currentChunkFileName)

    if err != nil {
        return err
        os.Exit(1)
    }

    defer newFileChunk.Close()

    chunkInfo, err := newFileChunk.Stat()

    if err != nil {
        return err
        os.Exit(1)
    }

    // calculate the bytes size of each chunk
    // we are not going to rely on previous data and constant

    var chunkSize int64 = chunkInfo.Size()
    chunkBufferBytes := make([]byte, chunkSize)

    //fmt.Println("Appending at position : [", writePosition, "] bytes")
    //writePosition = writePosition + chunkSize

    // read into chunkBufferBytes
    reader := bufio.NewReader(newFileChunk)
    _, err = reader.Read(chunkBufferBytes)

    if err != nil {
        return err
        os.Exit(1)
    }  
    return chunkBufferBytes
}


type Server struct {	
}


func (s *Server) GenerarPropuesta(chunks []*Chunk) {

    //"nombre/lugar; nombre/lugar; nombre/lugar; nombre/lugar; "

    totalPartsNum := len(chunks)

    for i := 0; i < totalPartsNum; i++ {
        if i < int(math.Ceil(float64(totalPartsNum)/3)){    

            fmt.Println("S1: ", strconv.Itoa(i))

        }else if i < int(2*math.Ceil(float64(totalPartsNum)/3)){   
            
            var conn2 *grpc.ClientConn
            conn2, err := grpc.Dial(":9002", grpc.WithInsecure())
            if err != nil {
                log.Fatalf("Could not connect to 9002: %s", err)
            }
            defer conn2.Close()

            c2 := NewDataServiceClient(conn2)

            _, err = c2.SendChunks(context.Background(), chunks[i])
            if err != nil {
                log.Fatalf("Error when calling Server 2: %s", err)
            } 
               
            fmt.Println("S2: ", strconv.Itoa(i)) 

        }else if i < int(3*math.Ceil(float64(totalPartsNum)/3)){
            
            var conn3 *grpc.ClientConn
            conn3, err := grpc.Dial(":9003", grpc.WithInsecure())
            if err != nil {
                log.Fatalf("Could not connect to 9003: %s", err)
            }
            defer conn3.Close()

            c3 := NewDataServiceClient(conn3)

            _, err = c3.SendChunks(context.Background(), chunks[i])
            if err != nil {
                log.Fatalf("Error when calling Server 3: %s", err)
            }
            
            fmt.Println("S3: ", strconv.Itoa(i))
        }
    }
}



func (s *Server) UploadChunks(stream DataService_UploadChunksServer) (error) {    

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

func (s *Server) RecuperarChunks(ctx context.Context, message *Message) (*Chunk, error) {

    data := s.read_chunks(string(chunk.Id))    
    return &Chunk{Id: message, Data: data}, nil
}