package main

import (
    "fmt"
    "strings"
    "strconv"
    "reflect"
    "net"
)

// Client is the main structure for the client
type Client struct {
    conn net.Conn // Connexion avec le serveur
    temps int     // Temps du client
    runner int    // Skin du client
    posX float64  // Position X du client
}

// Server is the main structure for the server
type Server struct {
    port string
    listener net.Listener
}

// Tableau de clients
var clients []Client

type GameState int
const (
    connexion GameState = iota
    choixPersos
    course
    score
)
var state GameState = connexion

const nbClients int = 4
const posXDepart float64 = 50.0

// NewServer creates a new server
func NewServer(port string) *Server {
    return &Server{port: port}
}

// Start starts the server
func (s *Server) Start() error {
    var err error
    s.listener, err = net.Listen("tcp", s.port)
    if err != nil {
        return err
    }
    fmt.Println("Server started on 127.0.0.1" + s.port)
    return nil
}

// Stop stops the server
func (s *Server) Stop() error {
    return s.listener.Close()
}


// Run runs the server
func (s *Server) Run() error {
    for { // Boucle infinie

        // Si le state est la connexion
        if state == connexion {
            // On accepte la connexion
            conn, err := s.listener.Accept()
            if err != nil {
                return err
            }
            fmt.Println("New client")

            // On crée un nouveau client
            client:= Client{conn: conn, temps: -1, runner: -1, posX: posXDepart}
            // On affiche l'adresse du client
            fmt.Println(conn.RemoteAddr().String())

            // On lance la goroutine qui va récupérer et gérer les messages du client
            go s.handleConnection(client)

            // On ajoute le client au tableau de clients
            clients = append(clients, client)

            // Si la partie est full après l'ajout du client
            if (len(clients) == nbClients) {
                // Alors on start le jeu
                sendNextState("")
            } else {
                sendNbClientPrets()
            }
        }

    }
}

/**
* Gestion des messages du client
*/
func (s *Server) handleConnection(client Client) {
    // On récupère les infos utiles
    conn := client.conn
    addr := conn.RemoteAddr().String()

    // On ferme la connexion à la fin de la fonction
    defer conn.Close()

    // Boucle infinie
    for {
        // On lit le message du client
        buf := make([]byte, 1024)
        n, err := conn.Read(buf)
        if err != nil {
            return
        }

        // On affiche le message du client
        fmt.Println("Client " + addr + " a écrit : " + string(buf[:n]))

        // On split le message en deux parties (séparées par "::")
        str := strings.Split(string(buf[:n]), "\n")[0]
        split := strings.Split(str, "::")

        // En fonction du type de message on appelle la fonction correspondante
        switch split[0] {
            case "runner": // State ChoixPersos
                saveRunner(client, split[1])
            case "temps": // State Course
                saveTemps(client, split[1])
            case "restart":
                saveRestart(client)
            case "position":
                savePosition(client, split[1])
        }
    }
}

func saveRunner(client Client, content string) {
    addr := client.conn.RemoteAddr().String()
    runner, err := strconv.Atoi(content)
    fmt.Println("Client " + addr + " a choisi le runner " + content)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Affichage du runner et du type de runner
    fmt.Println(runner)
    fmt.Println(reflect.TypeOf(runner))

    for i := range clients {
        if clients[i].conn.RemoteAddr().String() == addr {
            clients[i].runner = runner
            fmt.Println("Client " + addr + " a choisi le runner " + string(client.runner))
        }
    }

    // Si tous les clients on un runner alors on envoie les runners aux clients
    if allClientsHaveRunner() {
        sendRunners()
    } else {
        sendNbClientPrets()
    }
}

func saveTemps(client Client, temps string) {
    addr := client.conn.RemoteAddr().String()
    tempsInt, err := strconv.Atoi(temps)
    fmt.Println("Client " + addr + " a fini en " + temps)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Affichage du temps et du type de tempsInt
    fmt.Println(tempsInt)
    fmt.Println(reflect.TypeOf(tempsInt))

    for i := range clients {
        if clients[i].conn.RemoteAddr().String() == addr {
            clients[i].temps = tempsInt
            fmt.Println("Client " + addr + " a fini en " + string(client.temps))
        }
    }

    // Si tous les clients on un temps alors on envoie les temps aux clients
    if allClientsHaveTemps() {
        sendTemps()
    } else {
        sendNbClientPrets()
    }
}

func saveRestart(client Client) {
    addr := client.conn.RemoteAddr().String()

    // On reset
    for i := range clients {
        if clients[i].conn.RemoteAddr().String() == addr {
            clients[i].temps = -1
            clients[i].posX = posXDepart // Position de départ
        }
    }
    // On envoie le restart à tout le monde si tout le monde à restart
    if allClientsHaveRestart() {
        sendRestart()
    } else {
        sendNbClientPrets()
    }
}

func savePosition(client Client, position string) {
    addr := client.conn.RemoteAddr().String()
    posXFloat, err := strconv.ParseFloat(position, 64)

    if err != nil {
        fmt.Println(err)
        return
    }

    for i := range clients {
        if clients[i].conn.RemoteAddr().String() == addr {
            clients[i].posX = posXFloat
        }
    }

    // Envoie la position à tout le monde
    sendPosition()
}

func allClientsHaveTemps() bool {
    for i := range clients {
        if clients[i].temps == -1 && clients[i].runner != -1 {
            return false
        }
    }
    return true
}

func allClientsHaveRunner() bool {
    for i := range clients {
        if clients[i].runner == -1 {
            return false
        }
    }
    return true
}

func allClientsHaveRestart() bool {
    for i := range clients {
        if clients[i].temps != -1 && clients[i].posX != posXDepart {
            return false
        }
    }
    return true
}

func sendTemps() {
    str := ""
    for i := range clients {
        str += clients[i].conn.RemoteAddr().String() + "-" + strconv.Itoa(clients[i].temps) + ";"
    }

    // On passe à l'état suivant
    sendNextState(str)
}

func sendRunners() {
    str := ""

    for i := range clients {
        str += clients[i].conn.RemoteAddr().String() + "-" + strconv.Itoa(clients[i].runner) + ";"
    }

    // On passe à l'état suivant
    sendNextState(str)
}

func sendPosition() {
    str := ""

    for i := range clients {
        str += clients[i].conn.RemoteAddr().String() + "-" + fmt.Sprintf("%f",clients[i].posX) + ";"
    }

    for i := range clients {
        sendMessage(clients[i], "position::" + str)
    }
}

func sendRestart() {
    // On passe à l'état suivant
    sendNextState("")
}


// Envoie l'état actuel aux clients
func sendNextState(content string) {
    // Si tous les clients sont prêts alors on passe à l'état suivant
    state++

    fmt.Println("state::" + content)

    for i := range clients {
        sendMessage(clients[i], "state::" + content)
    }

    // Si on est à l'état des résultats on revient à l'état de connexion
    if state == score {
        state = choixPersos
    }
}

func sendNbClientPrets() {
    nbClientsPrets := 0

    for i := range clients {
        switch state {
            case connexion:
                nbClientsPrets++
            case choixPersos:
                if clients[i].runner != -1 && clients[i].temps == -1 {
                    nbClientsPrets++
                }
            case course:
                if clients[i].temps != -1 {
                    nbClientsPrets++
                }
        }
    }

    for i := range clients {
        sendMessage(clients[i], "nbClientsPrets::" + strconv.Itoa(nbClientsPrets))
    }
}

func sendMessage(client Client, message string) {
    client.conn.Write([]byte(message))
}

func main() {
    clients = make([]Client, 0)
    s := NewServer(":8080")
    if err := s.Start(); err != nil {
        fmt.Println(err)
        return
    }
    if err := s.Run(); err != nil {
        fmt.Println(err)
        return
    }
}