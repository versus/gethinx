package proxyserver

//GenerateLastBlockAverage is func for generate average ariphemic number block
func GenerateLastBlockAverage() {
	/*
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
		LastBlock.Hex = buffer.I2H(average)
		LastBlock.TimeUpdate = time.Now().Unix()
		monitoring.PromLastBlock.Set(float64(LastBlock.Dig))
		LastBlock.Mutex.Unlock()
	*/
}
