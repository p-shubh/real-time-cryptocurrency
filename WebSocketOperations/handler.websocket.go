package webSocketOperations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"trade/pkg/db"
	"trade/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CoinMarketHistory(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()

	internalConn, _, err := websocket.DefaultDialer.Dial(os.Getenv("WEBSOCKET_URL"), nil)
	if err != nil {
		log.Println("Failed to connect to internal websocket:", err)
		return
	}
	defer internalConn.Close()

	coin := os.Getenv("COIN")
	if coin != "ETH" {
		coin = "ETH"
	} else {
		coin = "BTC"
	}

	payload := map[string]string{"coin": coin}
	if err := internalConn.WriteJSON(payload); err != nil {
		log.Println("Failed to send payload to internal websocket:", err)
		return
	}
	var currentTime time.Time

	messageChan := make(chan []byte)

	go func() {
		defer close(messageChan)
		for {
			_, message, err := internalConn.ReadMessage()
			if err != nil {
				log.Println("Error reading message from internal websocket:", err)
				return
			}
			messageChan <- message
		}
	}()

	for message := range messageChan {

		var datas []models.CryptoHistory
		if err := json.Unmarshal(message, &datas); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		for _, data := range datas {

			timestamp, err := time.Parse("2006-01-02 15:04:05", data.Timestamp)
			if err != nil {
				log.Println("Error parsing timestamp:", err)
				continue
			}

			// Check if the timestamp is greater than the current timestamp
			if timestamp.After(currentTime) {
				currentTime = timestamp

				// add insert db opearation
				InsertCryptoHistory(data)
			}

			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error writing message to external websocket:", err)
				return
			}
		}
	}
}
func ManualTimeInsterCoinMarketHistory(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()

	internalConn, _, err := websocket.DefaultDialer.Dial(os.Getenv("WEBSOCKET_URL"), nil)
	if err != nil {
		log.Println("Failed to connect to internal websocket:", err)
		return
	}
	defer internalConn.Close()

	payload := map[string]string{"coin": os.Getenv("COIN")}
	if err := internalConn.WriteJSON(payload); err != nil {
		log.Println("Failed to send payload to internal websocket:", err)
		return
	}
	var currentTime time.Time

	messageChan := make(chan []byte)

	go func() {
		defer close(messageChan)
		for {
			_, message, err := internalConn.ReadMessage()
			if err != nil {
				log.Println("Error reading message from internal websocket:", err)
				return
			}
			messageChan <- message
		}
	}()

	for message := range messageChan {

		var payloads []models.CryptoHistory
		if err := json.Unmarshal(message, &payloads); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		for _, payload := range payloads {

			timestamp := currentTime.Add(time.Second * 5)

			// Check if the timestamp is greater than the current timestamp
			if timestamp.After(currentTime) {
				currentTime = timestamp
				payload.Timestamp = currentTime.Format("2006-01-02 15:04:05")

				// Insert the payload record in the database
				if err := db.DB.Create(&payload).Error; err != nil {
					log.Println("Error saving to database:", err)
				}
			}

			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error writing message to external websocket:", err)
				return
			}
			time.Sleep(time.Second)
		}
	}
}

func WebSocketForGetCoinMarketHistory(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()

	for {

		// Example: Send data every 1 seconds
		ticker := time.NewTicker(1 * time.Second)
		count := 0
		var (
			data                          models.CryptoHistory
			Timestamp, formattedTimestamp time.Time
			err1                          error
		)
		// defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if count == 0 {
					data = GetWebSocketDbOperationForEvery5minFirstTimestamp()
					formattedTimestamp, err1 = time.Parse(time.RFC3339, data.Timestamp)
					if err1 != nil {
						log.Println("Error parsing timestamp:", err1)
					}
					formattedTimestampp := formattedTimestamp.Format("2006-01-02 15:04:05")

					Timestamp, err = time.Parse("2006-01-02 15:04:05", formattedTimestampp)
					if err != nil {
						log.Println("Error converting formatted timestamp to time.Time:", err)
						return
					}
					fmt.Println("formattedTimestampp : ", Timestamp)

					count++
				}

				startTime := Timestamp
				endTime := startTime.Add(+time.Minute * 5)
				data := GetWebSocketDbOperationForEvery5min(startTime, endTime)

				reqBodyBytes := new(bytes.Buffer)
				json.NewEncoder(reqBodyBytes).Encode(data)

				conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

				if err := conn.WriteMessage(websocket.TextMessage, reqBodyBytes.Bytes()); err != nil {
					fmt.Println("Failed to write the response messages : ", err.Error())
				}
				Timestamp = endTime
			}
		}
	}
}

func ApiForGetCoinMarketHistory(c *gin.Context) {

	// Prepare the response
	data := GetCryptoHistoryCalculation()

	c.JSON(200, gin.H{
		"data": data,
	})
}
