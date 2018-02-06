package cli

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/versus/gethinx/scheduler"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		log.Println("Client got:", string(buf[0:n]))
	}
}

func SocketCli(reload bool, config *scheduler.Config) {
	c, err := net.Dial("unix", config.SocketPath)
	if err != nil {
		log.Fatal("Dial error", err)
	}
	defer c.Close()

	go reader(c)
	if reload {
		_, err := c.Write([]byte("reload"))
		if err != nil {
			log.Fatal("Write error:", err)
		}
		time.Sleep(3 * time.Second)
	}
}
