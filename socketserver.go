package main

import (
	"log"
	"net"
)

//https://gist.github.com/hakobe/6f70d69b8c5243117787fd488ae7fbf2
func RequestSocketServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		log.Println("Server got:", string(data))
		if string(data) == "reload" {
			ReloadBackendServers(flagConfigFile)
			data = []byte("server was reloaded")
		}
		_, err = c.Write(data)
		if err != nil {
			log.Println("Writing client error: ", err)
		}
	}
}

func StartSocketServer(ln net.Listener) {

	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Println("Accept error: ", err)
		}

		go RequestSocketServer(fd)
	}

}
