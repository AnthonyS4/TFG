package main

import "os/exec"

func main() {
	Directorio := "/home/anthony/elasticsearch-7.4.2/bin"
	exec.Command("bash", "-c", "cd "+Directorio)
	Commando := "./elasticsearch"
	exec.Command("bash", "-c", Commando)

}
