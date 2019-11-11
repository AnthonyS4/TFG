package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

func obtainConfiguration() *map[string]string {
	/*
		Input: ~
		Output: The reference of the map that will contains the configurations of the program
		Execution: This function reads the file config.yml and make a map with the attibutes defined in the file
	*/
	configuration := make(map[string]string)
	data, error := ioutil.ReadFile("config.yml")
	checkError(error) //Checking if the file exists
	lines := checkAttributes(data)
	for i := 0; i < len(lines); i++ { //This loop will read and store the keys and values defined in the config.yml
		keyAndValue := strings.Split(lines[i], ":")
		configuration[keyAndValue[0]] = keyAndValue[1]
	}
	return &configuration
}

func executeTshark(config *map[string]string) {
	/*
		Input: The configuration map
		Output: ~
		Execution: This function executes the tshark command for an measure of time given in the config.yml
	*/
	commandExecutionTshark, executionTimeLimit, writer := startTshark(makeCommandTshark(config), config) //startTshark will begin the process
	go stopTshark(commandExecutionTshark, executionTimeLimit, writer)
	if bytes, e := commandExecutionTshark.Output(); e == nil {
		//Convert the []bytes to JSON
		halfs := strings.Split(commandExecutionTshark.Args[2], " > ")
		ioutil.WriteFile(halfs[1], bytes, 0)
	} else { //An error ocurred
		checkError(e)
	}
	fmt.Println(fmt.Sprintf("Tshark executed with a duration of %f", executionTimeLimit))
}

func stopTshark(commandExecutionTshark *exec.Cmd, executionTimeLimit float64, writer *io.WriteCloser) {
	now := time.Now()
	for time.Now().Sub(now).Seconds() < executionTimeLimit { //Loop for execution time control
	}
	checkError(commandExecutionTshark.Process.Kill())
}

func startTshark(tsharkCommand string, config *map[string]string) (*exec.Cmd, float64, *io.WriteCloser) {
	/*
		Input: The command string and the map of configurations.
		Output:	The reference to the Cmd that started the tshark command and the duration of the execution.
		Execution: Using a Reader for read the password in the std input of the process, that is who sudo can continue. It starts the process and returns it with the measure of time
	*/
	//buffer := strings.NewReader((*config)["PASSWORD"])
	tsharkCommand = "sudo ls -l > /home/anthony/TFG/archivos/ls.txt"
	commandExecutionTshark := exec.Command("bash", "-c", tsharkCommand)
	//commandExecutionTshark.Stdin = buffer //With this the process can obtain the value of the key "PASSWORD"
	writerPipe, e := commandExecutionTshark.StdinPipe()
	checkError(e)
	executionTimeLimit := getDuration(config)
	executionTimeLimit = 2.0
	writerPipe.Write([]byte((*config)["PASSWORD"]))
	//commandExecutionTshark.Start()
	writerPipe.Close()
	return commandExecutionTshark, executionTimeLimit, &writerPipe
}

func getDuration(config *map[string]string) float64 {
	/*
		Input: The configuration map
		Output:	The amount of seconds defined in config.yml, if it is not defined => we use the default duration(5 seconds)
		Execution: We parse the keys "TIME_SECONDS" and "TIME_MINUTES", if these are defined => then it stores the add in seconds of both in the variable duration
	*/
	duration := 5.0
	if seconds, error := strconv.ParseFloat((*config)["TIME_SECONDS"], 64); error == nil && seconds > 0.0 {
		//The error is nil => There is a amount of seconds defined, we have to check if it's negative
		duration = seconds //Modify the duration
	} else {
		checkError(errors.New("Error in the parsing of TIME_SECONDS"))
	}
	if minutes, error := strconv.ParseFloat((*config)["TIME_MINUTES"], 64); error == nil && minutes > 0.0 {
		//The error is nil => There is a amount of minutes defined, we have to check if it's negative
		duration += minutes * 60 //Modify the duration
	} else {
		checkError(errors.New("Error in the parsing of TIME_MINUTES"))
	}
	return duration
}

func makeCommandTshark(config *map[string]string) string {
	/*
		Input: The configuration map
		Output:	The string of the command for the execution of tshark with the args defined in config.yml
		Execution:	It checks the map and concats the  begin of command with the output_dirname/filename
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
		//The env. variable doesn't exists, then stop the execution
		checkError(errors.New("The env. variable " + variable + " is not defined"))
	} //Else then we return the new path
	return variableDirectory
}

func desployElasticSearchNetwork(config *map[string]string) {
	/*
		Input: The configuration map.
		Output: ~
		Execution: This function start a ${nodesNumber} goroutines with a shared reference of a waitGroup, every goroutine will start the execution of a elasticsearch node. It will wait the goroutines's end.
	*/
	nodesNumber := obtainNodesNumber(config)
	var waitGroup sync.WaitGroup
	for i := 0; i < nodesNumber; i++ {
		waitGroup.Add(1)
		go startESNode(makeCommandES(config), config, i+1, &waitGroup)
	}
	waitGroup.Wait() //It waits until the execution of all the goroutines => End of the nodes net
}

func obtainNodesNumber(config *map[string]string) int {
	/*
		Input: The configuration map.
		Output: The nodes numbers of the elasticsearch network.
		Execution: It checks if NODES_NUMBER is defined, if it is then it checks the validity of the parsed number.
	*/

	if strings.Compare((*config)["NODES_NUMBER"], "") == 0 {
		return 1
	} else {
		nodesNumber, e := strconv.ParseInt((*config)["NODES_NUMBER"], 10, 32)
		if nodesNumber < 1 || e != nil {
			checkError(errors.New("	Error in the parsing of NODES_NUMBER"))
		}
		return int(nodesNumber)
	}
}

func getNodeTimeOut(config *map[string]string) float64 {
	/*
		Input: The configuration map.
		Output: The numbers of seconds that will be the node timeout.
		Execution: It checks if NODES_TIMEOUT is defined, if it is then it checks the validity of the parsed number.
	*/

	if strings.Compare((*config)["NODE_TIMEOUT"], "") == 0 {
		return 1.0
	} else {
		exeTime, e := strconv.ParseFloat((*config)["NODE_TIMEOUT"], 64)
		if exeTime < 0.0 || e != nil {
			checkError(errors.New("	Error in the parsing of NODE_TIMEOUT"))
		}
		return exeTime
	}
}

func startESNode(command string, config *map[string]string, idNode int, waitGroup *sync.WaitGroup) {
	/*
		Input:The command to execute, configuration map, the identifier of this goroutine, the "barrier" waitgroup.
		Output: ~
		Execution: It start the process with the input command, it wait a timeControl, then it executes the waitGroup.Done().
	*/
	defer (*waitGroup).Done() //This instruction will execute after the end of this function
	commandExecutionES := exec.Command("bash", "-c", command)
	commandExecutionES.Start()
	timeControl(commandExecutionES, config, idNode)
}

func timeControl(commandExecution *exec.Cmd, config *map[string]string, idNode int) {
	/*
		Input: The reference to the Cmd that started a node, the configuration, and the identifier of the actual node
		Output: ~
		Execution: It waits the executionTime using a loop, then it kills the node process, and prints the end.
	*/
	now := time.Now()
	executionTime := getNodeTimeOut(config)
	for time.Now().Sub(now).Seconds() < executionTime { //Loop for execution time control
	}
	checkError(commandExecution.Process.Kill()) //Check if there is any error
	fmt.Println(fmt.Sprintf("ES Node %d executed with a duration of %f", idNode, executionTime))
}

func makeCommandES(config *map[string]string) string {
	/*
		Input: The configuration map.
		Output: The string defined by the command to execute for start a node.
		Execution:	It checks the ELASTICSEARCH_PATH and looks if it defined and concats the directory of it with the binary location
	*/
	if strings.Compare((*config)["ELASTICSEARCH_PATH"], "") == 0 {
		checkError(errors.New("The ELASTICSEARCH_PATH is not defined"))
	}
	command := obtainDirectory((*config)["ELASTICSEARCH_PATH"]) + "/bin/elasticsearch"
	return command
}

func filebeatExecution(config *map[string]string) {
	//commandExecutionFilebeat := exec.Command("bash", "-c", makeCommandFilebeat(config))
	fmt.Println(makeCommandFilebeat(config))
}

func makeCommandFilebeat(config *map[string]string) string {
	if strings.Compare((*config)["FILEBEAT_PATH"], "") == 0 {
		checkError(errors.New("The FILEBEAT_PATH is not defined"))
	}
	command := obtainDirectory((*config)["FILEBEAT_PATH"]) + "filebeat -e"
	return command
}

//
func main() {
	configuration := obtainConfiguration() //Obtains a reference of a map with the configurations
	executeTshark(configuration)           //Executes tshark
	//filebeatExecution(configuration)
	//desployElasticSearchNetwork(configuration) //Start the nodes execution

	//executeKibana(args)
}
