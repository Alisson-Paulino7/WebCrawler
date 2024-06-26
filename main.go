package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Alisson-Paulino7/WebCrawler/infra"
	"github.com/Alisson-Paulino7/WebCrawler/website"

	"golang.org/x/net/html"
)

var mu sync.Mutex

var(
	link string
	action string
)

func init() {
	flag.StringVar(&link, "url", 
	"https://go.dev", "url para iniciar WebCrawler")
	flag.StringVar(&action, "action", "Website", "Ação a ser executada")
}

func main() {
	flag.Parse()

	switch action {
		case "Website":
			website.Run()
		case "webcrawler":
			done := make(chan bool)
			go visitLink(link)

			<-done

		default:
			fmt.Printf("Ação '%s' é inválida", action)
	}

}

func visitLink (link string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Erro recuperado: ", r)
		}
	}()

	fmt.Printf("Visitando: %s\n", link)

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
				if   err  != nil || link.Scheme == "" {
					continue
				}
				mu.Lock()
				if infra.CheckLink(link.String()) {
					fmt.Printf("Link já visitado: %v\n", link)
					mu.Unlock()
					continue
				}
				mu.Unlock()

				visitedLink := infra.VisitedLink{
					Website    : link.Hostname(),
					Link       : link.String(),
					VisitedDate: time.Now(),
				}

				infra.Insert("links", visitedLink)
				go visitLink(link.String())
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c)
	}
}