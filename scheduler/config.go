package scheduler

import "time"

type Config struct {
	Age        int
	Cats       []string
	Pi         float64
	Perfection []int
	DOB        time.Time // requires `import time`
}
