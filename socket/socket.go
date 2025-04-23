package socket

import (
	"log"
	"net"
	"searchx-indexer/handlers"
)

func Listen(address, port, protocol string) net.Listener {
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		log.Fatalf("[!] Error starting socket listener: %v", err)
	}
	log.Printf("[*] Server on %s:%s (%s)\n", address, port, protocol)
	return listener
}

func SocketAccept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("[!] Error accepting connection:", err)
			continue
		}
		log.Println("[+] New connection accepted")
		go SocketConnection(conn)
	}
}

func SocketConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("[!] Error reading from socket:", err)
			return
		}

		handlers.ProcessMessage(buffer[:n])
	}
}
