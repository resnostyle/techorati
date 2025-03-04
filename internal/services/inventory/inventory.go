package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "pong"})
    })

    // New route for /hello
    r.GET("/chicken", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
    })

    // New route for /goodbye
    r.GET("/steak", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Goodbye!"})
    })

    // New route for /status
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "running"})
    })

    r.Run(":8081") // Start the server on port 8081
} 