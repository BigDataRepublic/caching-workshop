package main

import (
	"caching-workshop/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"

	"log"
)

var (
	ListenAddr  = "localhost:8080"
	RedisAddr   = "localhost:6379"
	RedisPasswd = os.Getenv("REDIS_PASSWORD")
)

func initRouter(database *db.Database) *gin.Engine {
	r := gin.Default()

	r.POST("/points", func(c *gin.Context) {
		var userJson db.User
		if err := c.ShouldBindJSON(&userJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := database.SaveUser(&userJson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": userJson})
	})

	r.GET("/points/:username", func(c *gin.Context) {
		username := c.Param("username")
		user, err := database.GetUser(username)
		if err != nil {
			if err == db.ErrNil {
				c.JSON(http.StatusNotFound, gin.H{"error": "No record found for " + username})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	r.GET("/leaderboard", func(c *gin.Context) {
		leaderboard, err := database.GetLeaderboard()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
	})

	r.GET("/rival/:username", func(c *gin.Context) {
		username := c.Param("username")
		user, err := database.GetRival(username)
		if err != nil {
			if err == db.ErrNil {
				c.JSON(http.StatusNotFound, gin.H{"error": "No record found for " + username})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	return r
}

func main() {
	database, err := db.NewDatabase(RedisAddr, RedisPasswd)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	router := initRouter(database)
	router.Run(ListenAddr)
}
