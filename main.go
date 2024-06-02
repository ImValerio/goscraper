package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/html"
)

type ScrapeDto struct {
	Urls []string   `json:"urls"`
	Tags [][]string `json:"tags"`
}

type ServerConfig struct {
	Port        string
	RedisClient *redis.Client
	Ctx         context.Context
}

var serverConfig ServerConfig

func loadServerConfig() {

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "5000"
	}
	port = ":" + port

	serverConfig = ServerConfig{
		Port: port,
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     redisHost + ":6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		Ctx: context.Background(),
	}
}
func main() {
	loadServerConfig()

	// Check if connected to the Redis server successfully
	_, err := serverConfig.RedisClient.Ping(serverConfig.Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	slog.Info("Connected to Redis server")

	router := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)
	router.Use(middleware.Logger)

	router.Get("/", home)
	router.Post("/scrape", scrapeUrl)

	slog.Info("Server started on port: " + serverConfig.Port)
	http.ListenAndServe(serverConfig.Port, router)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{"message": "welcome"}`))
}

func scrapeUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var req ScrapeDto
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup

	var miners []*Miner

	for i, url := range req.Urls {
		var tags []string
		if i < len(req.Urls)-1 {
			tags = req.Tags[i]
		} else {
			tags = req.Tags[len(req.Tags)-1]
		}

		miners = append(miners, &Miner{url, tags, []string{}})
	}

	for _, miner := range miners {
		handleMiner(miner, req, &wg)
	}

	wg.Wait()
	encoder := json.NewEncoder(w)

	encoder.Encode(miners)
}

func handleMiner(miner *Miner, req ScrapeDto, wg *sync.WaitGroup) {
	id := miner.generateId()

	val, err := serverConfig.RedisClient.Get(serverConfig.Ctx, id).Result()
	if err == redis.Nil {
		wg.Add(1)
		miner.scrapeUrl(wg)
	} else if err != nil {
		panic(err)
	} else {
		slog.Info("cache hit")
		miner.Res = strings.Split(val, "|")
	}
}

func (m *Miner) generateId() string {
	return m.Url + strings.Join(m.Tags, "")
}

func (m *Miner) storeInCache() {
	key := m.generateId()
	value := strings.Join(m.Res, "|") // Store as a comma-separated string
	err := serverConfig.RedisClient.Set(serverConfig.Ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func loadTokens(tkz *html.Tokenizer) []*html.Token {
	var tokens []*html.Token
	for {
		tokenType := tkz.Next()
		token := tkz.Token()

		if tokenType == html.ErrorToken {
			err := tkz.Err()
			if err == io.EOF {
				break
			}
			log.Fatalf("error tokenizing HTML: %v", tkz.Err())
		}
		tokens = append(tokens, &token)
	}

	return tokens
}

func clearData(token *html.Token) string {
	content := token.Data
	content = strings.ReplaceAll(content, "\\n", "\n")
	content = strings.ReplaceAll(content, "\\t", "\t")
	content = strings.TrimSpace(content) // Trim leading and trailing whitespace

	return content
}
