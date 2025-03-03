package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal"
	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils"
)

func main() {
	configFile := flag.String("config", "./config.json", "Path to the json configuration")
	host := flag.String("host", "", "Host to replace the one in the configuration file")
	showHelp := flag.Bool("help", false, "Show help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Use: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	configuration, err := internal.LoadConfigurationFromFile(*configFile)
	if err != nil {
		fmt.Println("There was an error when trying to load the configuration file:")
		fmt.Println(err)
		os.Exit(1)
	}

	if *host != "" {
		configuration.ReplaceHost(*host)
	}

	client := utils.NewClient(configuration.Host, "v1")
	configuration.Sync(client)
}
