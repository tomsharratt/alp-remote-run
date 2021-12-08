package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type ExecuteRequest struct {
	Files []File
}

type File struct {
	Name    string
	Path    string
	Content string
}

func execute(w http.ResponseWriter, req *http.Request) {
	var er ExecuteRequest

	err := decodeJSONBody(w, req, &er)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "ExecutionRequest: %+v", er)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/execute", execute)

	http.ListenAndServe(":8080", mux)
}
