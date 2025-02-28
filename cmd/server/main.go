package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var secretKey = []byte("supersecretkey")
var requestCache = make(map[string]bool) // Lưu các nonce để ngăn chặn replay

// Tạo HMAC cho dữ liệu
func generateHMAC(data string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Middleware kiểm tra HMAC
func validateHMAC(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(strings.NewReader(string(body))) // Reset body để đọc lại

		timestamp := r.Header.Get("Timestamp")
		nonce := r.Header.Get("Nonce")
		hmacHeader := r.Header.Get("HMAC")

		// Ngăn chặn replay attack bằng cách kiểm tra nonce
		if _, exists := requestCache[nonce]; exists {
			http.Error(w, "Replay attack detected", http.StatusUnauthorized)
			return
		}

		// Xác thực chữ ký HMAC
		expectedHMAC := generateHMAC(string(body) + timestamp + nonce)
		if hmacHeader != expectedHMAC {
			http.Error(w, "Invalid HMAC", http.StatusUnauthorized)
			return
		}

		// Lưu lại nonce để chống replay
		requestCache[nonce] = true
		next.ServeHTTP(w, r)
	})
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Request accepted: Secure data")
}

func main() {
	r := mux.NewRouter()
	r.Handle("/secure", validateHMAC(http.HandlerFunc(protectedHandler))).Methods("POST")

	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
