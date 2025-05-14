package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"searchx-indexer/banner"
	"searchx-indexer/database"
	"searchx-indexer/handlers"
	"searchx-indexer/socket"
	"syscall"
)

func init() {
	go exit()
	banner.Show()
}

func main() {
	address := flag.String("address", "0.0.0.0", "set listening address (default: 0.0.0.0)")
	port := flag.String("port", "5000", "set listening port (default: 5000)")
	httpPort := flag.String("http-port", "8080", "set HTTP server port (default: 8080)")
	protocol := flag.String("protocol", "TCP", "set listening protocol (default: tcp)")

	flag.Parse()

	manager := database.GetManager()
	defer manager.Close()

	// Conectar ao banco de dados combolist
	if err := manager.Connect("combolist-db", "postgres", "127.0.0.1", 5433, "docker", "docker", "searchx_combolist"); err != nil {
		log.Println("[!]", err)
	}

	// Conectar ao banco de dados searchx
	if err := manager.Connect("searchx-db", "postgres", "127.0.0.1", 5435, "docker", "docker", "searchx"); err != nil {
		log.Println("[!]", err)
	}

	// Conectar ao banco de dados discord
	if err := manager.Connect("discord-db", "postgres", "127.0.0.1", 5434, "docker", "docker", "discord_db"); err != nil {
		log.Println("[!]", err)
	}

	// Iniciar servidor TCP
	go func() {
		ln := socket.Listen(*address, *port, *protocol)
		for {
			socket.SocketAccept(ln)
		}
	}()

	// Iniciar servidor HTTP
	http.HandleFunc("/discord", handlers.HandleDiscordMessage)
	log.Printf("[+] HTTP server listening on port %s\n", *httpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", *httpPort), nil); err != nil {
		log.Fatalf("[!] Failed to start HTTP server: %s\n", err)
	}
}

func exit() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("[+] Signal received: %s.\n", sig)

	database.GetManager().Close()

	os.Exit(0)
}
