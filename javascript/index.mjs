import { WebSocketService } from "./WebSocketService.mjs";
import path from "path";

import { config } from "dotenv";
const dotenvAbsolutePath = path.join("../.env.local");
const dotenv = config({
  path: dotenvAbsolutePath,
});
if (dotenv.error) {
  throw dotenv.error;
}

const ws = new WebSocketService(
  process.env.accountId,
  process.env.clientId,
  process.env.clientSecret,
  process.env.url
);

await ws.connect();
