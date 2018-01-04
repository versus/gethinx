package scheduler

import (
	"bytes"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/looplab/fsm"
	"github.com/versus/gethinx/lib"
)

type Upstream struct {
	Host         string
	Target       string
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

	_, err = url.Parse(target.String())
	if err != nil {
		log.Fatalln("Can't get url from ", target, err.Error())
	}

	upstream := &Upstream{
		Host:   host,
		Target: target.String(),
		Port:   uint16(uintPort),
		Weight: uint8(uintWeight),
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

func (u *Upstream) GetURL() (*url.URL, error) {
	url, err := url.Parse(u.Target)
	if err != nil {
		log.Fatalln("Can't get url from ", u.Target, err.Error())
	}
	return url, err
}
