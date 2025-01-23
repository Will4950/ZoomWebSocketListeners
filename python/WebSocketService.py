import base64
import json
import requests
import threading
from websocket import WebSocketApp

class WebSocketService:
    def __init__(self, account_id, client_id, client_secret, url):
        self.url = url
        self.account_id = account_id
        self.client_id = client_id
        self.client_secret = client_secret
        self.oauth_url = "https://zoom.us/oauth/"
        self.ws = None
        self.heartbeat_thread = None

    def get_access_token(self):
        try:
            oauth_token = base64.b64encode(f"{self.client_id}:{self.client_secret}".encode()).decode()
            headers = {"Authorization": f"Basic {oauth_token}"}
            response = requests.post(
                f"{self.oauth_url}token?grant_type=account_credentials&account_id={self.account_id}",
                headers=headers
            )
            response.raise_for_status()
            data = response.json()
            return data.get("access_token")
        except Exception as e:
            print(f"{e}\nUnable to get access token. Check credentials.")
            return None

    def send_heartbeat(self):
        if self.ws:
            self.ws.send(json.dumps({"module": "heartbeat"}))
        self.heartbeat_thread = threading.Timer(30, self.send_heartbeat)
        self.heartbeat_thread.start()

    def on_open(self, ws):
        print("Connected")
        self.send_heartbeat()

    def on_close(self, ws, close_status_code, close_msg):
        print("Connection closed")
        if self.heartbeat_thread:
            self.heartbeat_thread.cancel()

    def on_message(self, ws, message):
        try:
            data = json.loads(message)
        except json.JSONDecodeError:
            print("Invalid JSON received")
            return

        if data.get("module") == "message":
            print(json.dumps(data.get("content"), indent=4))
            if data.get("content", {}).get("event") == "user.created":
                self.new_user_created_handler()

    def new_user_created_handler(self):
        print("\n\nA new user was created")
        print("Do some processing\n\n")

    def connect(self):
        try:
            access_token = self.get_access_token()
            if not access_token:
                raise Exception("Unable to get access token")

            self.ws = WebSocketApp(
                f"{self.url}&access_token={access_token}",
                on_open=self.on_open,
                on_close=self.on_close,
                on_message=self.on_message
            )
            self.ws.run_forever()
        except Exception as e:
            print(e)