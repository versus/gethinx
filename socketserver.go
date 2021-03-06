package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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

func StartSocketServer() {

	syscall.Unlink(conf.SocketPath)
	ln, err := net.Listen("unix", conf.SocketPath)
	if err != nil {
		log.Fatal("Listen error: ", conf.SocketPath, err.Error())
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		ln.Close()
		os.Exit(0)
	}(ln, sigc)

	for {
		fd, err := ln.Accept()
		if err == nil {
			go RequestSocketServer(fd)
		}
	}

}
