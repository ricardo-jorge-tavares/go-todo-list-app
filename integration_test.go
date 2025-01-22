package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"local.com/todo-list-app/test"
)

func TestIntegrationSuite(t *testing.T) {

	ts, teardown := test.RunTestServer()
	defer teardown(t)

	var userId = "f6645494-56f5-40e4-a3f9-b6ae8935006b"

	t.Run("should create a todo for an existing user", func(t *testing.T) {

		requestBody := bytes.NewBuffer([]byte(`
		{
			"description": "Taks to do something"
		}`))

		resp, _ := http.Post(fmt.Sprintf("%s/api/%s/todo/", ts.URL, userId), "application/json", requestBody)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status code 200, got %v", resp.StatusCode)
		}

		var response struct {
			Id string `json:"id"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		if response.Id == "" {
			t.Errorf("Expected non empty id, got %v", response.Id)
		}

	})

	t.Run("should fail with an invalid body request", func(t *testing.T) {

		requestBody := bytes.NewBuffer([]byte(``))

		resp, _ := http.Post(fmt.Sprintf("%s/api/%s/todo/", ts.URL, userId), "application/json", requestBody)

		if resp.StatusCode != 412 {
			t.Errorf("Expected status code 412, got %v", resp.StatusCode)
		}

	})

	t.Run("should fail if description is empty", func(t *testing.T) {

		requestBody := bytes.NewBuffer([]byte(`
		{
			"description": ""
		}`))

		resp, _ := http.Post(fmt.Sprintf("%s/api/%s/todo/", ts.URL, userId), "application/json", requestBody)

		if resp.StatusCode != 412 {
			t.Errorf("Expected status code 412, got %v", resp.StatusCode)
		}

	})

	t.Run("it should return ok when insert new post successfully", func(t *testing.T) {

		resp, _ := http.Get(fmt.Sprintf("%s/api/%s/", ts.URL, userId))

		type GetUserTodoListResponse struct {
			Id          string    `json:"id"`
			Description string    `json:"description"`
			IsCompleted bool      `json:"isCompleted"`
			Rank        int       `json:"rank"`
			CreatedAt   time.Time `json:"createdAt"`
		}

		var returnData []GetUserTodoListResponse
		err := json.NewDecoder(resp.Body).Decode(&returnData)
		if err != nil {
			t.Errorf("Error parsing response, got %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Expected status code 412, got %v", resp.StatusCode)
		}

		if len(returnData) != 1 {
			t.Errorf("Expected to have 1 user, got %v", len(returnData))
		}

	})

}
