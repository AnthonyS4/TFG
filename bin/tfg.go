package main

import (
	"os"

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
	flags.Read(os.Args[1:]) //Read flags, need to implement the execution of it
	//configuration := others.ObtainConfiguration()    //Obtains a reference of a map with the configurations
	//wireshark.Execute(configuration, &tsharkCommand) //Executes tshark
	//tsharkCommand = "blabal/ablaa7asdb/ > /home/anthony/TFG/archivos/paquetes.json"
	//filebeat.Execution(configuration, tsharkCommand)
	//elasticsearch.DeployElasticSearchNetwork(configuration) //Start the nodes execution

	//executeKibana(args)
}
