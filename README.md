# Zoom WebSocket Listener (examples)

> **Note**
> 
> The following sample application is a personal, open-source project shared by the app creator and not an officially supported Zoom Video Communications, Inc. sample application. Zoom Video Communications, Inc., its employees and affiliates are not responsible for the use and maintenance of this application. Please use this sample application for inspiration, exploration and experimentation at your own risk and enjoyment. You may reach out to the app creator and broader Zoom Developer community on https://devforum.zoom.us/ for technical discussion and assistance, but understand there is no service level agreement support for this application. Thank you and happy coding!

This project is a collection of **Zoom WebSocket listeners** implemented in **JavaScript**, **Python**, and **C# .NET Core**. Each application performs the following tasks:

1. **Fetches an Access Token**: Authenticates using OAuth credentials to retrieve an access token for connecting to Zoom WebSockets.
2. **Connects to a WebSocket**: Establishes a real-time communication channel with Zoom's WebSocket service.
3. **Sends Heartbeats**: Maintains the connection by periodically sending heartbeat messages.
4. **Custom Event Handling**: Listens for the `user.created` event and triggers a customizable function to handle this event.

## Features
- **Modular WebSocketService Class**: Each implementation contains a `WebSocketService` class that manages authentication, connection, and message handling.
- **Customizable Functionality**: Modify the `newUserCreatedHandler` function in each `WebSocketService` to perform specific actions when a new user is created.
- **Multi-language Support**: Choose your preferred language from JavaScript, Python, or C# .NET Core to integrate with Zoom WebSockets.

---

## File Structure
The project includes:

```
project/
├── .env.local               # Environment variables
├── javascript/               # JavaScript implementation
│   ├── WebSocketService.mjs  # Contains the WebSocketService class
│   └── index.mjs             # Entry point script
├── python/                   # Python implementation
│   ├── WebSocketService.py   # Contains the WebSocketService class
│   └── index.py              # Entry point script
├── csharp/                   # C# .NET Core implementation
│   ├── WebSocketService.cs   # Contains the WebSocketService class
│   └── Program.cs            # Entry point script
```

---

## Setup and Usage

### 1. Environment Variables
Each application expects environment variables to be defined in the `.env.local` file located in the project root:

```dotenv
accountId=your_account_id
clientId=your_client_id
clientSecret=your_client_secret
url=your_websocket_url
```

Ensure the `.env.local` file is correctly placed and configured.

### 2. Installation

#### JavaScript
Navigate to the `javascript` directory, install dependencies, and run the application:
```bash
cd javascript
npm install
node index.mjs
```

#### Python
Navigate to the `python` directory, install dependencies, and run the application:
```bash
cd python
pip install -r requirements.txt
python index.py
```

#### C# .NET Core
Navigate to the `csharp` directory, restore dependencies, and run the application:
```bash
cd csharp
dotnet restore
dotnet run
```

---

## Customizing the Event Handler
The `newUserCreatedHandler` function in the `WebSocketService` class is designed to be modified. This function is triggered whenever a `user.created` event is received.

### Example (Python):
```python
def new_user_created_handler(self):
    print("\n\nA new user was created")
    print("Perform custom processing here\n\n")
```

### Example (C#):
```csharp
private void NewUserCreatedHandler()
{
    Console.WriteLine("\n\nA new user was created");
    Console.WriteLine("Perform custom processing here\n\n");
}
```

### Example (JavaScript):
```javascript
newUserCreatedHandler() {
    console.log("\n\nA new user was created");
    console.log("Perform custom processing here\n\n");
}
```

## Installation Links
Below are the links to download and install the required tools for each implementation:

- **Node.js (JavaScript)**: [Download Node.js](https://nodejs.org/)
- **Python**: [Download Python](https://www.python.org/downloads/)
- **.NET Core (C#)**: [Download .NET Core](https://dotnet.microsoft.com/download)


