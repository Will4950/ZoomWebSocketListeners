package main

import (
	"encoding/json"
	"fmt"
	"go/v2/websocketservice"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func newUserCreatedHandler() {
	log.Printf("\n\n\n\nA new user was created")
	log.Printf("Do some processing\n\n\n\n")
}

func main() {
	if err := godotenv.Load("../.env.local"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	accountID := os.Getenv("accountId")
	clientID := os.Getenv("clientId")
	clientSecret := os.Getenv("clientSecret")
	wsURL := os.Getenv("url")

	if accountID == "" || clientID == "" || clientSecret == "" || wsURL == "" {
		log.Fatalf("Missing required environment variables: ACCOUNT_ID, CLIENT_ID, CLIENT_SECRET, WS_URL")
	}

	service := websocketservice.NewWebSocketService(accountID, clientID, clientSecret, wsURL)

	accessToken, err := service.GetAccessToken()
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
	}

	err = service.ConnectWebSocket(accessToken)
	if err != nil {
		log.Fatalf("Error connecting to WebSocket: %v", err)
	}
	defer service.CloseWebSocket()

	service.SendHeartbeat()

	// Main application loop
	for {
		_, message, err := service.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		fmt.Printf("Received message: %s\n", message)

		var data map[string]interface{}
		err = json.Unmarshal([]byte(message), &data)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}

		content, ok := data["content"].(string)
		if !ok {
			continue
		}

		var contentData map[string]interface{}
		err = json.Unmarshal([]byte(content), &contentData)
		if err != nil {
			continue
		}

		event, ok := contentData["event"].(string)
		if !ok {
			continue
		} else {
			log.Printf("Event: %s", event)
			if event == "user.created" {
				newUserCreatedHandler()
			}
		}
	}

	log.Println("WebSocket connection closed.")
}
