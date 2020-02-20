package src

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func PrepareExecuteCmd(
	ctx context.Context,
	globalState *RootState,
	localState *CmdState,
) error {

	params := strings.Split(localState.Cmd, " ")
	cmd := exec.CommandContext(ctx, params[0], params[1:]...)
	cmd.Stdout = os.Stdout // only scr write to stdout
	cmd.Stderr = os.Stdout

	go func() {
		defer localState.StopSignal.Done()
		log := NewLogger(localState)
		<-globalState.StartSignal

		if err := cmd.Run(); err != nil {
			if localState.Cmd == globalState.MainCmd.Cmd {
				Panicer(fmt.Sprintf(`Main command failed (%s)`, localState.Cmd), err)
			} else {
				log.Err("%v", err)
			}
		}
	}()

	return nil
}
