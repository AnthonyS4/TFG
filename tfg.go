package main

import (
	"./elasticsearch"
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
	//filebeatExecution(configuration)
	elasticsearch.DeployElasticSearchNetwork(configuration) //Start the nodes execution

	//executeKibana(args)
}
