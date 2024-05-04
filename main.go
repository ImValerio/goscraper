package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"golang.org/x/net/html"
)

type Miner struct {
	Url string   `json:"url"`
	Res []string `json:"res"`
}

type ScrapeDto struct {
	Urls []string `json:"urls"`
	Tags []string `json:"tags"`
}

func main() {
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

	http.ListenAndServe("localhost:5000", router)
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
		wg.Add(1)
		miner.scrapeUrl(&wg, req.Tags)
	}

	wg.Wait()
	encoder := json.NewEncoder(w)

	encoder.Encode(miners)
}

func (m *Miner) printTokens(wg *sync.WaitGroup, tags []string) {
	defer wg.Done()
	res, err := http.Get(m.Url)

	if err != nil {
		log.Fatal(err.Error())
	}

	tkz := html.NewTokenizer(res.Body)

	counter := 0

	for counter < 100 {

		// tokenBefore := tkz.Token()
		tokenType := tkz.Next()
		token := tkz.Token()
		if tokenType == html.StartTagToken {
			fmt.Println("-> ", token.Data)
		}
		if tokenType == html.ErrorToken {
			err := tkz.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			log.Fatalf("error tokenizing HTML: %v", tkz.Err())
		}
		counter++
	}

}

func (m *Miner) scrapeUrl(wg *sync.WaitGroup, tags []string) {
	defer wg.Done()
	res, err := http.Get(m.Url)

	if err != nil {
		log.Fatal(err.Error())
	}

	tkz := html.NewTokenizer(res.Body)

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
	fmt.Println("found ", len(tokens), " tokens")

	index := 0
	prevIndex := 0
	counter := 0

	for index < len(tokens)-1 {

		token := tokens[index]

		if token.Type != html.StartTagToken {
			index++
			continue
		}

		if token.Data == tags[counter] {

			if counter == 0 {
				prevIndex = index + 1
			}
			// fmt.Println("----> ", token.Data, token.ype)
			if counter == len(tags)-1 {

				innerIndex := index
				for token.Type != html.TextToken {

					if innerIndex+1 > len(tokens)-1 {
						fmt.Println("BREAK")
						counter = 0
						index = prevIndex
						break
					}
					token = tokens[innerIndex]
					innerIndex++
				}

				if innerIndex+1 > len(tokens)-1 {
					continue
				}
				content := token.Data
				content = strings.ReplaceAll(content, "\\n", "\n")
				content = strings.ReplaceAll(content, "\\t", "\t")

				// Trim leading and trailing whitespace
				content = strings.TrimSpace(content)
				fmt.Println("CONTENT:", content)
				if content != "" {
					m.Res = append(m.Res, content)
				}
				counter = 0
				index = prevIndex
				continue

			} else {
				index++
				counter++
				continue
			}

		} else {
			if counter == 0 {
				index++
				continue
			}
			index = prevIndex
			counter = 0
		}

	}

}
