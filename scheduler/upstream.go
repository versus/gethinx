package scheduler

import (
	"bytes"
	"github.com/looplab/fsm"
	"github.com/versus/gethinx/lib"
	"log"
	"net/url"
	"strconv"
	"time"
)

type Upstream struct {
	Host         string
	Url          *url.URL
	Port         uint16
	HexLastBlock string
	LastBlock    int64
	TimeUpdate   int64
	Weight       uint8
	Backup       bool
	State        string
	FSM          *fsm.FSM
}

func NewUpstream(host string, port string, weight string) *Upstream {

	uintPort, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		log.Fatalln("Can't convert port to uint16", err.Error())
	}

	uintWeight, err := strconv.ParseUint(weight, 10, 8)
	if err != nil {
		log.Fatalln("Can't convert weight to uint8", err.Error())
	}
	if uintWeight == 0 {
		uintWeight = 1
	}

	target := bytes.NewBufferString("http://")
	target.WriteString(host)
	target.WriteString(":")
	target.WriteString(port)

	url, err := url.Parse(host)
	if err != nil {
		log.Fatalln("Can't get url from ", target, err.Error())
	}

	upstream := &Upstream{
		Host:   host,
		Port:   uint16(uintPort),
		Weight: uint8(uintWeight),
		Url:    url,
	}

	upstream.FSM = fsm.NewFSM(
		"down",
		fsm.Events{
			{Name: "up", Src: []string{"down"}, Dst: "active"},
			{Name: "down", Src: []string{"active"}, Dst: "down"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { upstream.enterState(e) },
		},
	)

	return upstream
}

func (u *Upstream) enterState(event *fsm.Event) {
	log.Printf("The upstream to %s is %s\n", u.State, event.Dst)
}

func (u *Upstream) UpdateLastBlock(host string, hexLastBlock string) error {
	var err error
	u.Host = host
	u.LastBlock, err = lib.H2I(hexLastBlock)
	if err != nil {
		log.Fatalln("Error convert block to int", err.Error())
		return err
	}
	u.HexLastBlock = hexLastBlock
	u.TimeUpdate = time.Now().Unix()
	if u.FSM.Current() == "down" {
		u.FSM.Event("up")
	}
	return err
}

func getUrl(req lib.MyRequestBody) (*url.URL, error) {
	_ = req.Body
	target := "http://127.0.0.1:8080"
	url, err := url.Parse(target)
	return url, err
}
