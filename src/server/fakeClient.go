package main

import (
    "fmt"
    "net"
    "flag"
)

// Fonction goroutin qui attends la reception d'un message du serveur et l'affiche
func handleConnection(conn net.Conn, ch chan bool) {
    defer conn.Close()
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(buf[:n]))
}

func main() {

    //Flag pour préciser l'adresse/port
    var ip string = ""

    flag.StringVar(&ip, "ip", "127.0.0.1:8080" ,"Précise l'ip du serveur")
    flag.Parse()


    // Se connecte au serveur sur le port 8080
    conn, err := net.Dial("tcp", ip)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Nouveau chan
    ch := make(chan bool)

    // Attend la reception d'un message du serveur
    go handleConnection(conn, ch)

    <-ch
}