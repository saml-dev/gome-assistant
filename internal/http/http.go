// http is used to interact with the home assistant
// REST API. Currently only used to retrieve state for
// a single entity_id
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HttpClient struct {
	url   string
	token string
}

func NewHttpClient(url *url.URL, token string) *HttpClient {
	// Shallow copy the URL to avoid modifying the original
	u := *url
	u.Path = "/api"
	if u.Scheme == "ws" {
		u.Scheme = "http"
	}
	if u.Scheme == "wss" {
		u.Scheme = "https"
	}

	return &HttpClient{
		url:   u.String(),
		token: token,
	}
}

func (c *HttpClient) GetState(entityID string) ([]byte, error) {
	return c.GetStateWithContext(context.Background(), entityID)
}

// GetStateWithContext retrieves an entity state and stops the request when
// ctx is canceled.
func (c *HttpClient) GetStateWithContext(ctx context.Context, entityID string) ([]byte, error) {
	resp, err := get(ctx, c.url+"/states/"+entityID, c.token)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *HttpClient) GetStates(entityIDs []string) ([]byte, error) {
	template, err := statesTemplate(entityIDs)
	if err != nil {
		return nil, err
	}

	resp, err := post(c.url+"/template", c.token, map[string]string{
		"template": template,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func statesTemplate(entityIDs []string) (string, error) {
	stateCalls := make([]string, 0, len(entityIDs))
	for _, entityID := range entityIDs {
		quotedEntityID, err := json.Marshal(entityID)
		if err != nil {
			return "", err
		}
		stateCalls = append(stateCalls, "states("+string(quotedEntityID)+")")
	}

	return "{{ [" + strings.Join(stateCalls, ", ") + "] | to_json }}", nil
}

func (c *HttpClient) States() ([]byte, error) {
	resp, err := get(context.Background(), c.url+"/states", c.token)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func get(ctx context.Context, url, token string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.New("Error creating HTTP request: " + err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Error on response.\n[ERROR] -" + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Error while reading the response bytes:" + err.Error())
	}

	return body, nil
}

func post(url, token string, data any) ([]byte, error) {
	postBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, errors.New("Error creating HTTP request: " + err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Error on response.\n[ERROR] -" + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Error while reading the response bytes:" + err.Error())
	}

	return body, nil
}

// func post(url string, token string, data any) ([]byte, error) {
// 	postBody, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req, err := http.NewRequest("GET", url, bytes.NewBuffer(postBody))
// 	if err != nil {
// 		return nil, errors.New("Error building post request: " + err.Error())
// 	}

// 	req.Header.Add("Authorization", "Bearer "+token)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, errors.New("Error in post response: " + err.Error())
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == 401 {
// 		panic("ERROR: Auth token is invalid. Please double check it or create a new token in your Home Assistant profile")
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return body, nil
// }
