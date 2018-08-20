package test

import (
	"bytes"
	"net/http"
	"testing"
)

func BenchmarkFixedWindowCounterPost(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := "http://localhost:8000/v1.0/fwc/publish"

		var jsonStr = []byte(`{
			"AuthorID": 1,
			"Body":     "Go golang rocks! ",
			"Tags":     ["intro", "golang"]
		}`)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

	}
}

func BenchmarkSlidingCounterPost(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := "http://localhost:8000/v1.0/swc/publish"

		var jsonStr = []byte(`{
			"AuthorID": 1,
			"Body":     "Go golang rocks! ",
			"Tags":     ["intro", "golang"]
		}`)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

	}
}

func BenchmarkSlidingWindowLogPost(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := "http://localhost:8000/v1.0/swl/publish"

		var jsonStr = []byte(`{
			"AuthorID": 1,
			"Body":     "Go golang rocks! ",
			"Tags":     ["intro", "golang"]
		}`)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

	}
}

func BenchmarkTokenBucketPost(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := "http://localhost:8000/v1.0/tb/publish"

		var jsonStr = []byte(`{
			"AuthorID": 1,
			"Body":     "Go golang rocks! ",
			"Tags":     ["intro", "golang"]
		}`)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

	}
}
