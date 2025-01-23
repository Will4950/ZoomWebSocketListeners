import os
from dotenv import load_dotenv
from WebSocketService import WebSocketService

load_dotenv("../.env.local")
account_id = os.getenv("accountId")
client_id = os.getenv("clientId")
client_secret = os.getenv("clientSecret")
url = os.getenv("url")

if not all([account_id, client_id, client_secret, url]):
    raise ValueError("Missing one or more required environment variables.")

ws = WebSocketService(account_id, client_id, client_secret, url)
ws.connect()
