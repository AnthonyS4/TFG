package main

import (
	"fmt"

	"../flags"
)

/*
	Input:
	Output:
	Execution:
*/

//
func main() {
	//var tsharkCommand string
	configuration := flags.Execute() //Read flags, need to implement the execution of it
	//configuration := others.ObtainConfiguration() //Obtains a reference of a map with the configurations
	fmt.Println("Key: ", "END", "Value: ", (*configuration)["END"])
	fmt.Println("Key: ", "BEGIN", "Value: ", (*configuration)["BEGIN"])
	fmt.Println("Key: ", "TIMELIMIT", "Value: ", (*configuration)["TIMELIMIT"])
	fmt.Println("Key: ", "NOW", "Value: ", (*configuration)["NOW"])
	//wireshark.Execute(configuration, &tsharkCommand) //Executes tshark
	//tsharkCommand = "blabal/ablaa7asdb/ > /home/anthony/TFG/archivos/paquetes.json"
	//filebeat.Execution(configuration, tsharkCommand)
	//elasticsearch.DeployElasticSearchNetwork(configuration) //Start the nodes execution

	//executeKibana(args)
}
