package main

import (
    "fmt"
    "net"
)

// Struct client
type Client struct {
    conn net.Conn // Connexion avec le serveur
}

func (g *Game) ConnectToServer() {
    // Se connecte au serveur sur le port 8080
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println(err)
        return
    }
    // Lie le client au jeu
    g.Client = &Client{conn: conn}

    // Lance la goroutine qui attend les messages du serveur
    go g.handleConnection()

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
        if value == "start" {
            g.stateServer = 4
        }
    }
}


func (g *Game) SendResults() {

    fmt.Println("Sending results (" + string(g.runners[0].runTime.Milliseconds()) + ")")
    // Results
    s, ms := GetSeconds(g.runners[0].runTime.Milliseconds())

    // Envoi des résultats au server
    g.Client.conn.Write([]byte(string(s) + "-" + string(ms)))
}

func (g *Game) sendRun() {
    g.Client.conn.Write([]byte("run"))
}

func (g *Game) finishRun(temps int) {
    // temps en string
    tps := fmt.Sprintf("%d", temps)
    g.Client.conn.Write([]byte(tps))
}