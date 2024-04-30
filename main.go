package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

	http.ListenAndServe(":5000", router)
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
		miner.scrapeUrl(&wg)
	}

	wg.Wait()
	encoder := json.NewEncoder(w)

	encoder.Encode(miners)
}

func (m *Miner) scrapeUrl(wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := http.Get(m.Url)

	if err != nil {
		log.Fatal(err.Error())
	}

	tkz := html.NewTokenizer(res.Body)
	for {
		//get the next token type
		tokenType := tkz.Next()

		if tokenType == html.StartTagToken {
			token := tkz.Token()

			if token.Data == "title" {
				for tokenType != html.TextToken {
					tokenType = tkz.Next()

				}
				if tokenType == html.TextToken {
					content := tkz.Token().Data
					// fmt.Println("CONTENT:", content)

					m.Res = append(m.Res, content)
				}

			}
		}
		if tokenType == html.ErrorToken {
			err := tkz.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			log.Fatalf("error tokenizing HTML: %v", tkz.Err())
		}

	}

}
