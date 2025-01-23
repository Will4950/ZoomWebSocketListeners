use base64::prelude::*;
use futures_util::SinkExt;
use futures_util::StreamExt;
use reqwest::Client;
use serde_json::Value;
use std::time::Duration;
use tokio::sync::mpsc;
use tokio::time;
use tokio_tungstenite::connect_async;
use tungstenite::protocol::Message;

pub struct WebSocketService {
    account_id: String,
    client_id: String,
    client_secret: String,
    url: String,
    oauth_url: String,
}

impl WebSocketService {
    pub fn new(account_id: String, client_id: String, client_secret: String, url: String) -> Self {
        Self {
            account_id,
            client_id,
            client_secret,
            url,
            oauth_url: "https://zoom.us/oauth/".to_string(),
        }
    }

    pub async fn get_access_token(&self) -> Option<String> {
        let oauth_token =
            BASE64_STANDARD.encode(format!("{}:{}", self.client_id, self.client_secret));
        let client = Client::new();

        let response = client
            .post(format!(
                "{}token?grant_type=account_credentials&account_id={}",
                self.oauth_url, self.account_id
            ))
            .header("Authorization", format!("Basic {}", oauth_token))
            .send()
            .await;

        match response {
            Ok(resp) => {
                if let Ok(json) = resp.json::<Value>().await {
                    return json["access_token"].as_str().map(|s| s.to_string());
                }
            }
            Err(e) => {
                eprintln!("Error fetching access token: {}", e);
            }
        }

        None
    }

    pub async fn connect(&self) {
        if let Some(access_token) = self.get_access_token().await {
            let ws_url = format!("{}&access_token={}", self.url, access_token);
            match connect_async(&ws_url).await {
                Ok((websocket, _)) => {
                    println!("Connected to WebSocket.");

                    // Split the websocket into sender and receiver
                    let (mut ws_sender, mut ws_receiver) = websocket.split();

                    let (tx, mut rx) = mpsc::channel::<Message>(32);
                    let tx_clone = tx.clone();

                    // Spawn task for sending messages
                    tokio::spawn(async move {
                        while let Some(msg) = rx.recv().await {
                            ws_sender.send(msg).await.unwrap();
                        }
                    });

                    self.send_heartbeat(tx_clone);

                    // Use receiver in the main loop
                    while let Some(msg) = ws_receiver.next().await {
                        match msg {
                            Ok(Message::Text(text)) => self.message_handler(&text).await,
                            Ok(_) => {}
                            Err(e) => {
                                eprintln!("WebSocket error: {}", e);
                                break;
                            }
                        }
                    }

                    println!("Connection closed.");
                }
                Err(e) => {
                    eprintln!("Failed to connect to WebSocket: {}", e);
                }
            }
        } else {
            eprintln!("Unable to get access token.");
        }
    }

    fn send_heartbeat(&self, tx: mpsc::Sender<Message>) {
        tokio::spawn(async move {
            loop {
                time::sleep(Duration::from_secs(30)).await;
                if tx
                    .send(Message::Text("{\"module\": \"heartbeat\"}".to_string()))
                    .await
                    .is_err()
                {
                    break;
                }
            }
        });
    }

    pub async fn message_handler(&self, message: &str) {
        match serde_json::from_str::<Value>(message) {
            Ok(data) => {
                if let Some(module) = data["module"].as_str() {
                    if module == "message" {
                        if let Ok(pretty_json) = serde_json::to_string_pretty(&data["content"]) {
                            println!("Received content:\n{}", pretty_json);
                        }
                        if let Some(content) = data["content"].as_object() {
                            if let Some(event) = content["event"].as_str() {
                                if event == "user.created" {
                                    self.new_user_created_handler().await;
                                }
                            }
                        }
                    }
                }
            }
            Err(_) => {
                eprintln!("Invalid JSON received.");
            }
        }
    }

    async fn new_user_created_handler(&self) {
        println!("\n\nA new user was created.");
        println!("Perform custom processing here\n\n");
    }
}
