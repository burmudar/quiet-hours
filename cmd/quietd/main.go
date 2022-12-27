package main

import (
	"github.com/burmudar/quiet-hours/server"
)

func main() {
	s := server.New("127.0.0.1", server.ListenPort)

	server.QuietHours[8] = 0
	server.QuietHours[9] = 0
	server.QuietHours[10] = 0
	server.QuietHours[11] = 0
	server.QuietHours[12] = 0
	server.QuietHours[13] = 0
	s.Listen()
}
