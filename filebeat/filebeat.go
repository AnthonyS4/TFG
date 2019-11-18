package filebeat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"../others"
)

// FilebeatExecution :Function that executes filebeat using filebeat.yml
func Execution(config *map[string]string, tsharkCommand string) {
	packetsDest := strings.Split(tsharkCommand, " > ")[1]
	changeConfigPermissions(config)
	changeYML(config, packetsDest)
	//commandExecutionFilebeat := exec.Command("bash", "-c", makeCommandFilebeat(config))
	fmt.Println(makeCommandFilebeat(config))
}

func changeYML(config *map[string]string, packetsDest string) {
	pwdOutput, _ := exec.Command("bash", "-c", "pwd").Output()
	//Get our working directory of this project, add the filebeatModel.yml location
	source := strings.Replace(string(pwdOutput), "\n", "", 1) + "/filebeat/filebeatModel.yml"
	//Modify the model.yml with the tshark output directory.
	newData := modifySource(source, packetsDest)
	ioutil.WriteFile(source, newData, 0644)
	dest := others.ObtainDirectory((*config)["FILEBEAT_PATH"]) + "/filebeat.yml"
	exec.Command("bash", "-c", "sudo cp "+source+" "+dest).Run()
}

func modifySource(sourceDirectory string, packetsDest string) []byte {
	data, errorRead := ioutil.ReadFile(sourceDirectory)
	others.CheckError(errorRead)
	lines := strings.Split(string(data), "\n")
	//In the lines numerb 59 of the filebeatModel.yml is the packets directory, it modifies it
	lines[59] = "    - " + packetsDest
	//Return the byte array of the modified source.
	return newSource(lines, sourceDirectory)
}

func newSource(lines []string, sourceDirectory string) []byte {
	newSource := ""
	for i := 0; i < len(lines)-1; i++ {
		newSource += lines[i] + "\n"
	}
	newSource += lines[len(lines)-1]
	return []byte(newSource)
}

func changeConfigPermissions(config *map[string]string) {
	changeCommand := "sudo chown root" + others.ObtainDirectory((*config)["FILEBEAT_PATH"]) + "filebeat.yml"
	changePermissions := exec.Command("bash", "-c", changeCommand)
	changePermissions.Run()
}

func makeCommandFilebeat(config *map[string]string) string {
	/*
		Input:
		Output:
		Execution:
	*/
	if strings.Compare((*config)["FILEBEAT_PATH"], "") == 0 {
		others.CheckError(errors.New("The FILEBEAT_PATH is not defined"))
	}
	command := "sudo " + others.ObtainDirectory((*config)["FILEBEAT_PATH"]) + "filebeat -e"
	return command
}
