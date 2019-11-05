package main

import (
	"fmt"
	"os/exec"
	"time"
)

//With this function we define the enviorement variable SUDO_ASKPASS, sudo command will use this variable.
func defineAskpass() {
	var askpassCommand string
	askpassCommand = "SUDO_ASKPASS=\"/home/anthony/go/scr/wireshark/askpass\""
	exec.Command(askpassCommand)
}

//With this function we start the execution of tshark with elevate privilegies
func executeTshark() *exec.Cmd {
	var tsharkCommand string
	tsharkCommand = "${SUDO_ASKPASS} | sudo -S tshark -i wlo1 -T ek > paquetes.json"
	var commandExecutionTshark *exec.Cmd
	commandExecutionTshark = exec.Command(tsharkCommand)
	return commandExecutionTshark
}

//With this function we give 3 seconds of data recollection, after this it ends the process
func waitEndProcess(tshark *exec.Cmd) {
	var start time.Time
	var timeOut float64
	timeOut = 3.0
	start = time.Now()
	for time.Since(start).Seconds() < timeOut {
	}
	if tshark == nil {
		fmt.Println("Nil pointer")
	} else {
		tshark.Process.Kill()
		fmt.Println("Killed process")
	}
}

//
func main() {
	defineAskpass()
	var processTshark *exec.Cmd
	processTshark = executeTshark()
	waitEndProcess(processTshark)
}
