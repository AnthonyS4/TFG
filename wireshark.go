package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
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

//This function will check if there is any error in the reading
// of the configuration file, if it exists then it invokes panic(error)
func checkData(e error) {
	if e != nil {
		panic(e)
	}
}

func removeElement(vector []string, index int) []string {
	/*
		Input: An array of strings and the index of the element to be erased
		Output: An array of the string that will have the keys and values of the map
		Execution: It checks the existence of the substring ":" in the lines of the file, and erases the strings that don't have one
	*/
	vector[index] = vector[len(vector)-1]
	return vector[:len(vector)-1]
}

func checkAttributes(data []byte) []string {
	/*
		Input: The array of bytes read in the config.yml file
		Output: An array of the string that will have the keys and values of the map
		Execution: It checks the existence of the substring ":" in the lines of the file, and erases the strings that don't have one
	*/
	configurationStrings := strings.Split(string(data), "\n")
	for i := 0; i < len(configurationStrings); i++ {
		if strings.Count(configurationStrings[i], ":") != 1 {
			configurationStrings = removeElement(configurationStrings, i)
		}
	}
	return configurationStrings[:len(configurationStrings)-1]
	//It returns the parametters except the last one because it is a line that contains the end of file
}

func obtainConfiguration(configuration *map[string]string) {
	/*
		Input: The reference of the map that will contains the configurations of the program
		Output: ~
		Execution: This function reads the file config.yml and make a map with the attibutes defined in the file
	*/
	data, error := ioutil.ReadFile("config.yml")
	checkData(error)
	lines := checkAttributes(data)
	fmt.Println(len(lines))
	for i := 0; i < len(lines); i++ { //This loop will read and store the keys and values defined in the config.yml
		fmt.Println(lines[i])
		key := strings.Split(lines[i], ":")[0]
		value := strings.Split(lines[i], ":")[1]
		(*configuration)[key] = value
	}
}

//
func main() {
	//	defineAskpass()
	//	var processTshark *exec.Cmd
	//	processTshark = executeTshark()
	//	waitEndProcess(processTshark)
	configuration := make(map[string]string)
	obtainConfiguration(&configuration)
	executeTshark(obtainConfiguration)
}
