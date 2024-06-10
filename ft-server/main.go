package main 

import (
    "fmt"
    "os"
    "bufio"
)

func file_io() {
    file, err := os.Open("./resource/info.txt")
    
    if err != nil {
        fmt.Println("open error:", err)
        return 
    }
    defer file.Close()

    data := make([]byte, 100)
    
    count, err := file.Read(data)
    if err != nil {
        fmt.Println("read err", err)
        return
    }

    fmt.Printf("read flie %d byte, the content is : %q\n", count, data[:count])

    fmt.Println("begin to copy file")

    file_cp, err := os.OpenFile("./resource/copy.txt", os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0644)
    if err != nil {
        fmt.Println("OpenFile error:", err)
        return
    }
    defer file_cp.Close()

    if _, err := file_cp.Write(data); err != nil {
        fmt.Println("Write error:", err)
        return
    }
}

func buffer_io() {
    file, err := os.Open("./resource/info.txt")
    if err != nil {
        fmt.Println("Open Error:", err)
        return
    }
    defer file.Close()

    reader := bufio.NewReader(file)

    file_cp, err := os.OpenFile("./resource/copy.txt", os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0644) 
    if err != nil {
        fmt.Println("OpenFile Error:", err) 
        return
    }
    defer file_cp.Close()

    writer := bufio.NewWriter(file_cp)

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            if err.Error() == "EOF" {
                fmt.Println("finish copy")
                writer.Flush()
                break
            }
            fmt.Println("ReadString Error:", err)
            return
        }
        fmt.Print(line)
        writer.WriteString(line)
    }

    // !!!
}

func test_io() {
    buffer_io()
}
