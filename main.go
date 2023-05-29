package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/MateSousa/internal/postgres"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stdout, "%s - [%s] \"%s %s %s\" %s\n", r.RemoteAddr, time.Now().Format(time.RFC3339), r.Method, r.URL.Path, r.Proto, r.UserAgent())
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

