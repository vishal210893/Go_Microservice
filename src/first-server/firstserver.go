package main

import (
	"log"
	"net/http"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

type api struct {
	addr string
}

func (s *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/":
			w.Write([]byte("Index page"))
			return
		case "/users":
			w.Write([]byte("Hello user"))
			return
		}
	case http.MethodPost:
		w.Write([]byte("Hello post"))
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Hello world"))
}

func main() {
	api := &api{
		addr: ":8080",
	}

	server := &http.Server{
		Addr:    api.addr,
		Handler: api,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
