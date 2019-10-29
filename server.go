package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
	// find port selected by service, default to 8080
	var port string
	var found bool
	port, found = os.LookupEnv("WEBSITES_PORT")
	if found == false {
		port, found = os.LookupEnv("PORT")
	}
	if found == true {
		port = fmt.Sprintf(":%s", port)
	} else {
		port = ":8080"
	}

	md, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	htmlFlags := mdhtml.CommonFlags | mdhtml.HrefTargetBlank
	opts := mdhtml.RendererOptions{Flags: htmlFlags}
	renderer := mdhtml.NewRenderer(opts)

	htmlContent := markdown.ToHTML(md, parser, renderer)

	http.Handle("/", WithLogging(ReturnHtml(htmlContent)))

	fmt.Printf("Listening on port http://localhost%s\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log.Println(r.Method, r.URL.Path, time.Since(time.Now()), r.RemoteAddr, r.UserAgent())
		}()
		next.ServeHTTP(w, r)
	})
}

func ReturnHtml(htmlContent []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(htmlContent)
	})
}
