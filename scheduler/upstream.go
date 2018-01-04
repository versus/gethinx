package scheduler

import (
	"github.com/looplab/fsm"
	"log"
	"github.com/versus/gethinx/lib"
	"time"
)

type Upstream  struct {
	Host string
	HexLastBlock string
	LastBlock int64
	Timeupdate int64
	Weight uint
    Backup bool
	State  string
	FSM *fsm.FSM
}

func NewUpstream(to string) *Upstream {
	u := &Upstream{
		State: to,
	}

	u.FSM = fsm.NewFSM(
		"down",
		fsm.Events{
			{Name: "up", Src: []string{"down"}, Dst: "active"},
			{Name: "down", Src: []string{"active"}, Dst: "down"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { u.enterState(e) },
		},
	)

	return u
}

func (u *Upstream) enterState(event *fsm.Event) {
	log.Printf("The upstream to %s is %s\n", u.State, event.Dst)
}

func (u *Upstream) UpdateLastBlock(host string, hexlastblock string) (error) {
	var err error
	u.Host = host
	u.LastBlock, err = lib.H2I(hexlastblock)
	if err != nil {
		log.Fatalln("Error convert block to int", err.Error())
		return err
	}
	u.HexLastBlock = hexlastblock
	u.Timeupdate = time.Now().Unix()
	if u.FSM.Current() == "down" {
		u.FSM.Event("up")
	}
	return err
}

