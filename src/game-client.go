package main

import (
    "fmt"
    "net"
)

// Struct client
type Client struct {
    conn net.Conn // Connexion avec le serveur
}

func (g *Game) connectToServer() {
    // Se connecte au serveur sur le port 8080
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println(err)
        return
    }
    // Lie le client au jeu
    g.Client = Client{conn: conn}

    go g.handleConnection()

    defer conn.Close()
}

func (g *Game) handleConnection() {
    for {
        buf := make([]byte, 1024)
        n, err := g.Client.conn.Read(buf)
        if err != nil {
            return
        }
        value := string(buf[:n])
        fmt.Println(value)
        if value == "4 clients" {
            g.Client.conn.Write([]byte("4 clients"))
        }
    }
}