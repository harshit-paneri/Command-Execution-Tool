package main

import (
	"bytes"
	"embed"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

var frontend embed.FS

func main() {
	r := gin.Default()

	// Serve static files
	r.StaticFS("/", http.FS(frontend))

	// API endpoint for command execution
	r.POST("/execute", func(c *gin.Context) {
		var request struct {
			Command string `json:"command"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cmd := exec.Command("sh", "-c", request.Command)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  err.Error(),
				"stdout": stdout.String(),
				"stderr": stderr.String(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		})
	})

	r.Run(":31337")
}