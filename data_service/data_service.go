package data_service

import (
    "log"
    //"fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "golang.org/x/net/context"
)




// Guardar los chunks en el disco de la machina
func save_chunks(chunk_name string, data []byte) {
	//fmt.Println(chunk_name)
	dir := "data_service/chunks/" + chunk_name + ".chunk"

	// create file
	if err := os.MkdirAll(filepath.Dir(dir), 0755); err != nil {
        log.Fatal(err)
    }    	
        
	// write/save buffer to disk
	ioutil.WriteFile(dir, data, os.ModeAppend)	
}

type Server struct {	
}

func (s *Server) UploadChunks(ctx context.Context, message *Book) (*Book, error) {

	res := make([]byte, 1) // inutil respuesta

	save_chunks(string(message.Chunks), message.Data)

    log.Printf("Received chunk: %s", message.Chunks)
    //fmt.Println(message.Data)
    return &Book{Chunks: message.Chunks, Data: res}, nil //&Book{}, nil 
}
