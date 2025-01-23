package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"local.com/todo-list-app/test"
)

var executionCounter = 0
var itemsCreated = 0

func BenchmarkSuite(b *testing.B) {

	executionCounter++

	ts, teardown := test.RunTestServer()
	defer teardown(b)

	userId := "e5a91497-ae4c-48cd-a73d-1d43b34a7b9b"

	b.Run("should create a todo for an existing user", func(b *testing.B) {

		// b.ReportAllocs()
		// b.ResetTimer()

		for i := 0; i < b.N; i++ {

			requestBody := bytes.NewBuffer([]byte(`{"description": "Taks to do something"}`))

			_, err := http.Post(fmt.Sprintf("%s/api/%s/todo/", ts.URL, userId), "application/json", requestBody)
			if err != nil {
				b.Errorf("Failed request: %v", err)
			}

			itemsCreated++

		}

	})

	fmt.Println("Items created: ", itemsCreated)
}
