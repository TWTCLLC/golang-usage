package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron"
)

var logFile *os.File
var rotateProcessing bool

func StartExe() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("missing args... please start with shell\n")
		fmt.Printf(`eg. ".\\logRotation duration exe_path exe_args..."`)
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		return
	}
	duration := args[1]

	_, err := os.Stat("./log")
	if err != nil {
		os.Mkdir("./log", 0666)
	}

	_, err = time.ParseDuration(duration)
	if err != nil {
		fmt.Printf("invalid duration: %v\n", err)
		return
	}

	crontab := cron.New()
	crontab.AddFunc("@every "+duration, rotate)
	crontab.Start()
	rotate()

	cmd := exec.Command(args[2], args[3:]...)
	op, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Stdout err: %v\n", err)
		return
	}
	s := bufio.NewScanner(op)
	cmd.Start()

	for {
		if s.Scan() || rotateProcessing {
			continue
		}
		fmt.Printf(s.Text() + "\n")
		logFile.WriteString(s.Text() + "\n")
	}
}

func rotate() {
	rotateProcessing = true
	logFile.Close()
	fileName := time.Now().Format("0601020304")
	cmd := exec.Command("powershell", "new-item", ".\\log\\"+fileName+".log")
	cmd.Start()
	cmd.Wait()
	var err error
	logFile, err = os.OpenFile("./log/"+fileName+".log", os.O_APPEND, 0660)
	if err != nil {
		fmt.Printf("open file err: %v\n", err)
		return
	}
	rotateProcessing = false
}
