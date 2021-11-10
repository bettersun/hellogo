package server

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8012",
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/hello", hello)

	log.Println("ListenAndServe 127.0.0.1:8012")
	server.ListenAndServe()
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome, bettersun")
}

func hello(w http.ResponseWriter, r *http.Request) {
	resp := fmt.Sprintf("[%v] Hello, world.", r.Method)
	fmt.Fprintf(w, resp)
}
