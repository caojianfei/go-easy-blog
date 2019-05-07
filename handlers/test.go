package handlers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	_, err := fmt.Fprintf(w, "test router")
	if err != nil {
		log.Fatal(err)
	}
}