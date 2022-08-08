package main

import (
	"caching-workshop/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"

	"log"
)

var (
	ListenAddr  = "localhost:8080"
	RedisAddr   = "localhost:6379"
	RedisPasswd = os.Getenv("REDIS_PASSWORD")
	FastRivals  = strings.ToLower(os.Getenv("FAST_RIVALS"))
)

func initRouter(database *db.Database) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/points", func(c *gin.Context) {
		var userJson db.User
		if err := c.ShouldBindJSON(&userJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := database.SaveUser(&userJson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
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

func backgroundTasks(database *db.Database) {
	for range time.Tick(time.Second * 10) {
		err := database.UpdateRivals()
		if err != nil {
			log.Printf("Background task failed: %s", err)
		}
	}
}

func main() {
	database, err := db.NewDatabase(RedisAddr, RedisPasswd)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	if FastRivals == "true" {
		log.Println("Fast rivals enabled, will run background process to update rival cache")
		go backgroundTasks(database)
	}

	router := initRouter(database)
	err = router.Run(ListenAddr)
	if err != nil {
		log.Fatalf("Router crashed with: %s", err)
	}
}
