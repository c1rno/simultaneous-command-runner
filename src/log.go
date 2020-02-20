package src

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func NewLogger(cfg *CmdState) *Logger {
	l := &Logger{
		cfg: cfg,
	}

	if cfg == nil {
		l = loggerWithTimeOpts(l)
		l = loggerWithIdentityOpts(l, "root")
	} else {
		if !cfg.WithoutTime {
			l = loggerWithTimeOpts(l)
		}
		if cfg.Name != "" {
			l = loggerWithIdentityOpts(l, cfg.Name)
		}
	}
	return l
}

type Logger struct {
	cfg      *CmdState
	prefixes []func() string
}

// Only for main & watching cmds output
func (l *Logger) Ok(pattern string, params ...interface{}) {
	l.logFn(os.Stdout, pattern, params)
}

// for important internal messages
func (l *Logger) Err(pattern string, params ...interface{}) {
	l.logFn(os.Stderr, pattern, params)
}

// for regular internal messages
func (l *Logger) Debug(pattern string, params ...interface{}) {
	if l.cfg.Debug {
		l.logFn(os.Stderr, pattern, params)
	}
}

func (l *Logger) logFn(wr io.Writer, pattern string, params []interface{}) {
	prefix := ""
	for _, fn := range l.prefixes {
		prefix += fn()
	}
	fmt.Fprintf(
		wr,
		fmt.Sprintf("%s %s\n", prefix, pattern),
		params...,
	)
}

func loggerWithTimeOpts(l *Logger) *Logger {
	l.prefixes = append(l.prefixes, func() string {
		return fmt.Sprintf("[%s]", time.Now().Format(time.RFC3339))
	})
	return l
}

func loggerWithIdentityOpts(l *Logger, identity string) *Logger {
	l.prefixes = append(l.prefixes, func() string {
		return fmt.Sprintf("[%s]", identity)
	})
	return l
}

func Panicer(params ...interface{}) {
	spew.Config.MaxDepth = 10
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableCapacities = true
	spew.Config.SortKeys = true
	tmp := ""
	for _, p := range params {
		tmp += spew.Sdump(p) + " "
	}
	panic(tmp)
}
