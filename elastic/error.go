package elastic

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type ElasticError struct {
	Data   map[string]any
	Err    error
	Status int
	Reason string
}

func (j *ElasticError) Error() string {
	return fmt.Sprintf("[%d] %s: %v", j.Status, j.Reason, j.Err)
}

func (j *ElasticError) Unwrap() error {
	return j.Err
}

// Utility function to create a structured error from a http response
func newElasticError(res *esapi.Response) *ElasticError {
	var reason string
	data := map[string]any{}
	err := json.NewDecoder(res.Body).Decode(&data)
	if err == nil {
		if respErr, ok := data["error"].(map[string]any); ok {
			if rawReason, ok := respErr["reason"].(string); ok {
				reason = rawReason
			}
		}
	}

	return &ElasticError{
		Err:    err,
		Data:   data,
		Status: res.StatusCode,
		Reason: reason,
	}
}
