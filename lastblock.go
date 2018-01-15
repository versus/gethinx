package main

import "github.com/versus/gethinx/lib"

//GenerateLastBlockAverage is func for generate average ariphemic number block
func GenerateLastBlockAverage() {
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
	LastBlock.Mutex.Lock()
	LastBlock.Dig = average
	LastBlock.Hex = lib.I2H(average)
	LastBlock.Mutex.Unlock()
}
