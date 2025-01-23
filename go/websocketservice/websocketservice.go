package websocketservice

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	AccountID    string
	ClientID     string
	ClientSecret string
	URL          string
	OAuthURL     string
	Conn         *websocket.Conn
	Heartbeat    *time.Ticker
}

func NewWebSocketService(accountID, clientID, clientSecret, url string) *WebSocketService {
	return &WebSocketService{
		AccountID:    accountID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		URL:          url,
		OAuthURL:     "https://zoom.us/oauth/",
	}
}

func (service *WebSocketService) GetAccessToken() (string, error) {
	oauthToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", service.ClientID, service.ClientSecret)))

	url := fmt.Sprintf("%stoken?grant_type=account_credentials&account_id=%s", service.OAuthURL, service.AccountID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", oauthToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response")
	}

	return accessToken, nil
}

func (service *WebSocketService) ConnectWebSocket(accessToken string) error {

	wsURL := fmt.Sprintf("%s&access_token=%s", service.URL, accessToken)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %v", err)
	}

	service.Conn = conn
	return nil
}

func (service *WebSocketService) SendHeartbeat() {
	service.Heartbeat = time.NewTicker(30 * time.Second)
	go func() {
		for range service.Heartbeat.C {
			if err := service.Conn.WriteMessage(websocket.TextMessage, []byte("heartbeat")); err != nil {
				log.Println("Failed to send heartbeat:", err)
				service.Heartbeat.Stop()
				return
			}
		}
	}()
}

func (service *WebSocketService) CloseWebSocket() {
	if service.Heartbeat != nil {
		service.Heartbeat.Stop()
	}
	if service.Conn != nil {
		service.Conn.Close()
	}
}
