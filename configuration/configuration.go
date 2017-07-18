package configuration


import (
	"io/ioutil"
	"os"
	"fmt"
	"bufio"

	"gopkg.in/yaml.v2"
	"github.com/op/go-logging"
)


// Global vars
var config ConfigType
var log *logging.Logger

// Instance configuration
type ConfigType struct {

	Server_port string
	Log_file string
	Log_format string
	Mail_to string
	Slack_api_token string
	Slack_channel_name string
}



/**
 * Load configuration yaml file
 */
func LoadConfiguration(filename string) ConfigType {

	// Set config_pre
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	printBootLogo()
	fmt.Printf("--> Configuration loaded values: %#v\n", config)


	// Set logger
	format := logging.MustStringFormatter( config.Log_format )
	logbackend1 := logging.NewLogBackend(os.Stdout, "", 0)
	logbackend1Formatted := logging.NewBackendFormatter(logbackend1, format)
	f, err := os.OpenFile(config.Log_file, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		defer f.Close()
	}
	logbackend2 := logging.NewLogBackend(f, "", 0)
	logbackend2Formatted := logging.NewBackendFormatter(logbackend2, format)
	logging.SetBackend(logbackend1Formatted, logbackend2Formatted)

	log = logging.MustGetLogger("formhandler")
	log.Info("Application started successfuly.")

	return config
}


// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}


func printBootLogo() {

	lines, _ := readLines("boot_logo.txt")
	for _, line := range lines {
		fmt.Println(line)
	}

}


func GetLog() *logging.Logger {
	return log
}