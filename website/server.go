package website

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/Alisson-Paulino7/WebCrawler/infra"
)

type DataLinks struct {
	Links []infra.VisitedLink
}

func Run() {
	tmpl, err := template.ParseFiles("website/templates/index.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	links, err := infra.FindAllLinks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(links)
	data := DataLinks{Links: links}

	tmpl.Execute(w, &data)
	})
	
	if err := http.ListenAndServe(":8080", nil); err == nil {
		fmt.Println("Servidor rodando na porta 8080")
	}
	
}