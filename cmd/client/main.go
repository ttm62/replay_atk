package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"
)

var secretKey = []byte("supersecretkey")

func generateHMAC(data string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func main() {
	data := `{"message": "Hello, secure world!"}`
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := fmt.Sprintf("%d", time.Now().UnixNano()) // Nonce duy nhất cho mỗi request
	hmacSignature := generateHMAC(data + timestamp + nonce)

	fmt.Println(data, timestamp, nonce)

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
	fmt.Println(string(body))
}
