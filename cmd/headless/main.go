package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal"
	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/pubsub"
	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils"
	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils/Llog"
)

func main() {
	Llog.Init()

	configFile := flag.String("config", "./config.json", "Path to the json configuration")
	host := flag.String("host", "", "Host to replace the one in the configuration file")
	showHelp := flag.Bool("help", false, "Show help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Use: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	Llog.Debug(fmt.Sprintf("Using as 'config' flag value '%s'", *configFile))
	Llog.Debug(fmt.Sprintf("Using as 'host' flag value '%v'", *host))
	Llog.Debug(fmt.Sprintf("Using as 'showHelp' flag value '%v'", *showHelp))

	if *showHelp {
		Llog.Debug("Showing help menu and exiting")
		flag.Usage()
		os.Exit(0)
	}

	configuration, err := internal.LoadConfigurationFromFile(&utils.FileReader{}, *configFile)
	if err != nil {
		fmt.Println("There was an error when trying to load the configuration file:")
		fmt.Println(err)
		os.Exit(1)
	}

	if *host != "" {
		Llog.Debug(
			fmt.Sprintf(
				"Host given, trying to replace actual value from '%s' to '%s'",
				configuration.Host,
				*host,
			),
		)
		configuration.ReplaceHost(*host)
		Llog.Debug(fmt.Sprintf("Using host '%s'", *host))
	}

	client := utils.NewClient(configuration.Host, "v1")
	configuration.Sync(client)

	topicsList, err := pubsub.ListTopics(client, configuration.Projects[0].Name)
	if err != nil {
		fmt.Println("There was some error while trying to list the topics")
		os.Exit(1)
	}

	for _, topic := range topicsList {
		Llog.Debug(topic.String())
	}
}
