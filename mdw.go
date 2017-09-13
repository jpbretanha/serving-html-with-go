package main

import (
	"log"
	"net/http"
	"html/template"
	"path"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// middleware provides a convenient mechanism for filtering HTTP requests
// entering the application. It returns a new handler which performs various
// operations and finishes with calling the next HTTP handler.
type middleware func(http.HandlerFunc) http.HandlerFunc
type data struct {
        Data template.HTML
}
// chainMiddleware provides syntactic sugar to create a new middleware
// which will be the result of chaining the ones received as parameters.
func chainMiddleware(mw ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Logged connection from %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}

func withTracing(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Tracing request for %s", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Println("reached home")


	fp := path.Join("templates", "teste.html")

	tmpl, _ := template.ParseFiles(fp)

	val := data {Data:template.HTML("<a href='/teste'>clique aqui</a>")} 
    tmpl.Execute(w, val);
}

func teste(w http.ResponseWriter, r *http.Request) {
	log.Println("reached teste")
	t := template.New("Test")
    t, _ = t.Parse("<html><body>Página de teste {{.Data}}</body></html>")
	val := data {Data:template.HTML("<a href='/'>Voltar</a>")} 
    t.Execute(w, val)
}

func main() {
	mw := chainMiddleware(withLogging, withTracing)
	http.Handle("/", mw(home))
	http.Handle("/teste", mw(teste))
	log.Fatal(http.ListenAndServe(":8080", nil))
}