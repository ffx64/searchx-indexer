package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"searchx-indexer/banner"
	"searchx-indexer/database"
	"searchx-indexer/socket"
	"syscall"
)

var (
	VersionNumber = "0.0"
	VersionName   = ""
)

func init() {
	go exit()

	banner.Show()
}

func main() {
	address := flag.String("address", "0.0.0.0", "set listening address (default: 0.0.0.0)")
	port := flag.String("port", "5000", "set listening port (default: 1000)")
	protocol := flag.String("protocol", "TCP", "set listening protocol (default: tcp)")

	flag.Parse()

	manager := database.GetManager()
	defer manager.Close()

	if err := manager.Connect("combolist-db", "postgres", "127.0.0.1", 5433, "docker", "docker", "searchx_combolist"); err != nil {
		log.Println("[!]", err)
	}

	if err := manager.Connect("searchx-db", "postgres", "127.0.0.1", 5435, "docker", "docker", "searchx"); err != nil {
		log.Println("[!]", err)
	}

	ln := socket.Listen(*address, *port, *protocol)

	for {
		socket.SocketAccept(ln)
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
