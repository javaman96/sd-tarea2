package main

import (    
    "bufio"
    "fmt"      
    "os"
    "log"
    "strings"
    "strconv"
    "archive/zip"
    "io"
    "path/filepath"
)

func unzip(zipfile string){            

    reader, err := zip.OpenReader(zipfile)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer reader.Close()

    for _, f := range reader.Reader.File {

        zipped, err := f.Open()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        defer zipped.Close()

        // get the individual file name and extract the current directory
        path := filepath.Join("./", f.Name)

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
            fmt.Println("Creating directory", path)
        } else {
            writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, f.Mode())

            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            defer writer.Close()

            if _, err = io.Copy(writer, zipped); err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            fmt.Println("Decompressing : ", path)
        }
    }
}

func main() {
    newFileName := "MagoOz.zip"
    _, err := os.Create(newFileName)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    //set the newFileName file to APPEND MODE!!
    // open files r and w

    file, err := os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // IMPORTANT! do not defer a file.Close when opening a file for APPEND mode!
    // defer file.Close()

    // just information on which part of the new file we are appending
    var writePosition int64 = 0

    // Cambiar esta forma por un identificador 
    // para todos los chunks del mismo libro
    // Igual como el ensamblado de fragmentos de capa IP
    ////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////

    totalPartsNum := 0
    chunk_dir, err := os.Open(".")
    if err != nil {
        log.Fatalf("failed opening directory: %s", err)
    }
    defer chunk_dir.Close()
 
    chunk_list,_ := chunk_dir.Readdirnames(0) // 0 to read all files and folders

    var input string
    // Pensarlo como un diccionario en el NameNode nombre: Id, o al reves
    fmt.Printf("\n Ingrese Id del archivo a descargar: ")
    fmt.Scanln(&input)

    for _, name := range chunk_list {
        if strings.Contains(name, input + "_") {
            totalPartsNum++
        }
    }
    ////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////    

    for j := uint64(0); j < uint64(totalPartsNum); j++ {

        //read a chunk
        currentChunkFileName := input + "_" + strconv.FormatUint(j, 10)

        newFileChunk, err := os.Open(currentChunkFileName)

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        defer newFileChunk.Close()

        chunkInfo, err := newFileChunk.Stat()

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // calculate the bytes size of each chunk
        // we are not going to rely on previous data and constant

        var chunkSize int64 = chunkInfo.Size()
        chunkBufferBytes := make([]byte, chunkSize)

        //fmt.Println("Appending at position : [", writePosition, "] bytes")
        writePosition = writePosition + chunkSize

        // read into chunkBufferBytes
        reader := bufio.NewReader(newFileChunk)
        _, err = reader.Read(chunkBufferBytes)

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
        // write/save buffer to disk
        //ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)

        n, err := file.Write(chunkBufferBytes)

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        file.Sync() //flush to disk

        // free up the buffer for next cycle
        // should not be a problem if the chunk size is small, but
        // can be resource hogging if the chunk size is huge.
        // also a good practice to clean up your own plate after eating

        chunkBufferBytes = nil // reset or empty our buffer

        fmt.Println("Written ", n, " bytes")

        //fmt.Println("Recombining part [", j, "] into : ", newFileName)
    }

    // now, we close the newFileName
    file.Close()

    // Decompress
    unzip("MagoOz.zip")

    // Delete remaining zip files
    chunk_dir, err = os.Open(".")
    if err != nil {
        log.Fatalf("failed opening directory: %s", err)
    }
    defer chunk_dir.Close()
 
    chunk_list,_ = chunk_dir.Readdirnames(0) // 0 to read all files and folders

    for _, name := range chunk_list {
        if strings.Contains(name, ".zip") {
            os.Remove(name)
        }
    }
}