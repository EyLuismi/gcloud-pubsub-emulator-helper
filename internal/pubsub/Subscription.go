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

// String returns a JSON string representation of the Subscription.
func (t *Subscription) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("error marshaling Subscription: %v", err)
	}
	return string(b)
}

// GetSubscription retrieves a subscription by its resource name.
func GetSubscription(
	client utils.ClientInterface,
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
	client utils.ClientInterface,
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
	client utils.ClientInterface,
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
	client utils.ClientInterface,
	project, subscriptionResourceName, topicResourceName string,
	labels *Labels,
) error {
	exists, err := IsSubscriptionPresent(client, project, subscriptionResourceName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// Prepare the request body.
	type CreateSubscriptionBody struct {
		Topic  string `json:"topic"`
		Labels Labels `json:"labels"`
	}

	createSubscriptionBody := CreateSubscriptionBody{Topic: topicResourceName}

	if labels != nil {
		createSubscriptionBody.Labels = *labels
	}

	rawBody, err := json.Marshal(createSubscriptionBody)
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
	client utils.ClientInterface,
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
