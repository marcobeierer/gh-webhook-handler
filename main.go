package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
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
	secret = args[2]
	script = args[3]

	http.HandleFunc("/", runScript)
	http.ListenAndServe(address, nil)
}

func runScript(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	signature := req.Header.Get("X-Hub-Signature")

	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(body)
	expectedSignature := mac.Sum(nil)

	if hmac.Equal([]byte(signature), expectedSignature) {
		log.Println("wrong secret provided", signature, expectedSignature)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	cmd := exec.Command(script)

	err = cmd.Run()
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
