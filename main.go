package main

import (
	"brawler-manager/calculator"
	"github.com/gin-gonic/gin"
	"net/http"
)

type matchContext struct {
	Bans     []string `json:"bans"`
	Friendly []string `json:"friendly"`
	Enemy    []string `json:"enemy"`
	Map      string   `json:"map"`
}

func getBrawlerPicks(c *gin.Context) {
	var match matchContext
	if err := c.BindJSON(&match); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	brawlersRanked := calculator.GenerateRanks(match.Bans, match.Friendly, match.Enemy, match.Map)
	c.IndentedJSON(http.StatusOK, brawlersRanked)
}

func main() {
	router := gin.Default()
	router.GET("/brawlerpicks", getBrawlerPicks)

	router.Run("localhost:8080")
}
