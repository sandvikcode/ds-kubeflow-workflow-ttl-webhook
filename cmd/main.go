package main

import (
	"fmt"
	"html"
	"io/ioutil"
	m "k8s-kubeflow-mutate-webhook/pkg/mutate"
	"log"
	"net/http"
	"time"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello %q", html.EscapeString(r.URL.Path))
}

func handleMutate(w http.ResponseWriter, r *http.Request) {

	log.Println("New request noticed")
	// read the body / request
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	// mutate the request
	log.Println("Ready to mutate the input")
	mutated, err := m.Mutate(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}
	log.Println("Mutate is done will return the output")
	// and write it back
	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
	log.Println("Sucess")
}

func main() {

	log.Println("The service is started")
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/mutate", handleMutate)

	log.Println("The server is started")
	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}

	log.Fatal(s.ListenAndServeTLS("/etc/tls-secret/tls.crt", "/etc/tls-secret/tls.key"))

}
