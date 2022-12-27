package main

import (
	"fmt"
	"log"
	"net"

	"github.com/burmudar/quiet-hours/models"
	"github.com/burmudar/quiet-hours/server"
)

func main() {
	addr := "127.0.0.1"
	port := server.ListenPort
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	req := models.QuietQuery{
		Whoami: "quietctl",
	}

	data, err := req.Marshal()
	if err != nil {
		log.Fatalf("failed to marshal Quite Query: %v", err)
	}

	log.Printf("[DEBUG] sending: %v", data)
	n, err := conn.Write(data)
	if err != nil {
		log.Fatalf("failed to write marshal data to connection: %v", err)
	}
	log.Printf("[INFO] Sent %d bytes", n)

	var recv [128]byte
	n, err = conn.Read(recv[:])
	if err != nil {
		log.Fatalf("failed to receiving response: %v", err)
	}
	log.Printf("[INFO] Received %d bytes", n)

	// reset size to what we received
	data = recv[:n]

	resp := models.QuietResponse{}

	_, err = resp.Unmarshal(data)
	if err != nil {
		log.Fatalf("failed to unmarshal response to QuietResponse: %v", err)
	}

	log.Printf("[INFO] Response: %#v", resp)
}
