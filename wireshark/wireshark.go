package wireshark

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"../others"
)

// Execute :This function executes the tshark command for an measure of time given in the config.yml
func Execute(config *map[string]string, reference *string) {
	/*
		Input: The configuration map
		Output: ~
	*/
	(*reference) = makeCommandTshark(config)
	startTshark((*reference), config) //startTshark will begin the process
	fmt.Println("Tshark executed!")
}

func startTshark(tsharkCommand string, config *map[string]string) {
	/*
		Input: The command string and the map of configurations.
		Output:	The reference to the Cmd that started the tshark command and the duration of the execution.
		Execution: Using a Reader for read the password in the std input of the process, that is who sudo can continue. It starts the process and returns it with the measure of time
	*/
	sudoCommandExecution := exec.Command("bash", "-c", "sudo echo Capture of packets")
	sudoCommandExecution.Run()
	tsharkExecution := exec.Command("bash", "-c", tsharkCommand)
	fmt.Println(tsharkCommand)
	tsharkExecution.Run()
}

func makeCommandTshark(config *map[string]string) string {
	/*
		Input: The configuration map
		Output:	The string of the command for the execution of tshark with the args defined in config.yml
		Execution:	It checks the map and concats the  begin of command with the output_dirname/filename
	*/
	command := "sudo tshark -a duration:" + strconv.Itoa(getDuration(config)) + " -i " + (*config)["NETWORK_INTERFACE"] + " -T ek > " + others.ObtainDirectory((*config)["PACKETS_OUTPUT_DIRNAME"])
	if strings.Compare((*config)["PACKETS_OUTPUT_FILENAME"], "") == 0 {
		//This attribute is not defined, then we use the predefined name
		command += "packets.json"
	} else {
		//Add the filename
		command += (*config)["PACKETS_OUTPUT_FILENAME"] + ".json"
	}
	return command
}

func getDuration(config *map[string]string) int {
	/*
		Input: The configuration map
		Output:	The amount of seconds defined in config.yml, if it is not defined => we use the default duration(5 seconds)
		Execution: We parse the keys "TIME_SECONDS" and "TIME_MINUTES", if these are defined => then it stores the add in seconds of both in the variable duration
	*/
	duration := 10
	if seconds, error := strconv.ParseInt((*config)["TIME_SECONDS"], 10, 0); error == nil && seconds >= 0 {
		//The error is nil => There is a amount of seconds defined, we have to check if it's negative
		secs := int(seconds)
		duration = secs //Modify the duration
	} else {
		others.CheckError(errors.New("Error in the parsing of TIME_SECONDS"))
	}
	if minutes, error := strconv.ParseFloat((*config)["TIME_MINUTES"], 64); error == nil && minutes >= 0.0 {
		//The error is nil => There is a amount of minutes defined, we have to check if it's negative
		duration += int(minutes * 60) //Modify the duration
	} else {
		others.CheckError(errors.New("Error in the parsing of TIME_MINUTES"))
	}
	if duration <= 0 {
		duration = 10
	}
	duration = 10
	return duration
}
