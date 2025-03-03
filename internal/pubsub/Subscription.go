package pubsub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils"
)

// Subscription represents a Pub/Sub subscription.
type Subscription struct {
	Name   string `json:"name"`
	Labels Labels `json:"labels"`
}

// GetSubscription retrieves a subscription by its resource name.
func GetSubscription(
	client utils.Client,
	project, subscriptionResourceName string,
) (*Subscription, error) {
	response, err := client.Get(subscriptionResourceName)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusNotFound:
		return nil, errors.New("subscription not found")
	case http.StatusOK:
		var sub Subscription
		if err := json.Unmarshal(response.Body, &sub); err != nil {
			return nil, err
		}
		return &sub, nil
	default:
		return nil, fmt.Errorf("unexpected status code %d in GetSubscription", response.StatusCode)
	}
}

// listSubscriptionsResponse is the internal structure for unmarshalling the ListSubscriptions response.
type listSubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

// ListSubscriptions retrieves all subscriptions for a given project.
func ListSubscriptions(
	client utils.Client,
	project string,
) ([]Subscription, error) {
	// Build the URL for listing subscriptions.
	url := fmt.Sprintf("projects/%s/subscriptions", project)
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		var res listSubscriptionsResponse
		if err := json.Unmarshal(response.Body, &res); err != nil {
			return nil, err
		}
		return res.Subscriptions, nil
	default:
		return nil, fmt.Errorf("unexpected status code %d in ListSubscriptions", response.StatusCode)
	}
}

// IsSubscriptionPresent checks if a subscription exists.
func IsSubscriptionPresent(
	client utils.Client,
	project, subscriptionResourceName string,
) (bool, error) {

	response, err := client.Get(subscriptionResourceName)
	if err != nil {
		return false, err
	}

	switch response.StatusCode {
	case http.StatusNotFound:
		return false, nil
	case http.StatusOK:
		return true, nil
	default:
		return false, fmt.Errorf("unexpected status code %d in IsSubscriptionPresent", response.StatusCode)
	}
}

// CreateSubscription creates a subscription for a topic if it does not exist.
func CreateSubscription(
	client utils.Client,
	project, subscriptionResourceName, topicResourceName string,
) error {
	exists, err := IsSubscriptionPresent(client, project, subscriptionResourceName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// Prepare the request body.
	type createSubscriptionBody struct {
		Topic string `json:"topic"`
	}
	rawBody, err := json.Marshal(createSubscriptionBody{Topic: topicResourceName})
	if err != nil {
		return err
	}
	response, err := client.Put(subscriptionResourceName, rawBody)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error creating subscription: status code %d", response.StatusCode)
	}

	return nil
}

// DeleteSubscription deletes a subscription.
func DeleteSubscription(
	client utils.Client,
	project, subscriptionResourceName string,
) error {
	// Build the full resource name.
	response, err := client.Delete(subscriptionResourceName)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting subscription: status code %d", response.StatusCode)
	}

	return nil
}

// GetResourceNameForSubscription generates the full resource name for a subscription.
func GetResourceNameForSubscription(project, subscription string) string {
	return fmt.Sprintf("projects/%s/subscriptions/%s", project, subscription)
}
