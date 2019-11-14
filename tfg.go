package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"./filebeat"
	"./others"
)

/*
	Input:
	Output:
	Execution:
*/

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
			others.CheckError(errors.New("	Error in the parsing of NODES_NUMBER"))
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
			others.CheckError(errors.New("	Error in the parsing of NODE_TIMEOUT"))
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
	others.CheckError(commandExecution.Process.Kill()) //Check if there is any error
	fmt.Println(fmt.Sprintf("ES Node %d executed with a duration of %f", idNode, executionTime))
}

func makeCommandES(config *map[string]string) string {
	/*
		Input: The configuration map.
		Output: The string defined by the command to execute for start a node.
		Execution:	It checks the ELASTICSEARCH_PATH and looks if it defined and concats the directory of it with the binary location
	*/
	if strings.Compare((*config)["ELASTICSEARCH_PATH"], "") == 0 {
		others.CheckError(errors.New("The ELASTICSEARCH_PATH is not defined"))
	}
	command := others.ObtainDirectory((*config)["ELASTICSEARCH_PATH"]) + "/bin/elasticsearch"
	return command
}

//
func main() {
	configuration := others.ObtainConfiguration() //Obtains a reference of a map with the configurations
	//wireshark.ExecuteTshark(configuration)        //Executes tshark
	filebeat.FilebeatExecution(configuration)
	//desployElasticSearchNetwork(configuration) //Start the nodes execution

	//executeKibana(args)
}
