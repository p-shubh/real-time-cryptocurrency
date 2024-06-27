
# Cryptocurrency Data Retrieval and API

This project involves developing a Go application using the Gin or Fiber framework to retrieve real-time cryptocurrency data for BTC, USDT, and ETH from a WebSocket, store this data in a PostgreSQL database, and provide API endpoints for retrieving top volumes, losers, and gainers for the last 24 hours at 5-minute intervals.

## Requirements

1. **WebSocket Data Retrieval:**

   - Connect to a WebSocket providing real-time data for BTC, USDT, and ETH.
   - Retrieve the data every millisecond.
   - Ensure that only unique timestamp data is stored.
2. **Database Storage:**

   - Use PostgreSQL to store the data.
   - Maintain only the last 24 hours of data.
3. **Batching Insert:**

   - Implement batching insert for the WebSocket data to improve performance.
4. **API Endpoints:**

   - Create a REST API or WebSocket API using the Gin or Fiber framework.
   - Provide endpoints to get the top volumes, losers, and gainers for the last 24 hours at 5-minute intervals.
   - Calculate gainers and losers based on the rate from the WebSocket data.
5. **Environment Configuration:**

   - Use a `.env` file to store database credentials and other configuration settings.

## Deliverables

1. **Source Code:**

   - Well-documented and structured source code in a zip file or in a GitHub link.
   - Include a README file with setup instructions.
2. **.env File:**

   - Template for the `.env` file with placeholders for database credentials.

## Example .env File

```env
POSTGRES_USER=your_postgres_user
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_DB=your_postgres_db
POSTGRES_HOST=your_postgres_host
POSTGRES_PORT=your_postgres_port
WEBSOCKET_URL=your_websocket_url
```

## Evaluation Criteria

- Correctness and efficiency of WebSocket data retrieval and storage.
- Proper use of PostgreSQL.
- Functionality and performance of the API endpoints.
- Code quality and documentation.

## WebSocket URL

Use the following WebSocket URL to get the coin data:

```
wss://stream.bit24hr.in/coin_market_history/
```

You need to pass the following message in the WebSocket URL to get the data for specific coins:

```json
{ 
 "coin": "coin_name" 
}
```

Replace `coin_name` with one of the following: `"BTC"`, `"ETH"`, `"USDT"`.

## Setup Instructions

### Prerequisites

- Go (version 1.18 or higher)
- PostgreSQL
- Git

### Installation

1. **Clone the Repository:**

   ```sh
   git clone <repository-url>
   cd <repository-directory>
   ```
2. **Setup Environment Variables:**

   Create a `.env` file in the root directory and add your configuration settings based on the provided `.env` template.
3. **Install Dependencies:**

   ```sh
   go mod tidy
   ```
4. **Run the Application:**

   ```sh
   go run main.go
   ```

## API Endpoints

### WebSocket Routing

```go
func WebSocketRouting(c *gin.Engine) {
    // Websockets
    c.GET("/insertCoinMarketHistory", webSocketOperations.CoinMarketHistory)
    c.GET("/manualInsertCoinMarketHistory", webSocketOperations.ManualTimeInsertCoinMarketHistory)
    c.GET("/WebSocketForCalGetCoinMarketHistory", webSocketOperations.WebSocketForGetCoinMarketHistory)

    // REST API
    c.GET("/ApiGetCoinMarketHistory", webSocketOperations.ApiForGetCoinMarketHistory)

    // Health Check
    c.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
}
```

### Example API Requests

- **Insert Coin Market History:**

  ```
  GET /insertCoinMarketHistory
  ```
- **Manual Insert Coin Market History:**

  ```
  GET /manualInsertCoinMarketHistory
  ```
- **WebSocket for Coin Market History:**

  ```
  GET /WebSocketForCalGetCoinMarketHistory
  ```
- **Get Coin Market History via API:**

  ```
  GET /ApiGetCoinMarketHistory
  ```
- **Health Check:**

  ```
  GET /ping
  ```

## Conclusion

This README provides a comprehensive guide for setting up and running the Go application for real-time cryptocurrency data retrieval and API endpoint creation. Make sure to configure your environment variables correctly and follow the setup instructions to get the application running.
