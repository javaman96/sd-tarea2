package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "math"
    "os"
    "io"
    "strconv"
    "strings"
    "archive/zip"
    "encoding/hex"
)

func encodeString(titulo string) string {
    src := []byte(titulo)
    encodedStr := hex.EncodeToString(src)
    return encodedStr
}

func ZipFile(filename string) error {      

    newfile, err := os.Create(strings.Split(filename,".pdf")[0] + ".zip")
    if err != nil {
        return err
    }
    defer newfile.Close()

    zipit := zip.NewWriter(newfile)

    defer zipit.Close()

    zipfile, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer zipfile.Close()

    // get the file information
    info, err := zipfile.Stat()
    if err != nil {
        return err
    }

    header, err := zip.FileInfoHeader(info)
    if err != nil {
        return err
    }

    header.Method = zip.Deflate

    writer, err := zipit.CreateHeader(header)
    if err != nil {
        return err
    }
    _, err = io.Copy(writer, zipfile)
    return err
}

func main() {

    var input string
    fmt.Printf("\n Ingrese nombre completo del archivo (archivo.pdf): ")
    fmt.Scanln(&input)

    input = strings.TrimSpace(input)

    chunkname := encodeString(input) + "_"
    fmt.Println(chunkname)

    err := ZipFile(input)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fileToBeChunked := "./"+strings.Split(input,".pdf")[0] + ".zip"    

    file, err := os.Open(fileToBeChunked)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    defer file.Close()

    fileInfo, _ := file.Stat()

    var fileSize int64 = fileInfo.Size()

    fmt.Println(fileSize)

    const fileChunk = 250 * (1 << 10) // 250KB

    // calculate total number of parts the file will be chunked into

    totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

    fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

    for i := uint64(0); i < totalPartsNum; i++ {

        partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
        partBuffer := make([]byte, partSize)

        file.Read(partBuffer)

        // write to disk
        fileName := chunkname + strconv.FormatUint(i, 10)
        _, err := os.Create(fileName)

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // write/save buffer to disk
        ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)

        fmt.Println("Split to: ", fileName)
    }

    // Delete remaining zip files
    chunk_dir, err := os.Open(".")
    if err != nil {
        log.Fatalf("failed opening directory: %s", err)
    }
    defer chunk_dir.Close()
 
    chunk_list,_ := chunk_dir.Readdirnames(0) // 0 to read all files and folders

    for _, name := range chunk_list {
        if strings.Contains(name, ".zip") {
            os.Remove(name)
        }
    }
}