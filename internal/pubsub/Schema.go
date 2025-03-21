package pubsub

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils"
)

type Schema struct {
	Id                 string `json:"id,omitempty"`
	Name               string `json:"name"`
	Type               string `json:"type"`
	Definition         string `json:"definition"`
	RevisionId         string `json:"revisionId,omitempty"`
	RevisionCreateTime string `json:"revisionCreateTime,omitempty"`
}

// String returns a JSON string representation of the Schema.
func (s *Schema) String() string {
	jsonBytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func CreateSchema(client utils.ClientInterface, project, schemaId, name, schemaType, definition string) error {
	type CreateSchemaRequest struct {
		Name       string `json:"name"`
		Type       string `json:"type"`
		Definition string `json:"definition"`
	}

	body, err := json.Marshal(
		CreateSchemaRequest{
			Name:       GetResourceNameForSchema(project, name),
			Type:       schemaType,
			Definition: definition,
		},
	)
	if err != nil {
		return err
	}

	prettyBody, err := json.MarshalIndent(
		CreateSchemaRequest{
			Name:       GetResourceNameForSchema(project, schemaId),
			Type:       schemaType,
			Definition: definition,
		},
		"", "  ",
	)
	if err != nil {
		return err
	}
	fmt.Printf("Creating schema with body:\n%s\n", string(prettyBody))

	response, err := client.Post(
		fmt.Sprintf("projects/%s/schemas?schemaId=%s", project, schemaId),
		body,
	)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return fmt.Errorf("unexpected status code %d in CreateSchema", response.StatusCode)
	}
}

func IsSchemaPresent(client utils.ClientInterface, schemaResourceName string) (bool, error) {
	response, err := client.Get(schemaResourceName)
	if err != nil {
		return false, err
	}

	switch response.StatusCode {
	case http.StatusNotFound:
		return false, nil
	case http.StatusOK:
		return true, nil
	default:
		return false, fmt.Errorf("unexpected status code %d in IsSchemaPresent", response.StatusCode)
	}
}

func ListSchemas(client utils.ClientInterface, project string) ([]Schema, error) {
	url := fmt.Sprintf("projects/%s/schemas?view=FULL", project)
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	type listSchemasResponse struct {
		Schemas       []Schema `json:"schemas"`
		NextPageToken string   `json:"nextPageToken"`
	}

	switch response.StatusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("project not found")
	case http.StatusOK:
		var res listSchemasResponse
		if err := json.Unmarshal(response.Body, &res); err != nil {
			return nil, err
		}
		return res.Schemas, nil
	default:
		return nil, fmt.Errorf("unexpected status code %d in ListSchemas", response.StatusCode)
	}
}

func GetSchemaRevisionIdBySchemaId(client utils.ClientInterface, project, schemaId string) (string, error) {
	response, err := client.Get(fmt.Sprintf("projects/%s/schemas/%s", project, schemaId))
	if err != nil {
		return "", err
	}

	switch response.StatusCode {
	case http.StatusOK:
		var schema Schema
		if err := json.Unmarshal(response.Body, &schema); err != nil {
			return "", err
		}
		return schema.RevisionId, nil
	default:
		return "", fmt.Errorf("unexpected status code %d in GetSchemaRevisionIdBySchemaId", response.StatusCode)
	}
}

func GetResourceNameForSchema(project, schema string) string {
	return fmt.Sprintf("projects/%s/schemas/%s", project, schema)
}
