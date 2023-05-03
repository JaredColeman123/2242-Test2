package main

import (
	"log"
	"net/http"
	"os"

	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/justinas/alice"
)

func home(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing home handler ...")
	w.Write([]byte("success"))
}

func request(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing request handler ...")
	w.Write([]byte("request is in application format"))
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing login handler ...")
	w.Write([]byte("login is successful"))
}

func logger(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing logger handler ...")
	w.Write([]byte("log to the file successful"))

}

func simple(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing simple handler ...")
	w.Write([]byte("Alice used successfully"))
}

//Middleware Function
//middleware 1

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware One ...")
		next.ServeHTTP(w, r)
		log.Println("Executing middleware One again ...")
	})
}

//Middleware Function 2
func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware Two ...")
		if r.URL.Path != "/" {
			return
		}

		next.ServeHTTP(w, r)
		log.Println("Executing middleware Two again ...")
	})
}

//enforce JSON content-type middleware

func enforceJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("Executing enforce JSON middleware ...")
		contentType := r.Header.Get("Content-Type")

		if contentType == "" {
			log.Println("No Content-Type header found in request")
			return
		} else if contentType != "application/json" {
			log.Println("Invalid Content-Type. Expected 'application/json' ")
			return
		}
		log.Println("Request is of type 'application/json' ")
		next.ServeHTTP(w, r)
	})
}

//Main Function

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", middlewareOne(middlewareTwo(http.HandlerFunc(home))))

	mux.Handle("/request", enforceJSON(http.HandlerFunc(request)))

	authMiddleware := httpauth.SimpleBasicAuth("username", "secret")
	mux.Handle("/login", authMiddleware(http.HandlerFunc(login)))

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/log", handlers.LoggingHandler(logFile, http.HandlerFunc(logger)))

	mux.Handle("/alice", alice.New(authMiddleware, middlewareOne).Then(http.HandlerFunc(simple)))

	log.Println("starting server on port 8000...")
	err = http.ListenAndServe(":8000", mux)

	if err != nil {
		log.Fatal(err)
	}

}
