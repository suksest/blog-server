package test

import (
	"bytes"
	"net/http"
	"testing"
)

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
		url := "http://localhost:8000/v1.0/auth/user/login"

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJoYXNoIjoiUmlkIiwiZXhwIjoxNTM0NDk5MzAwLCJqdGkiOiJtYWluX3VzZXJfaWQifQ.xGUbh_w-g0u9GwM_BdxpdkG12IHTaC3RcMY87cR-UGHayhOR42-t_rgR9AHySZ7kOLvAFr7ETwYpvheBUe_F0Q")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

	}
}

func BenchmarkFixedWindowCounter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := "http://localhost:8000/v1.0/auth/user/login"

		var jsonStr = []byte(`{
			"email": "rid@shade.com",
			"password":"airin123"
		}`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

	}
}
