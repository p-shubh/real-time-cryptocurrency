package router

import (
	"log"
	webSocketOperations "trade/WebSocketOperations"

	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	WebSocketRouting(r)
	if err := r.Run(":9090"); err != nil {
		log.Println()
		return
	}
}

func WebSocketRouting(c *gin.Engine) {

	//Websockets
	c.GET("/insertCoinMarketHistory", webSocketOperations.CoinMarketHistory)
	c.GET("/manualInsertCoinMarketHistory", webSocketOperations.ManualTimeInsterCoinMarketHistory)
	c.GET("/WebSocketForCalGetCoinMarketHistory", webSocketOperations.WebSocketForGetCoinMarketHistory)
	//REST API
	c.GET("/ApiGetCoinMarketHistory", webSocketOperations.ApiForGetCoinMarketHistory)

}
