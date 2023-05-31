package main

import (
    "fmt"
    "strings"
    "strconv"
    "time"
    "net"
)

// Struct client
type Client struct {
    conn net.Conn // Connexion avec le serveur
}

/**
* Se connecte au serveur
*/
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

    for {
    }
}

func (g *Game) SendRunner() {
    // Envoie au serveur le choix du personnage
    runner := g.runners[0].colorScheme

    fmt.Println("runner::" + strconv.Itoa(runner))
    _, err := g.Client.conn.Write([]byte("runner::" + strconv.Itoa(runner) + "\n"))
    if err != nil {
        fmt.Println(err)
        return
    }
}

/**
* Attend les messages du serveur
*/
func (g *Game) handleConnection() {
    for {
        buf := make([]byte, 1024)
        n, err := g.Client.conn.Read(buf)
        if err != nil {
            return
        }

        value := string(buf[:n])

        split := strings.Split(value, "\n")
        split = strings.Split(split[0], "::")

        fmt.Println("Received: " + split[0] + " - " + split[1])

        switch split[0] {
            case "state":
                g.UpdateState(split[1])
        }
    }
}

func (g *Game) UpdateState(content string) {
    currentState := g.stateServer

    switch currentState {
        case 0: // Connexion
            g.stateServer = 1
        case 1: // Choix du personnage
            // Enregistrement des skins choisis
            g.RetrieveSkins(content)
            g.stateServer = 2
            fmt.Println("stateServer = 2")
    }
}

func (g *Game) RetrieveSkins(skins string) {
    // skins = "addr:port-skin;addr:port-skin;..."
    client := g.Client
    split := strings.Split(skins, ";")
    index := 1
    for i := range split {
        if (split[i] == "") {
            continue
        }

        fmt.Println("split[i] = " + split[i])
        split2 := strings.Split(split[i], "-")
        addr := split2[0]
        skin, err := strconv.Atoi(split2[1])

        if err != nil {
            fmt.Println(err)
            return
        }

        fmt.Println("addr: " + addr + " = " + client.conn.LocalAddr().String())
        if (addr == client.conn.LocalAddr().String()) {
            g.runners[0].colorScheme = skin
        } else {
            g.runners[index].colorScheme = skin
            index++
        }
    }
}

func (g *Game) RetrieveTemps(temps string) {
    // temps = "addr:port-temps;addr:port-temps;..."
    client := g.Client
    split := strings.Split(temps, ";")
    index := 1
    for i := range split {
        if (split[i] == "") {
            continue
        }
        split2 := strings.Split(split[i], "-")
        addr := split2[0]
        tempsInt, err := strconv.Atoi(split2[1])

        if err != nil {
            fmt.Println(err)
            return
        }

        fmt.Println("addr: " + addr + " = " + client.conn.LocalAddr().String())
        if (addr == client.conn.LocalAddr().String()) {
            // Convertis tempsInt (qui est en MS) en Duration
            g.runners[0].runTime = time.Duration(tempsInt) * time.Millisecond
            fmt.Println("Player 0 time is " + strconv.Itoa(int(g.runners[0].runTime.Milliseconds())))
        } else {
            g.runners[index].runTime = time.Duration(tempsInt) * time.Millisecond
            fmt.Println("Player " + strconv.Itoa(index) + " time is " + strconv.Itoa(int(g.runners[index].runTime.Milliseconds())))
            index++
        }
    }
}


/**
* Envoie le temps réalisé au serveur
*/
func (g *Game) SendResults() {
    // Récupération du temps du runner courrant
    ms := g.runners[0].runTime.Milliseconds()
    msString := strconv.Itoa(int(ms))

    fmt.Println("Sending results (" + msString + ")")

    // Envoie le temps au serveur
    g.Client.conn.Write([]byte("temps::" + msString + "\n"))
}