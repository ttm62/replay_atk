package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Yêu cầu được ghi lại từ một lần gửi hợp lệ trước đó
	data := `{"message": "Hello, secure world!"}`
	timestamp := "1700000000" // Timestamp cũ
	nonce := "1234567890"     // Nonce đã được sử dụng
	hmacSignature := "ABC123FakeSignature"

	req, _ := http.NewRequest("POST", "http://localhost:8080/secure", bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Timestamp", timestamp)
	req.Header.Set("Nonce", nonce)
	req.Header.Set("HMAC", hmacSignature)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Replay Attack Response:", string(body))
}
