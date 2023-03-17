package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CriticalContexts []string `yaml:"criticalContexts"`
	Verbs            []string `yaml:"verbs"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: kubectl-confirm <context> <verb>")
		os.Exit(1)
	}

	context := os.Args[1]
	verb := os.Args[2]

	config := readConfig()
	criticalContexts := config.CriticalContexts
	verbs := config.Verbs

	verbMatch := false
	for _, v := range verbs {
		if v == verb {
			verbMatch = true
			break
		}
	}

	if !verbMatch {
		os.Exit(0)
	}

	for _, critical := range criticalContexts {
		g, err := glob.Compile(critical)
		if err != nil {
			fmt.Printf("Error in pattern compilation: %v\n", err)
			os.Exit(1)
		}
		if g.Match(context) {
			fmt.Printf("The context '%s' is critical and the verb '%s' is dangerous. Are you sure you want to proceed? (y/N): ", context, verb)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))

			if input != "y" && input != "yes" {
				fmt.Println("Operation canceled.")
				os.Exit(1)
			}
			break
		}
	}
}

func readConfig() Config {
	configPath := os.Getenv("KUBECTL_CONFIRM_CONFIG")
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting user home directory: %v\n", err)
			os.Exit(1)
		}
		configPath = filepath.Join(homeDir, "kube-confirm-config.yaml")
	}

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	config := Config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		os.Exit(1)
	}

	return config
}
