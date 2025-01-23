mod websocket_service;
use std::env;
use std::path::PathBuf;
use websocket_service::WebSocketService;

#[tokio::main]
async fn main() {
    let mut env_path = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    env_path.pop();
    env_path.push(".env.local");

    dotenv::from_path(env_path).ok();

    let account_id = env::var("accountId").expect("ACCOUNT_ID must be set");
    let client_id = env::var("clientId").expect("CLIENT_ID must be set");
    let client_secret = env::var("clientSecret").expect("CLIENT_SECRET must be set");
    let url = env::var("url").expect("WS_URL must be set");

    let service = WebSocketService::new(account_id, client_id, client_secret, url);

    service.connect().await;
}
