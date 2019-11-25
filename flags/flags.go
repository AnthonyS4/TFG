package flags

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

//Execute :This function read the flags and values inserted for the program execution
func Execute() *map[string]string {
	configuration := make(map[string]string)
	command := getCobraCommand(&configuration)
	if command.Execute() != nil {
		fmt.Println("Vaya bug")
		os.Exit(1)
	}
	return &configuration
}

func getCobraCommand(config *map[string]string) *cobra.Command {
	var rootCmd = &cobra.Command{}
	var auxCmd = &cobra.Command{
		Use:   "Use",
		Short: "Analyse network data",
		Long:  "Analyse is for capture",
		PreRun: func(cmd *cobra.Command, args []string) {
			exec.Command("bash", "-c", "sudo pwd").Run()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if timelimit, errorTimeL := rootCmd.Flags().GetString("timelimit"); errorTimeL != nil {
				(*config)["TIMELIMIT"] = ""
			} else {
				(*config)["TIMELIMIT"] = timelimit
			}
			if begin, errorBegin := rootCmd.Flags().GetString("begin"); errorBegin != nil {
				(*config)["BEGIN"] = ""
			} else {
				(*config)["NOW"] = ""
				(*config)["BEGIN"] = begin
			}
			if end, errorEnd := rootCmd.Flags().GetString("end"); errorEnd != nil {
				(*config)["END"] = ""
			} else {
				(*config)["END"] = end
			}
		},
	}
	rootCmd = auxCmd
	setTemplates(rootCmd)
	putFlags(config, rootCmd)
	return rootCmd
}

func putFlags(config *map[string]string, rootCmd *cobra.Command) {
	var timeLimit, begin, end string
	rootCmd.Flags().StringVarP(&timeLimit, "timelimit", "t", "", "Limite de tiempo de captura en segundos")
	rootCmd.Flags().StringVarP(&begin, "begin", "b", "", "Fijar la hora de comienzo de la captura")
	rootCmd.Flags().StringVarP(&end, "end", "e", "", "Fijar la hora de comienzo de la captura")
}

func setTemplates(rootCmd *cobra.Command) {
	rootCmd.SetHelpTemplate("OUTPUT FOR HELP FLAG")
	rootCmd.SetUsageTemplate("./rc analyse [-n now] ")
	rootCmd.SetVersionTemplate("Version 1.0.0")
}
