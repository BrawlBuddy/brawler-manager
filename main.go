package main

import (
	"brawler-manager/brawlers"
	"brawler-manager/calculator"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
)

type matchContext struct {
	Bans     []string `json:"bans"`
	Friendly []string `json:"friendly"`
	Enemy    []string `json:"enemy"`
	Map      string   `json:"map"`
}

var OneVone map[string]float32
var mapData map[string]map[string]float32
var useRates map[string]map[string]float32

func postBrawlerPicks(c *gin.Context) {
	var match matchContext
	if err := c.BindJSON(&match); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	brawlersRanked := calculator.GenerateRanks(match.Bans, match.Friendly, match.Enemy, match.Map, OneVone, mapData, useRates)
	c.IndentedJSON(http.StatusOK, brawlersRanked)
}

func returnRandomBrawlerTest(match matchContext) []brawlers.Brawler {
	brawlersList := calculator.CreatePool(brawlers.GetAllBrawlers(), match.Bans, match.Friendly, match.Enemy)
	result := []brawlers.Brawler{}
	for _, x := range brawlersList {
		result = append(result, brawlers.Brawler{Name: x, WinPct: rand.Float32()})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].WinPct > result[j].WinPct // sort in decreasing order
	})
	return result
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://brawlbuddy.github.io")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	gin.SetMode("release")
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.POST("/brawlerpicks", postBrawlerPicks)

	mapData = brawlers.GetMapData()
	OneVone = brawlers.GetMatchUps()
	useRates = brawlers.GetUseRates()
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
