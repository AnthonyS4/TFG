package others

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// ObtainConfiguration :This function reads the file config.yml and make a map with the attibutes defined in the file
func ObtainConfiguration() *map[string]string {
	/*
		Input: ~
		Output: The reference of the map that will contains the configurations of the program
	*/
	configuration := make(map[string]string)
	data, error := ioutil.ReadFile("config.yml")
	CheckError(error) //Checking if the file exists
	lines := checkAttributes(data)
	for i := 0; i < len(lines); i++ { //This loop will read and store the keys and values defined in the config.yml
		keyAndValue := strings.Split(lines[i], ":")
		configuration[keyAndValue[0]] = keyAndValue[1]
	}
	return &configuration
}

func checkAttributes(data []byte) []string {
	/*
		Input: The array of bytes read in the config.yml file
		Output: An array of the string that will have the keys and values of the map
		Execution: It checks the existence of the substring ":" in the lines of the file, and erases the strings that don't have one
	*/
	configurationStrings := strings.Split(string(data), "\n")
	for i := 0; i < len(configurationStrings); i++ {
		if strings.Count(configurationStrings[i], ":") != 1 {
			configurationStrings = removeElement(configurationStrings, i)
			i--
		}
	}
	return configurationStrings
	//It returns the parametters except the last one because it is a line that contains the end of file
}

// ObtainDirectory :It checks the lexically if there is any "${" => there is a env. variable then it gets the absolute path of it, else it returns the absolute path given or ""
func ObtainDirectory(path string) string {
	/*
		Input: The path given in the key(PACKETS_OUTPUT_DIRNAME) of the config map
		Output:	The absolute path that packets will stored.
	*/
	variable := ""
	if strings.Compare(path, "") != 0 {
		if strings.Contains(path, "${") && strings.Contains(path, "}") { //Lexical check
			variable += lookupEnvVariable(path) + strings.SplitAfter(path, "}")[1] + "/"
		} else { //There is no env. variables, then we return the defined directory
			variable = path + "/"
		}
	}
	return variable
}

func lookupEnvVariable(path string) string {
	/*
		Input: The path with the env. variable
		Output: The absolute path defined in the env. variable
		Execution: It separates the env. variable taking out the "{" and "}", then it does an request using Getenv and gets it.
	*/
	variable := strings.Split(path, "}")[0]
	variable = strings.Split(variable, "{")[1]
	variableDirectory := os.Getenv(variable)
	if strings.Compare(variableDirectory, "") == 0 {
		//The env. variable doesn't exists, then stop the execution
		CheckError(errors.New("The env. variable " + variable + " is not defined"))
	} //Else then we return the new path
	return variableDirectory
}

// CheckError :It checks if the input error is nil, if it's not then it calls panic(input)
func CheckError(e error) {
	/*
		Input: The error obtained after the read of the file
		Output: ~
	*/
	if e != nil {
		panic(e)
	}
}

func removeElement(vector []string, index int) []string {
	/*
				Input: An array of strings and the index of the element to be erased
				Output: An array of the string that will have the keys and values of the map
		// RemoveElement :Uses append between the elements before index and the following elements.
	*/
	return append(vector[:index], vector[index+1:]...)
}
