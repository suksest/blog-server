package test

import (
	"net/http"
	"testing"
)

var token string = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJoYXNoIjoiUmlkIiwiZXhwIjoxNTM0NzY5OTk3LCJqdGkiOiJtYWluX3VzZXJfaWQifQ.VCRKMNbUnW8VA2A-5Xv8jX5-fm8GDXL9Xwg1Ffy7Vj7wMkz8Ov97Lst2tliCdkstrMGBJd7XnweWF7nsON38sQ0"

func BenchmarkFixedWindowCounterGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := "http://localhost:8000/v1.0/fwc/posts"

		req, err := http.NewRequest("GET", url, nil)
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

func BenchmarkSlidingWindowCounterGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := "http://localhost:8000/v1.0/swc/posts"

		req, err := http.NewRequest("GET", url, nil)
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

func BenchmarkSlidingWindowLogGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := "http://localhost:8000/v1.0/swl/posts"

		req, err := http.NewRequest("GET", url, nil)
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

func BenchmarkTokenBucketGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := "http://localhost:8000/v1.0/tb/posts"

		req, err := http.NewRequest("GET", url, nil)
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

}
