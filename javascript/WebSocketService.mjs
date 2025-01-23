import WebSocket from "ws";
import { createHmac } from "node:crypto";

export class WebSocketService {
  constructor(accountId, clientId, clientSecret, url) {
    this.url = url;
    this.accountId = accountId;
    this.clientId = clientId;
    this.clientSecret = clientSecret;
    this.oauthUrl = "https://zoom.us/oauth/";
    this.ws = null;
    this.heartbeat = null;
  }

  async getAccessToken() {
    try {
      const oauthToken = Buffer.from(
        `${this.clientId}:${this.clientSecret}`
      ).toString("base64");

      const res = await fetch(
        `${this.oauthUrl}token?grant_type=account_credentials&account_id=${this.accountId}`,
        {
          method: "POST",
          headers: {
            Authorization: `Basic ${oauthToken}`,
          },
        }
      );

      const data = await res.json();
      return data.access_token;
    } catch (e) {
      console.error(
        `${e.message}\nUnable to get access token. Check credentials.`
      );
    }
  }

  sendHeartbeat() {
    this.ws.send(JSON.stringify({ module: "heartbeat" }));
    this.heartbeat = setTimeout(this.sendHeartbeat, 30000 * Math.random());
  }

  open() {
    console.log("connected");
    this.sendHeartbeat;
  }

  close() {
    console.log("connection closed");
    clearTimeout(this.heartbeat);
    this.ws = null;
  }

  async message(message) {
    let data = {};

    try {
      data = JSON.parse(message);
    } catch (e) {
      console.error("Invalid JSON received");
      return;
    }

    if (data?.module === "message") {
      console.log(JSON.parse(data.content));
      if (data?.content?.event === "user.created") this.newUserCreatedHandler;
    }
  }

  newUserCreatedHandler() {
    console.log("\n\n\n\nA new user was created");
    console.log("Do some processing\n\n\n\n");
  }

  async connect() {
    try {
      const accessToken = await this.getAccessToken();
      if (!accessToken) throw new Error("Unable to get access token");
      this.ws = new WebSocket(`${this.url}&access_token=${accessToken}`);
      this.ws.on("open", this.open);
      this.ws.on("close", this.close);
      this.ws.on("message", this.message);
    } catch (e) {
      console.error(e);
    }
  }
}
