package main

import (
	"bufio"
	"freechatgpt/internal/tokens"
	"os"
	"strings"

	chatgpt_types "freechatgpt/internal/chatgpt"
	"freechatgpt/internal/otp"

	"github.com/acheong08/endless"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var HOST string
var PORT string
var ACCESS_TOKENS tokens.AccessToken
var proxies []string

func checkProxy() {
	// first check for proxies.txt
	proxies = []string{}
	if _, err := os.Stat("proxies.txt"); err == nil {
		// Each line is a proxy, put in proxies array
		file, _ := os.Open("proxies.txt")
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Split line by :
			proxy := scanner.Text()
			proxy_parts := strings.Split(proxy, ":")
			if len(proxy_parts) > 1 {
				proxies = append(proxies, proxy)
			} else {
				continue
			}
		}
	}
	// if no proxies, then check env http_proxy
	if len(proxies) == 0 {
		proxy := os.Getenv("http_proxy")
		if proxy != "" {
			proxies = append(proxies, proxy)
		}
	}
}

func init() {
	_ = godotenv.Load(".env")

	HOST = os.Getenv("SERVER_HOST")
	PORT = os.Getenv("SERVER_PORT")
	if HOST == "" {
		HOST = "127.0.0.1"
	}
	if PORT == "" {
		PORT = "8080"
	}
	checkProxy()
	readAccounts()
	err := otp.InitOTP()
	if err != nil {
		println("Warning: Failed to initialize OTP service:", err.Error())
		println("OTP functionality will not be available")
	}
	scheduleTokenPUID()
}
func main() {
	defer chatgpt_types.SaveFileHash()
	defer otp.CloseOTP()
	router := gin.Default()

	router.Use(cors)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	admin_routes := router.Group("/admin")
	admin_routes.Use(adminCheck)

	/// Admin routes
	admin_routes.PATCH("/password", passwordHandler)
	admin_routes.PATCH("/tokens", tokensHandler)
	/// Public routes
	router.OPTIONS("/v1/chat/completions", optionsHandler)
	router.POST("/v1/chat/completions", Authorization, nightmare)
	router.OPTIONS("/v1/audio/speech", optionsHandler)
	router.POST("/v1/audio/speech", Authorization, tts)
	router.OPTIONS("/v1/audio/transcriptions", optionsHandler)
	router.POST("/v1/audio/transcriptions", Authorization, stt)
	router.OPTIONS("/v1/models", optionsHandler)
	router.GET("/v1/models", Authorization, simulateModel)
	endless.ListenAndServe(HOST+":"+PORT, router)
}
