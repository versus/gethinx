package scheduler

import "github.com/versus/gethinx/lib"

//GenerateLastBlockAverage is func for generate average ariphemic number block
func GenerateLastBlockAverage(backends map[string]Upstream, lastBlock *EthBlock) {
	sum, count := int64(0), int64(0)
	average := int64(0)

	for _, srv := range backends {
		if srv.FSM.Current() == "active" {
			sum = sum + srv.LastBlock
			count++
		}
	}
	if count != 0 {
		average = int64(sum / count)
	}
	lastBlock.Mutex.Lock()
	lastBlock.Dig = average
	lastBlock.Hex = lib.I2H(average)
	lastBlock.Mutex.Unlock()
}
