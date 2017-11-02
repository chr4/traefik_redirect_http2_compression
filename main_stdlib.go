package main

import (
	"compress/gzip"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/containous/mux"
	"github.com/urfave/negroni"
	"net/http"
)

// Compress is a middleware that allows redirection
type Compress struct{}

// ServerHTTP is a function used by Negroni
func (c *Compress) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	gzipHandler(next).ServeHTTP(rw, r)
}

func gzipHandler(h http.Handler) http.Handler {
	wrapper, err := gziphandler.GzipHandlerWithOpts(
		gziphandler.CompressionLevel(gzip.DefaultCompression),
		gziphandler.MinSize(1))
	if err != nil {
		fmt.Println(err)
	}
	return wrapper(h)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	url := "https://localhost:3000/end"
	w.Header().Set("Location", url)
	w.WriteHeader(302)

	fmt.Fprintf(w, "<html><body>You are being <a href=\""+url+"\">redirected</a>.</body></html>")
}

func end(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<html><body>
	<p><a href="/redirect">/redirect</a></p>
	</body></html>
	`)
	fmt.Fprintf(w, "Done!")
}

func main() {
	systemRouter := mux.NewRouter()
	negroniInstance := negroni.New()
	negroniInstance.Use(&Compress{})
	negroniInstance.UseHandler(systemRouter)

	systemRouter.Methods("GET", "HEAD").Path("/redirect").HandlerFunc(redirect)
	systemRouter.Methods("GET", "HEAD").Path("/end").HandlerFunc(end)

	http.ListenAndServeTLS(":3000", "selfsigned.crt", "selfsigned.key", negroniInstance)
}
