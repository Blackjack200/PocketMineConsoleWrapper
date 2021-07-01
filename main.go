package main

import (
	"github.com/chzyer/readline"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"unsafe"
)

func input(r rune) (rune, bool) {
	switch r {
	case readline.CharInterrupt:
		stdin.Write([]byte("stop\n"))
		stdin.Write([]byte("end\n"))
		stdin.Write([]byte("shutdown\n"))
		go func() {
			time.Sleep(time.Second * 30)
			proc.Process.Kill()
		}()
		return r, false
	}
	return r, true
}

var proc *exec.Cmd
var l *readline.Instance
var stdin io.WriteCloser

func main() {
	initReadline()
	initProc()
	p, _ := proc.StdoutPipe()
	go func() {
		buf := make([]byte, 512)
		for {
			le, err := p.Read(buf)
			if err != nil {
				continue
			}
			l.Write(buf[:le])
		}
	}()
	if err := proc.Start(); err != nil {
		log.Panic(err)
	}
	go func() {
		proc.Wait()
		s := !proc.ProcessState.Success()
		os.Exit(*(*int)(unsafe.Pointer(&s)))
	}()
	for {
		line, err := l.Readline()
		if err != nil {
			continue
		}
		stdin.Write([]byte(line + "\n"))
	}
}

func initProc() {
	WorkingPath, _ := os.Getwd()
	WorkingPath, _ = filepath.Abs(WorkingPath)
	proc = exec.Command("bash", WorkingPath+"/wrapper.sh")
	proc.Env = os.Environ()
	var err error
	stdin, err = proc.StdinPipe()
	if err != nil {
		log.Panic(err)
	}
}

func initReadline() {
	var err error
	l, err = readline.NewEx(&readline.Config{
		Prompt:      "> ",
		HistoryFile: "/tmp/pocketmine_console_wrapper_readline.tmp",

		HistorySearchFold:   true,
		FuncFilterInputRune: input,
	})
	if err != nil {
		log.Panic(err)
	}
}
