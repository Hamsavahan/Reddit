package main

import (
	"protoactor-simulation/simulation"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	system := actor.NewActorSystem()
	sim := simulation.NewSimulation(system, 1000) // Simulate 1000 users

	router := gin.Default()

	router.POST("/start-simulation", func(c *gin.Context) {
		go sim.Run()
		c.JSON(200, gin.H{"message": "Simulation started"})
	})

	router.GET("/status", func(c *gin.Context) {
		status := sim.GetStatus()
		c.JSON(200, status)
	})

	router.GET("/subreddits", func(c *gin.Context) {
		subreddits := sim.GetSubreddits()
		c.JSON(200, gin.H{"subreddits": subreddits})
	})

	router.GET("/posts", func(c *gin.Context) {
		posts := sim.GetPosts()
		c.JSON(200, gin.H{"posts": posts})
	})

	router.GET("/comments/:postID", func(c *gin.Context) {
		postID := c.Param("postID")
		comments := sim.GetComments(postID)
		c.JSON(200, gin.H{"comments": comments})
	})

	router.GET("/feed/:username", func(c *gin.Context) {
		username := c.Param("username")
		feed := sim.GetFeed(username)
		c.JSON(200, gin.H{"feed": feed})
	})

	router.GET("/direct-messages/:username", func(c *gin.Context) {
		username := c.Param("username")
		messages := sim.GetDirectMessages(username)
		c.JSON(200, gin.H{"direct_messages": messages})
	})

	router.Run("localhost:8080")
}
