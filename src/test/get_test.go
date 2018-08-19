package test

import (
	"net/http"
	"testing"
)

var token string = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJoYXNoIjoiUmlkIiwiZXhwIjoxNTM0NzY5OTk3LCJqdGkiOiJtYWluX3VzZXJfaWQifQ.VCRKMNbUnW8VA2A-5Xv8jX5-fm8GDXL9Xwg1Ffy7Vj7wMkz8Ov97Lst2tliCdkstrMGBJd7XnweWF7nsON38sQ0"

// func BenchmarkFixedWindowCounter(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		url := "http://localhost:8000/v1.0/auth/user/login"

// 		var jsonStr = []byte(`{
// 			"email": "rid@shade.com",
// 			"password":"airin123"
// 		}`)
// 		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
// 		req.Header.Set("X-Custom-Header", "myvalue")
// 		req.Header.Set("Content-Type", "application/json")

// 		client := &http.Client{}
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer resp.Body.Close()

// 	}
// }

func BenchmarkFixedWindowCounter(b *testing.B) {
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

func BenchmarkSlidingWindowCounter(b *testing.B) {
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

func BenchmarkSlidingWindowLog(b *testing.B) {
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
