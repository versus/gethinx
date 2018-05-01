package gethinx

import (
	"github.com/versus/gethinx/config"
	"io"
	"log"
	"net"
	"time"
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

func SocketCli(reload bool, config *config.Config) {
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
