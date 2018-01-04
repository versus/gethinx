package scheduler

//import "time"
import (
	"github.com/looplab/fsm"
)

type Upstream  struct {
	Host string
	HexLastBlock string
	LastBlock int64
	Timeupdate int
	Weight uint
    Backup bool
	State  string
	FSM *fsm.FSM
	}



//log.Println(time.Now().Unix())