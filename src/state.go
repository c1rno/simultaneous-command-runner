package src

import (
	"sync"
	"time"
)

type RootState struct {
	StartSignal  chan interface{}
	MainCmd      CmdState      `yaml:"run"`
	WatchingCmds []CmdState    `yaml:"watch"`
	TimeOut      time.Duration `yaml:"timeout"`
	WaitAll      bool          `yaml:"wait_all"`
}

type CmdState struct {
	Cmd         string `yaml:"cmd"`
	Debug       bool   `yaml:"debug"`
	WithoutTime bool   `yaml:"disable_time"`
	Name        string `yaml:"name"`
	StopSignal  *sync.WaitGroup
}
