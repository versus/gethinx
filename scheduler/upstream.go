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

	"sync"

	"sync/atomic"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/looplab/fsm"
	"github.com/versus/gethinx/lib"
)

// Upstream is host for reverseproxy request from ethclients
type Upstream struct {
	Port         uint16     `json:"-"`
	TimeUpdate   int64      `json:"lastupdate"`
	Weight       uint8      `json:"weight"`
	Backup       bool       `json:"-"`
	Host         string     `json:"-"`
	Hostname     string     `json:"hostname"`
	Target       string     `json:"url"`
	Token        string     `json:"-"`
	LastBlock    int64      `json:"digblock"`
	HexLastBlock string     `json:"lastblock"`
	State        string     `json:"-"`
	RealState    string     `json:"status"`
	FSM          *fsm.FSM   `json:"-"`
	Mutex        sync.Mutex `json:"-"`
}

// NewUpstream is constructor for Upstream
func NewUpstream(host string, port string, weight int, token string, hostname string) *Upstream {

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
		Host:       host,
		Target:     target.String(),
		Port:       uint16(uintPort),
		Weight:     uint8(weight),
		Token:      token,
		Hostname:   hostname,
		TimeUpdate: time.Now().Unix(),
	}

	upstream.FSM = fsm.NewFSM(
		"down",
		fsm.Events{
			{Name: "up", Src: []string{"down", "suspend", "backup"}, Dst: "active"},
			{Name: "down", Src: []string{"active", "backup", "suspend"}, Dst: "down"},
			{Name: "backup", Src: []string{"active", "suspend"}, Dst: "backup"},
			{Name: "suspend", Src: []string{"active", "backup"}, Dst: "suspend"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { upstream.enterState(e) },
		},
	)

	return upstream
}

func (u *Upstream) GetTargetLastBlock(ctx context.Context, LastBlock *EthBlock) {

	addr := fmt.Sprintf("http://%s:%d", u.Host, u.Port)
	conn, err := ethclient.Dial(addr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: ", err)
	}

	header, err := conn.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Println("Failed get HeaderByNumber: %v", err)
		u.LastBlock = 0
		if u.FSM.Current() == "up" {
			u.FSM.Event("suspend")
		}
		u.Mutex.Lock()
		u.LastBlock = 0
		u.HexLastBlock = lib.I2H(0)
		u.RealState = u.FSM.Current()
		u.Mutex.Unlock()
		return
	}

	bint := *header.Number

	if bint.IsInt64() {
		u.Mutex.Lock()
		u.LastBlock = bint.Int64()
		u.HexLastBlock = lib.I2H(u.LastBlock)
		if u.LastBlock >= atomic.LoadInt64(&LastBlock.Dig) {
			if u.FSM.Current() == "down" {
				u.FSM.Event("up")
			}
			if u.FSM.Current() == "suspend" {
				u.FSM.Event("up")
			}
		}
		log.Println(u.Target, " is ", u.FSM.Current())
	} else {
		u.LastBlock = 0
		u.HexLastBlock = lib.I2H(0)
		if u.FSM.Current() == "up" {
			if err = u.FSM.Event("suspend"); err != nil {
				log.Println("error change  state FSM to  Down: ", err.Error())
			}
		}
	}
	u.RealState = u.FSM.Current()
	u.TimeUpdate = time.Now().Unix()
	u.Mutex.Unlock()
}

func (u *Upstream) enterState(event *fsm.Event) {
	log.Printf("The upstream %s is %s\n", u.Target, event.Dst)
}

// UpdateLastBlock function for update some fileds in Upstrea: LastBlock value, TimeUpdate value and state to UP
func (u *Upstream) UpdateLastBlock(hexLastBlock string) error {
	var err error
	u.Mutex.Lock()
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
	u.RealState = u.FSM.Current()
	u.Mutex.Unlock()
	return err
}

// GetURL return url.URL from target filed shema://host:port
func (u *Upstream) GetURL() (*url.URL, error) {
	geturl, err := url.Parse(u.Target)
	if err != nil {
		log.Fatalln("Can't get url from ", u.Target, err.Error())
	}
	return geturl, err
}
