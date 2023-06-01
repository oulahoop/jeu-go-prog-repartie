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
    conn, err := net.Dial("tcp", "172.21.65.25:8080")
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

    _, err := g.Client.conn.Write([]byte("runner::" + strconv.Itoa(runner) + "\n"))
    if err != nil {
        fmt.Println(err)
        return
    }
}

/**
* Ecoute les messages du serveur
*/
func (g *Game) handleConnection() {
    for {
        // Récupération du message du serveur
        buf := make([]byte, 1024)
        n, err := g.Client.conn.Read(buf)
        if err != nil {
            return
        }

        // Protocole : "type::content\n"
        value := string(buf[:n])
        split := strings.Split(value, "\n")
        split = strings.Split(split[0], "::")

        // Affichage du message reçu
        fmt.Println("Received: " + split[0] + " - " + split[1])

        // En fonction du type de message on appelle la fonction correspondante
        switch split[0] {
            case "state":
                g.UpdateState(split[1])
            case "position":
                g.RetrievePosition(split[1])
            case "nbClientsPrets":
                joueurPret, _ := strconv.Atoi(split[1])
                g.joueurPret = joueurPret
            case "deplacementMenuRunner" :
                g.RetrieveDeplacementMenuRunner(split[1])
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
        case 2: // Fin de la partie
            // Enregistrement des temps
            g.RetrieveTemps(content)
            g.stateServer = 3
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

        split2 := strings.Split(split[i], "-")
        addr := split2[0]
        skin, err := strconv.Atoi(split2[1])

        if err != nil {
            fmt.Println(err)
            return
        }

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

        if (addr == client.conn.LocalAddr().String()) {
            // Convertis tempsInt (qui est en MS) en Duration
            g.runners[0].runTime = time.Duration(tempsInt) * time.Millisecond
        } else {
            g.runners[index].runTime = time.Duration(tempsInt) * time.Millisecond
            index++
        }
    }
}

func (g *Game) RetrievePosition(positions string) {
    // pos = "addr:port-posX;addr:port-posX;..."
    client := g.Client
    split := strings.Split(positions, ";")
    index := 1

    for i := range split {
        if (split[i] == "") {
            continue
        }
        split2 := strings.Split(split[i], "-")

        if (len(split2) < 2) {
            continue
        }

        addr := split2[0]
        posX, err := strconv.ParseFloat(split2[1], 64)

        if err != nil {
            fmt.Println(err)
            return
        }

        if (addr == client.conn.LocalAddr().String()) {
            // On ne modifie pas la position du joueur local
        } else {
            g.runners[index].xpos = posX
            index++
        }
    }
}

func (g *Game) SendPosition() {
    // Envoie au serveur la position du joueur
    _, err := g.Client.conn.Write([]byte("position::" + fmt.Sprintf("%f", g.runners[0].xpos) + "\n"))
    if err != nil {
        fmt.Println(err)
        return
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

func (g *Game) RestartGame() {
    // Envoie au serveur le choix du personnage
    _, err := g.Client.conn.Write([]byte("restart\n"))
    if err != nil {
        fmt.Println(err)
        return
    }
}

func (g *Game) RetrieveDeplacementMenuRunner(content string) {
    // content = "addr:port-preselection;addr:port-preselection;..."
    clientAddr := g.Client.conn.LocalAddr().String()
    index := 1

    split := strings.Split(content, "\n")
    split = strings.Split(split[0], ";")

    for i := range split {
        if (split[i] == "") {
            continue
        }

        split2 := strings.Split(split[i], "-")
        addr := split2[0]
        preselection, err := strconv.Atoi(split2[1])

        if err != nil {
            fmt.Println(err)
            return
        }

        if (addr != clientAddr) {
            g.runners[index].colorScheme = preselection
            index++;
        }
    }

}

func (g *Game) SendDeplacementMenuRunner() {
    // Envoie au serveur le choix du personnage
    _, err := g.Client.conn.Write([]byte("deplacementMenuRunner::" + strconv.Itoa(g.runners[0].colorScheme) + "\n"))
    if err != nil {
        fmt.Println(err)
        return
    }
}