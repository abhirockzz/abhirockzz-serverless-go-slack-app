package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/abhirockzz/funcy-go/pkg/function"
)

func main() {
	port, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		port = "8080"
	}
	http.HandleFunc("/api/funcy", function.Funcy)
	fmt.Println("Go server listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
