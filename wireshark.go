package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

/*
	Input:
	Output:
	Execution:
*/

func checkData(e error) {
	/*
		Input: The error obtained after the read of the file
		Output: ~
		Execution: It checks if the input is nil, if it's not then it calls panic(input)
	*/
	if e != nil {
		panic(e)
	}
}

func removeElement(vector []string, index int) []string {
	/*
		Input: An array of strings and the index of the element to be erased
		Output: An array of the string that will have the keys and values of the map
		Execution: Uses append between the elements before index and the following elements.
	*/
	return append(vector[:index], vector[index+1:]...)
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
	for i := 0; i < len(lines); i++ { //This loop will read and store the keys and values defined in the config.yml
		keyAndValue := strings.Split(lines[i], ":")
		(*configuration)[keyAndValue[0]] = keyAndValue[1]
	}
}

//With this function we start the execution of tshark with elevate privilegies
func executeTshark(config *map[string]string) {
	tsharkCommand := makeCommand(config)
	fmt.Println(tsharkCommand)
	//var commandExecutionTshark *exec.Cmd
	//commandExecutionTshark = exec.Command(tsharkCommand)
}

func makeCommand(config *map[string]string) string {
	command := "sudo -S tshark -i " + (*config)["NETWORK_INTERFACE"] + " -T ek > " + obtainDirectory((*config)["PACKETS_OUTPUT_DIRNAME"])
	if strings.Compare((*config)["PACKETS_OUTPUT_FILENAME"], "") == 0 {
		//This attribute is not defined, then we use the predefined name
		command += "packets.json"
	} else {
		//Add the filename
		command += (*config)["PACKETS_OUTPUT_FILENAME"] + ".json"
	}
	return command
}

func obtainDirectory(path string) string {
	variable := ""
	if strings.Compare(path, "") == 0 { //Output directory name not defined, then we return ""
		return variable
	} else { //Output directory is defined, check the existence of env. variables
		if strings.Contains(path, "${") && strings.Contains(path, "}") { //Lexical check
			variable += lookupEnvVariable(path) + strings.SplitAfter(path, "}")[1] + "/"
		} else { //There is no env. variables, then we return the defined directory
			variable = path + "/"
		}
		return variable
	}
}

func lookupEnvVariable(path string) string {
	variable := strings.Split(path, "}")[0]
	variable = strings.Split(variable, "{")[1]
	variableDirectory := os.Getenv(variable)
	if strings.Compare(variableDirectory, "") == 0 {
		//The env. variable doesn't exists, then we use the Homa variable
		variableDirectory = os.Getenv("HOME")
	} //Else then we return the new path
	return variableDirectory
}

//
func main() {
	configuration := make(map[string]string)
	obtainConfiguration(&configuration)
	executeTshark(&configuration)
}
