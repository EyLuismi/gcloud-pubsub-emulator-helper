package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/pubsub"
	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils"
)

type Configuration struct {
	Host                       string           `json:"host"`
	StartTimeoutMs             int              `json:"startTimeoutMs"`
	AvoidStartupCheck          bool             `json:"avoidStartupCheck"`
	Projects                   []pubsub.Project `json:"projects"`
	TimeBetweenStartupChecksMs int              `json:"timeBetweenStartupChecksMs"`
	DelayBeforeStartupCheckMs  int              `json:"delayBeforeStartupCheckMs"`
}

func (c Configuration) String() string {
	raw, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic("Can't marshal configuration!")
	}

	return string(raw)
}

func LoadConfigurationFromFile(fileReader utils.FileReaderInterface, filepath string) (Configuration, error) {
	// TODO: Create an intermediate configuration schema to decouple Configuration struct <=> file format
	configurationFileBytes, err := fileReader.Read(filepath)
	if err != nil {
		return Configuration{}, err
	}

	var configuration Configuration
	err = json.Unmarshal(configurationFileBytes, &configuration)
	if err != nil {
		return Configuration{}, err
	}

	configuration.Host = strings.Trim(configuration.Host, " ")

	if configuration.Host == "" {
		configuration.Host = "localhost:8085"
	}

	if strings.HasPrefix(configuration.Host, ":") {
		configuration.Host = "localhost" + configuration.Host
	}

	if configuration.DelayBeforeStartupCheckMs < 0 {
		configuration.DelayBeforeStartupCheckMs = 0
	}

	if configuration.StartTimeoutMs <= 0 {
		configuration.StartTimeoutMs = 30_000
	}

	if configuration.TimeBetweenStartupChecksMs <= 0 {
		configuration.TimeBetweenStartupChecksMs = 200
	}

	if !utils.IsValidHost(configuration.Host) {
		fmt.Println("The given host is invalid")
		os.Exit(1)
	}

	return configuration, nil
}

func (c Configuration) ReplaceHost(host string) Configuration {
	if !utils.IsValidHost(host) {
		fmt.Println("The given host is invalid")
		os.Exit(1)
	}

	c.Host = host
	return c
}

/**
*	Sync will remove everything in the emulator and then apply the configuration
* TODO: In the future, it should have an option to just update what is required
* 	to preserve data in those topics/subscriptions
 */
func (c *Configuration) Sync(client utils.ClientInterface) {
	// Wait until the emulator is running
	if !c.AvoidStartupCheck {
		startTime := time.Now()
		for {
			_, err := client.Get("")
			if err != nil {
				if time.Since(startTime).Milliseconds() > int64(c.StartTimeoutMs) {
					fmt.Println("Time to start the emulator has been exceeded. Exiting...")
					os.Exit(1)
				}
				time.Sleep(200 * time.Millisecond)
				continue
			}
			break
		}
	}

	// Cleaning first everything in the emulator
	for _, project := range c.Projects {
		topics, err := pubsub.ListTopics(client, project.Name)
		if err != nil {
			topics = []pubsub.Topic{}
		}

		for _, topic := range topics {
			pubsub.DeleteTopic(client, project.Name, topic.Name)
		}

		subscriptions, err := pubsub.ListSubscriptions(client, project.Name)
		if err != nil {
			subscriptions = []pubsub.Subscription{}
		}

		for _, subscription := range subscriptions {
			pubsub.DeleteSubscription(client, project.Name, subscription.Name)
		}
	}

	// Applying information in the configuration
	for _, project := range c.Projects {
		for _, topic := range project.Topics {
			pubsub.CreateTopic(
				client,
				project.Name,
				pubsub.GetResourceNameForTopic(
					project.Name,
					topic.Name,
				),
				&topic.Labels,
				&topic.MessageStoragePolicy,
			)
			for _, subscription := range topic.Subscriptions {
				pubsub.CreateSubscription(
					client,
					project.Name,
					pubsub.GetResourceNameForSubscription(
						project.Name,
						subscription.Name,
					),
					pubsub.GetResourceNameForTopic(
						project.Name,
						topic.Name,
					),
					&subscription.Labels,
				)
			}
		}
	}
}
