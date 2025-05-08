package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var RedisClient *redis.Client

// Handler for the root path
func homeHandler(w http.ResponseWriter, req bunrouter.Request) error {
	fmt.Fprintf(w, "Welcome to Secure Proxy Management API Server!")
	return nil
}

func init() {
	redisUrl, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		log.Fatal("REDIS_URL not set. Please set it in your environment variables.")
		return
	}
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Printf("Failed to parse Redis URL: %v\n", err)
		return
	}

	// Initialize Redis client
	RedisClient = redis.NewClient(opt)

	// Ping the Redis server to ensure it's reachable
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Printf("Connected to Redis: %s", pong)
}

func addKeyHandler(w http.ResponseWriter, req bunrouter.Request) error {
	// params := req.Params()
	// for key, val := range params.Map() {
	// 	log.Printf("Key: %s, Value: %s\n", key, val)
	// }
	user := req.Param("user")
	expireStr := req.Param("expire")
	if user == "" || expireStr == "" {
		return bunrouter.JSON(w, map[string]string{
			"error": "Missing 'user' or 'expire' parameter.",
		})
	}

	key := uuid.New().String() // Generate a UUID as the API key
	// key, err := GenerateAPIKey(32) // 生成32字符的API Key

	expireSeconds := -1 // Default to -1 (infinite expiration)

	s, err := parseExpire(expireStr)
	if err != nil {
		return bunrouter.JSON(w, map[string]string{
			"error": "Invalid expiration format. Use 'infinite', 'forever', or a timestamp.",
		})
	}
	if s != "-1" {
		expireDate, err := time.Parse("20060102150405", s)
		if err != nil {
			return bunrouter.JSON(w, map[string]string{
				"error": "Invalid expiration format. Use 'infinite', 'forever', or a timestamp.",
			})
		}
		if int64(expireDate.Sub(time.Now()).Seconds()) < 0 {
			return bunrouter.JSON(w, map[string]string{
				"error": "Expiration date is in the past.",
			})
		}
		expireSeconds = int(expireDate.Sub(time.Now()).Seconds())
	}

	if expireSeconds == -1 {
		err = RedisClient.Set(ctx, key, user, 0).Err()
	} else {
		err = RedisClient.Set(ctx, key, user, time.Duration(expireSeconds)*time.Second).Err()
	}
	if err != nil {
		return bunrouter.JSON(w, map[string]string{
			"error": "Failed to save key to Redis.",
		})
	}

	// Respond with the generated key
	return bunrouter.JSON(w, map[string]string{
		"user":   user,
		"key":    key,
		"expire": s,
	})
}

func main() {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
		// bunrouter.Use(bunrouterotel.NewMiddleware(
		// 	bunrouterotel.WithClientIP(),
		// )),
	)
	// handler := otelhttp.NewHandler(router, "")

	router.GET("/", homeHandler)

	router.POST("/api/add/:user/:expire", addKeyHandler)

	port, ok := os.LookupEnv("SECURE_PROXY_MGMT_PORT")
	if !ok {
		log.Fatal("SECURE_PROXY_PORT environment variable is required")
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting Secure Proxy Management Server on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}

}
