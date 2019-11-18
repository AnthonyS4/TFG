package flags

import (
	"errors"
	"strings"

	"../others"
)

//Read :This function read the flags and values inserted for the program execution
func Read(args []string) *map[string]string {
	flags := make(map[string]string)
	checkArgs(args)
	return &flags
}

func checkArgs(args []string) {
	allowed1, allowed2 := getAllowedFlags()
	for i := 0; i < len(args); i++ {
		if strings.Contains(args[i], "--") {
			check2(allowed2[:], args[i])
		} else {
			if strings.Contains(args[i], "-") {
				check1(allowed1[:], args[i])
			} else {
				others.CheckError(errors.New("The " + args[i] + " is not allowed as input"))
			}
		}
	}
}

func getAllowedFlags() ([]string, []string) {
	first := [...]string{"-n", "-t", "-h"}
	second := [...]string{"--help", "--now", "--timelimit"}
	return first[:], second[:]
}

func check2(allowed []string, arg string) {
	if !contains(allowed[:], arg) {
		others.CheckError(errors.New("The flag" + arg + " is not allowed"))
	}
}

func check1(allowed []string, arg string) {
	if !contains(allowed[:], arg) {
		others.CheckError(errors.New("The flag" + arg + " is not allowed"))
	}
}

func contains(vector []string, element string) bool {
	for i := 0; i < len(vector); i++ {
		if strings.Compare(vector[i], element) == 0 {
			return true
		}
	}
	return false
}
