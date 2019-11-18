package main

import (
	"./filebeat"
	"./others"
)

/*
	Input:
	Output:
	Execution:
*/

//
func main() {
	configuration := others.ObtainConfiguration() //Obtains a reference of a map with the configurations
	//wireshark.ExecuteTshark(configuration)        //Executes tshark
	tsharkCommand := "blabal/ablaa7asdb/ > /home/anthony/TFG/archivos/paquetes.json"
	filebeat.FilebeatExecution(configuration, tsharkCommand)
	//elasticsearch.DeployElasticSearchNetwork(configuration) //Start the nodes execution

	//executeKibana(args)
}
