package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

type Miner struct {
	Url string   `json:"url"`
	Res []string `json:"res"`
}

func (m *Miner) scrapeUrl(wg *sync.WaitGroup, tags []string) {
	defer wg.Done()
	defer storeInCache(m, tags)
	res, err := http.Get(m.Url)

	if err != nil {
		log.Fatal(err.Error())
	}

	tkz := html.NewTokenizer(res.Body)

	tokens := loadTokens(tkz)

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
				content := clearData(token)
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