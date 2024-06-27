package webSocketOperations

import (
	"fmt"
	"log"
	"time"
	"trade/pkg/db"
	"trade/pkg/models"
)

func InsertCryptoHistory(data models.CryptoHistory) {
	// Insert the data record in the database
	if err := db.DB.Create(&data).Error; err != nil {
		log.Println("Error saving to database:", err)
	}
}

func GetWebSocketDbOperation() models.RateStats {
	endTime := time.Now()
	startTime := endTime.Add(-10 * time.Minute)

	var data models.RateStats
	result := db.DB.Raw("SELECT MAX(rate) AS max_rate, MIN(rate) AS min_rate, SUM(volume) AS sum_volume FROM crypto_histories where timestamp >= '" + startTime.Format("2006-01-02 15:04:05") + "' AND timestamp <= '" + endTime.Format("2006-01-02 15:04:05") + "'").Debug().Scan(&data)
	if result.Error != nil {
		// handle error print the error
		fmt.Println("GetWebSocketDbOperation : ", result.Error)
		return models.RateStats{}
	}

	return data

}

func GetWebSocketDbOperationForEvery5min(startTime, endTime time.Time) models.RateStats {

	var data models.RateStats

	result := db.DB.Raw("SELECT MAX(rate) AS max_rate, MIN(rate) AS min_rate, SUM(volume) AS sum_volume FROM crypto_histories where timestamp >= '" + startTime.Format("2006-01-02 15:04:05") + "' AND timestamp <= '" + endTime.Format("2006-01-02 15:04:05") + "'").Debug().Scan(&data)
	if result.Error != nil {
		// handle error print the error
		fmt.Println("GetWebSocketDbOperation : ", result.Error)
		return models.RateStats{}
	}

	return data

}

// get the first timestamp
func GetWebSocketDbOperationForEvery5minFirstTimestamp() models.CryptoHistory {

	var firstRecord models.CryptoHistory
	result := db.DB.Order("timestamp ASC").First(&firstRecord)
	if result.Error != nil {
		// handle error print the error
		fmt.Println("GetWebSocketDbOperation : ", result.Error)

	}
	return firstRecord
}
func GetCryptoHistoryCalculation() []models.RateStats {
	var crypto_histories []models.CryptoHistory
	var rateStats []models.RateStats

	result := db.DB.Raw("SELECT * FROM crypto_histories ORDER BY timestamp ASC").Scan(&crypto_histories)
	if result.Error != nil {
		fmt.Println("Error retrieving crypto histories:", result.Error)
		return []models.RateStats{}
	}

	windowSize := 5 * time.Minute
	numRecords := len(crypto_histories)
	if numRecords == 0 {
		return rateStats
	}

	// Initialize window start and end times
	startTimRaw := crypto_histories[0].Timestamp

	formattedTimestamp, err1 := time.Parse(time.RFC3339, startTimRaw)
	if err1 != nil {
		log.Println("Error parsing timestamp:", err1)
	}
	formattedTimestampp := formattedTimestamp.Format("2006-01-02 15:04:05")

	startTime, err := time.Parse("2006-01-02 15:04:05", formattedTimestampp)
	if err != nil {
		log.Println("Error converting formatted timestamp to time.Time:", err)
	}

	endTime := startTime.Add(windowSize)

	// Initialize window statistics
	var maxRate, minRate, sumVolume float64
	count := 0

	for i := 0; i < numRecords; i++ {
		recordTimeString := crypto_histories[i].Timestamp

		formattedTimestamp, err1 := time.Parse(time.RFC3339, recordTimeString)
		if err1 != nil {
			log.Println("Error parsing timestamp:", err1)
		}
		recordTimeTimestampp := formattedTimestamp.Format("2006-01-02 15:04:05")

		recordTime, err := time.Parse("2006-01-02 15:04:05", recordTimeTimestampp)
		if err != nil {
			log.Println("Error converting formatted timestamp to time.Time:", err)
		}

		// Check if current record is within the current window
		for recordTime.After(endTime) {

			// Append statistics for the completed window
			rateStats = append(rateStats, models.RateStats{
				MaxRate:   maxRate,
				MinRate:   minRate,
				SumVolume: sumVolume,
			})

			// Move window to the next 5-minute interval
			startTime = endTime
			endTime = startTime.Add(windowSize)

			// Reset window statistics for the new window
			maxRate = 0
			minRate = 0
			sumVolume = 0
			count = 0
		}

		// Process current record within the current window
		data := crypto_histories[i]
		if count == 0 {
			maxRate = data.Rate
			minRate = data.Rate
		} else {
			if data.Rate > maxRate {
				maxRate = data.Rate
			}
			if data.Rate < minRate {
				minRate = data.Rate
			}
		}
		sumVolume += data.Volume
		count++
	}

	
	rateStats = append(rateStats, models.RateStats{
		MaxRate:   maxRate,
		MinRate:   minRate,
		SumVolume: sumVolume,
	})

	return rateStats
}
