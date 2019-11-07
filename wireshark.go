package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//With this function we define the enviorement variable SUDO_ASKPASS, sudo command will use this variable.
/*func defineAskpass() {
	var askpassCommand string
	askpassCommand = "SUDO_ASKPASS=\"/home/anthony/go/scr/wireshark/askpass\""
	exec.Command(askpassCommand)
}*/

//With this function we give 3 seconds of data recollection, after this it ends the process
/*func waitEndProcess(tshark *exec.Cmd) {
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
}*/

/*
	Input:
	Output:
	Execution:
*/

func checkError(e error) {
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
			i--
		}
	}
	return configurationStrings
	//It returns the parametters except the last one because it is a line that contains the end of file
}

func obtainConfiguration(configuration *map[string]string) {
	/*
		Input: The reference of the map that will contains the configurations of the program
		Output: ~
		Execution: This function reads the file config.yml and make a map with the attibutes defined in the file
	*/
	data, error := ioutil.ReadFile("config.yml")
	checkError(error)
	lines := checkAttributes(data)
	for i := 0; i < len(lines); i++ { //This loop will read and store the keys and values defined in the config.yml
		keyAndValue := strings.Split(lines[i], ":")
		(*configuration)[keyAndValue[0]] = keyAndValue[1]
	}
}

//With this function we start the execution of tshark with elevate privilegies
func executeTshark(config *map[string]string) {
	tsharkCommand := makeCommand(config)
	var buffer io.Reader
	buffer = strings.NewReader((*config)["PASSWORD"])
	commandExecutionTshark := exec.Command("bash", "-c", tsharkCommand)
	commandExecutionTshark.Stdin = buffer
	executionTimeLimit := getDuration(config)
	commandExecutionTshark.Start()
	now := time.Now()
	for time.Now().Sub(now).Seconds() < executionTimeLimit {
	}
	if err := commandExecutionTshark.Process.Kill(); err != nil {
		log.Fatal("failed to kill process: ", err)
	}
	fmt.Println(fmt.Sprintf("Executed with a duration of %f", executionTimeLimit))
}

func startTshark(command string, config *map[string]string) *exec.Cmd {

}

func getDuration(config *map[string]string) float64 {
	duration := 5.0
	if seconds, error := strconv.ParseFloat((*config)["TIME_SECONDS"], 64); error == nil && seconds > 0.0 {
		//The error is nil => There is a amount of seconds defined, we have to check if it's negative
		duration = seconds //Modify the duration
	}
	if minutes, error := strconv.ParseFloat((*config)["TIME_MINUTES"], 64); error == nil && minutes > 0.0 {
		//The error is nil => There is a amount of minutes defined, we have to check if it's negative
		duration += minutes * 60 //Modify the duration
	}
	return duration
}

func makeCommand(config *map[string]string) string {
	/*
		Input: The configuration map
		Output:	The string of the command for the execution of tshark with the args defined in config.yml
		Execution:	It checks the map and concats the output dirname/filename
	*/
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
	/*
		Input: The path given in the key(PACKETS_OUTPUT_DIRNAME) of the config map
		Output:	The absolute path that packets will stored
		Execution: It checks the lexically if there is any "${" => there is a env. variable then it gets the absolute path of it, else it returns the absolute path given or ""
	*/
	variable := ""
	if strings.Compare(path, "") != 0 {
		if strings.Contains(path, "${") && strings.Contains(path, "}") { //Lexical check
			variable += lookupEnvVariable(path) + strings.SplitAfter(path, "}")[1] + "/"
		} else { //There is no env. variables, then we return the defined directory
			variable = path + "/"
		}
	}
	return variable
}

func lookupEnvVariable(path string) string {
	/*
		Input: The path with the env. variable
		Output: The absolute path defined in the env. variable
		Execution: It separates the env. variable taking out the "{" and "}", then it does an request using Getenv and gets it.
	*/
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
