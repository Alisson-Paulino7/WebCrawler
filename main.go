package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/html"
)

var (
	links 	[]string
	visited map[string]bool = map[string]bool{}
	mu 		sync.Mutex
)

func main() {
	visitLink("https://unifapce.edu.br")

	fmt.Printf("Quantidade de links: %d", len(links))

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
				links = append(links, link.String())

				visitLink(link.String())
			}
		}
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c)
	}
}