package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	app "github.com/c1rno/simultaneous-command-runner/src"
	_ "go.uber.org/automaxprocs"
	"gopkg.in/yaml.v2"
)

func main() {
	rootLogger := app.NewLogger(nil)
	rootLogger.Err("Started")
	defer rootLogger.Err("Completed")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func() {
		<-ctx.Done()
		rootLogger.Err(`
Example:

scr << EOF
timeout: 100
wait_all: false
run:
    cmd: "echo 🦁🐯"
    debug: true
watch:
    - cmd: "tail -F /dev/null"
      debug: false
      disable_time: true
      name: useless
    - cmd: "kubectl logs -f -lname=my-app"
      debug: true
      name: k8s
    - cmd: "ls -lah /var/log/"
      debug: true
      name: ls
EOF
`)
		os.Exit(1)
	}()
	raw, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		app.Panicer(`Read stdin failed`, err)
	}

	state := &app.RootState{}
	err = yaml.Unmarshal(raw, state)
	if err != nil {
		app.Panicer(`Parse yaml failed`, err)
	}

	ctx = context.Background()
	if state.TimeOut != 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Second*state.TimeOut)
		defer cancel()
	}
	go func() {
		<-ctx.Done()
		app.Panicer(`Timeout!`, state.TimeOut.Seconds())
	}()

	state.StartSignal = make(chan interface{})

	state.MainCmd.StopSignal = &sync.WaitGroup{}
	state.MainCmd.StopSignal.Add(1)
	if err = app.PrepareExecuteCmd(ctx, state, &state.MainCmd); err != nil {
		app.Panicer(fmt.Sprintf(`Preparing "%s" failed`, state.MainCmd.Cmd), err)
	}
	for _, watchingCmd := range state.WatchingCmds {
		watchingCmd.StopSignal = &sync.WaitGroup{}
		watchingCmd.StopSignal.Add(1)
		if err = app.PrepareExecuteCmd(ctx, state, &watchingCmd); err != nil {
			app.Panicer(fmt.Sprintf(`Preparing "%s" failed`, watchingCmd.Cmd), err)
		}
	}

	close(state.StartSignal) // start all cmd

	state.MainCmd.StopSignal.Wait()
	if state.WaitAll {
		for _, watchingCmd := range state.WatchingCmds {
			watchingCmd.StopSignal.Wait()
		}
	}
}
