package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var secret string
var script string

func main() {
	args := os.Args

	if len(args) != 4 {
		fmt.Printf("usage: %s [address] [secret] [script]\n", args[0])
		return
	}

	address := args[1]
	secret = args[3]
	script = args[4]

	http.HandleFunc("/", runScript)
	http.ListenAndServe(address, nil)
}

func runScript(rw http.ResponseWriter, req *http.Request) {
	if secret != req.Header.Get("X-Hub-Signature") {
		log.Println("wrong secret provided")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	cmd := exec.Command(script)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
