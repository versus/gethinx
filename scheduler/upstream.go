package scheduler

import (
	"bytes"
	"log"
	"net/url"
	"strconv"
	"time"

	"strings"

	"fmt"

	"context"

	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/looplab/fsm"
	"github.com/versus/gethinx/lib"
)

// Upstream is host for reverseproxy request from ethclients
type Upstream struct {
	Port         uint16
	LastBlock    big.Int
	TimeUpdate   int64
	Weight       uint8
	Backup       bool
	Host         string
	Target       string
	HexLastBlock string
	State        string
	FSM          *fsm.FSM
}

// NewUpstream is constructor for Upstream
func NewUpstream(host string, port string, weight int) *Upstream {

	uintPort, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		log.Fatalln("Can't convert port to uint16", err.Error())
	}

	if weight == 0 {
		weight = 1
	}

	target := bytes.NewBufferString("")
	if !strings.Contains(host, "http://") {
		if _, err = target.WriteString("http://"); err != nil {
			log.Fatalln("Error in construction target ", err.Error())
		}
	}
	if _, err = target.WriteString(host); err != nil {
		log.Fatalln("Error in construction target ", err.Error())
	}
	if _, err = target.WriteString(":"); err != nil {
		log.Fatalln("Error in construction target ", err.Error())
	}
	if _, err = target.WriteString(port); err != nil {
		log.Fatalln("Error in construction target ", err.Error())
	}

	_, err = url.Parse(target.String())
	if err != nil {
		log.Fatalln("Can't get url from ", target, err.Error())
	}

	upstream := &Upstream{
		Host:   host,
		Target: target.String(),
		Port:   uint16(uintPort),
		Weight: uint8(weight),
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

func (u *Upstream) GetLastBlock(ctx context.Context) {

	addr := fmt.Sprintf("http://%s:%d", u.Host, u.Port)
	log.Println("addr is ", addr)
	conn, err := ethclient.Dial(addr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	header, err := conn.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("Failed get HeaderByNumber: %v", err)
	}

	u.LastBlock = *header.Number
}

func (u *Upstream) enterState(event *fsm.Event) {
	log.Printf("The upstream to %s is %s\n", u.State, event.Dst)
}

// UpdateLastBlock function for update some fileds in Upstrea: LastBlock value, TimeUpdate value and state to UP
func (u *Upstream) UpdateLastBlock(hexLastBlock string) error {
	var err error
	u.LastBlock, err = lib.H2I(hexLastBlock)
	if err != nil {
		log.Fatalln("Error convert block to int", err.Error())
		return err
	}
	u.HexLastBlock = hexLastBlock
	u.TimeUpdate = time.Now().Unix()
	if u.FSM.Current() == "down" {
		if err = u.FSM.Event("up"); err != nil {
			log.Fatalln("error change  state FSM to  UP: ", err.Error())
		}
	}
	return err
}

// GetURL return url.URL from target filed shema://host:port
func (u *Upstream) GetURL() (*url.URL, error) {
	url, err := url.Parse(u.Target)
	if err != nil {
		log.Fatalln("Can't get url from ", u.Target, err.Error())
	}
	return url, err
}
