package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Alisson-Paulino7/WebCrawler/infra"

	"golang.org/x/net/html"
)

var (
	visited map[string]bool = map[string]bool{}
	mu 		sync.Mutex
)

type VisitedLink struct {
	Website 	string    `bson:"website"`
	Link 		string    `bson:"link"`
	VisitedDate time.Time `bson:"visited_Date"`
}

func main() {
	visitLink("https://aprendagolang.com.br")

}

func visitLink (link string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Erro recuperado: ", r)
		}
	}()

	mu.Lock()
	if visited[link] {
		mu.Unlock()
		return
	}
	visited[link] = true
	mu.Unlock()

	fmt.Println("Visitando", link)

	resp, err := http.Get(link)
	if   err  != nil {
		panic(fmt.Sprintf("Erro ao fazer requisição HTTP na URL: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Status diferente de 200: %d", resp.StatusCode))
	}

	doc, err := html.Parse(resp.Body)
	if  err  != nil {
		panic(fmt.Sprintf("Erro ao fazer parse do HTML: %v", err))
	}

	extractLinks(doc)
}

func extractLinks(doc *html.Node) {
	if doc.Type == html.ElementNode && doc.Data == "a" {
		for _, a := range doc.Attr {
			if a.Key == "href" {
				link, err := url.Parse(a.Val)
				if   err  != nil || link.Scheme == ""{
					continue
				}
				
				visitedLink := VisitedLink{
					Website    : link.Hostname(),
					Link       : link.String(),
					VisitedDate: time.Now(),
				}

				infra.Insert("links", visitedLink)

				visitLink(link.String())
			}
		}
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c)
	}
}