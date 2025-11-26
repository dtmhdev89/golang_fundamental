package sample

import (
	"fmt"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request")
	defer log.Print("End Hello world request")
	fmt.Fprintf(w, "Hello World====")
}

type HandleViaStruct struct{}

func (*HandleViaStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request")
	defer log.Print("End Hello world request")
	fmt.Fprintf(w, "Hello World via Struct")
}
