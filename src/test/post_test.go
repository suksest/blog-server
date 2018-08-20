package test

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
