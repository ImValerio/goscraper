package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/html"
)

type ScrapeDto struct {
	Urls []string `json:"urls"`
	Tags []string `json:"tags"`
}

var rdb *redis.Client
var ctx context.Context

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx = context.Background()

	router := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)
	router.Use(middleware.Logger)

	router.Post("/scrape", scrapeUrl)

	addr := "localhost:5000"

	slog.Info("Server started on: " + addr)
	http.ListenAndServe(addr, router)
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

	for _, url := range req.Urls {
		miners = append(miners, &Miner{url, []string{}})
	}

	for _, miner := range miners {
		handleMiner(miner, req, &wg)
	}

	wg.Wait()
	encoder := json.NewEncoder(w)

	encoder.Encode(miners)
}

func handleMiner(miner *Miner, req ScrapeDto, wg *sync.WaitGroup) {
	id := generateId(miner.Url, req.Tags)

	val, err := rdb.Get(ctx, id).Result()
	if err == redis.Nil {
		wg.Add(1)
		miner.scrapeUrl(wg, req.Tags)
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("cache hit")
		miner.Res = strings.Split(val, "|")
	}
}

func generateId(url string, tags []string) string {
	return url + strings.Join(tags, "")
}

func storeInCache(miner *Miner, tags []string) {
	key := generateId(miner.Url, tags)
	value := strings.Join(miner.Res, "|") // Store as a comma-separated string
	err := rdb.Set(ctx, key, value, 0).Err()
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
				//end of the file, break out of the loop
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

	// Trim leading and trailing whitespace
	content = strings.TrimSpace(content)

	return content
}
