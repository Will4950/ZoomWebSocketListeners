using System;
using System.IO;
using System.Threading.Tasks;
using dotenv.net;

class Program
{
	static async Task Main(string[] args)
	{
		DotEnv.Load(new DotEnvOptions().WithEnvFiles("../.env.local"));

		string accountId = Environment.GetEnvironmentVariable("accountId");
		string clientId = Environment.GetEnvironmentVariable("clientId");
		string clientSecret = Environment.GetEnvironmentVariable("clientSecret");
		string url = Environment.GetEnvironmentVariable("url");

		if (string.IsNullOrWhiteSpace(accountId) || string.IsNullOrWhiteSpace(clientId) ||
			string.IsNullOrWhiteSpace(clientSecret) || string.IsNullOrWhiteSpace(url))
		{
			Console.WriteLine("Missing one or more required environment variables.");
			return;
		}

		var wsService = new WebSocketService(accountId, clientId, clientSecret, url);

		await wsService.ConnectAsync();
	}
}