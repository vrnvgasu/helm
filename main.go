package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

type EchoHandler struct {
	Hostname        string
	EchoResponse    string
	ResponseSignKey string
}

func (h EchoHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Первая строка — метод, путь и адрес пода
	fmt.Fprintf(w, "[%s] %s %s (from pod: %s)\n", r.Proto, r.Method, r.URL.Path, h.Hostname)

	var responseBody string

	// Если задана переменная EchoResponse, включаем её в ответ
	if h.EchoResponse != "" {
		responseBody = h.EchoResponse
	} else {
		// Копируем запрос клиента в ответ
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if len(body) > 0 {
			responseBody = string(body)
		} else {
			responseBody = "(empty body)"
		}
	}
	fmt.Fprintf(w, "%s", responseBody)

	if signature := h.sign(responseBody); signature != "" {
		fmt.Fprintf(w, "\n%s", signature)
	}
}

func (h EchoHandler) sign(content string) string {
	if h.ResponseSignKey == "" {
		return ""
	}

	hasher := hmac.New(sha256.New, []byte(h.ResponseSignKey))
	hasher.Write([]byte(content))
	return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
	hostname, _ := os.Hostname()
	echoResponse := os.Getenv("ECHO_SERVER_RESPONSE")
	echoResponseSignKey := os.Getenv("ECHO_SERVER_RESPONSE_SIGN_KEY")

	echoHandler := EchoHandler{
		Hostname:        hostname,
		EchoResponse:    echoResponse,
		ResponseSignKey: echoResponseSignKey,
	}

	http.HandleFunc("/", echoHandler.Handle)

	addr := net.JoinHostPort("", "8080")
	fmt.Println("Starting echo server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
