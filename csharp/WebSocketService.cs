using System;
using System.Net.Http;
using System.Net.WebSockets;
using System.Text;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;

public class WebSocketService
{
	private readonly string accountId;
	private readonly string clientId;
	private readonly string clientSecret;
	private readonly string url;
	private readonly string oauthUrl = "https://zoom.us/oauth/";
	private ClientWebSocket ws;
	private CancellationTokenSource heartbeatTokenSource;

	public WebSocketService(string accountId, string clientId, string clientSecret, string url)
	{
		this.accountId = accountId;
		this.clientId = clientId;
		this.clientSecret = clientSecret;
		this.url = url;
	}

	private async Task<string> GetAccessTokenAsync()
	{
		try
		{
			var authToken = Convert.ToBase64String(Encoding.UTF8.GetBytes($"{clientId}:{clientSecret}"));
			using var client = new HttpClient();
			client.DefaultRequestHeaders.Add("Authorization", $"Basic {authToken}");

			var response = await client.PostAsync(
				$"{oauthUrl}token?grant_type=account_credentials&account_id={accountId}",
				null);

			response.EnsureSuccessStatusCode();
			var content = await response.Content.ReadAsStringAsync();
			var json = JsonSerializer.Deserialize<JsonElement>(content);

			return json.GetProperty("access_token").GetString();
		}
		catch (Exception ex)
		{
			Console.WriteLine($"Error retrieving access token: {ex.Message}");
			throw;
		}
	}

	private async Task SendHeartbeatAsync()
	{
		while (!heartbeatTokenSource.Token.IsCancellationRequested)
		{
			var heartbeatMessage = JsonSerializer.Serialize(new { module = "heartbeat" });
			var buffer = Encoding.UTF8.GetBytes(heartbeatMessage);

			await ws.SendAsync(new ArraySegment<byte>(buffer), WebSocketMessageType.Text, true, CancellationToken.None);

			await Task.Delay(TimeSpan.FromSeconds(30), heartbeatTokenSource.Token);
		}
	}

	private void OnMessageReceived(string message)
	{
		try
		{
			var dataObj = JsonSerializer.Deserialize<JsonElement>(message);

			if (dataObj.TryGetProperty("module", out var module) && module.GetString() == "message")
			{
				var contentString = dataObj.GetProperty("content").GetString();
				var contentObj = JsonSerializer.Deserialize<JsonElement>(contentString);

				Console.WriteLine(JsonSerializer.Serialize(
					contentObj,
					new JsonSerializerOptions { WriteIndented = true }
				));

				if (contentObj.GetProperty("event").GetString() == "user.created")
				{
					NewUserCreatedHandler();
				}
			}
		}
		catch (Exception ex)
		{
			Console.WriteLine($"Invalid JSON received: {ex.Message}");
		}
	}

	private void NewUserCreatedHandler()
	{
		Console.WriteLine("\n\nA new user was created");
		Console.WriteLine("Do some processing\n\n");
	}

	public async Task ConnectAsync()
	{
		try
		{
			var accessToken = await GetAccessTokenAsync();

			ws = new ClientWebSocket();
			await ws.ConnectAsync(new Uri($"{url}&access_token={accessToken}"), CancellationToken.None);

			Console.WriteLine("Connected");

			heartbeatTokenSource = new CancellationTokenSource();
			_ = Task.Run(SendHeartbeatAsync);

			var buffer = new byte[1024 * 4];

			while (ws.State == WebSocketState.Open)
			{
				var result = await ws.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);

				if (result.MessageType == WebSocketMessageType.Text)
				{
					var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
					OnMessageReceived(message);
				}
				else if (result.MessageType == WebSocketMessageType.Close)
				{
					Console.WriteLine("Connection closed");
					await ws.CloseAsync(WebSocketCloseStatus.NormalClosure, string.Empty, CancellationToken.None);
				}
			}
		}
		catch (Exception ex)
		{
			Console.WriteLine($"Error: {ex.Message}");
		}
	}
}
