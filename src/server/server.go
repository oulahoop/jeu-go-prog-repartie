package main

import (
    "fmt"
    "net"
)

// Client is the main structure for the client
type Client struct {
    conn net.Conn
}

// Server is the main structure for the server
type Server struct {
    port string
    listener net.Listener
}

// Tableau de clients
var clients []Client

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
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            return err
        }
        fmt.Println("New client")
        go s.handleConnection(conn)

        clients = append(clients, Client{conn: conn})

        fmt.Println(len(clients))
        if (len(clients) == 4) {
            for i := range clients {
                clients[i].conn.Write([]byte("4 clients"))
            }
        }
    }
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()

    for {
        buf := make([]byte, 1024)
        n, err := conn.Read(buf)
        if err != nil {
            return
        }
        fmt.Println(string(buf[:n]))
    }
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