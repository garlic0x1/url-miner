package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// constructs payloads and creates an identifier for the injection point
func buildPayload(params []string, u string) string {
	str := u + "?"
	for i, s := range params {
		hash := fmt.Sprintf("zzxy%d", i)
		str += s + "=" + hash + "&"
	}
	return str
}

func mine(params []string, u string, results chan string) {
	text := request(buildPayload(params, u))

	for i, param := range params {
		hash := fmt.Sprintf("zzxy%d", i)
		if strings.Contains(text, hash) {
			results <- fmt.Sprintf("[reflected] %s=%s", param, hash)
		}
	}
}

// send keys that effect the response to results
func poet(u string, wordlist string, nparams int, results chan string) {
	//baseline := request(u)
	var params []string
	c := 0

	file, err := os.Open(wordlist)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if c < nparams {
			params = append(params, scanner.Text())
			c++
		} else {
			mine(params, u, results)

			// reset
			params = []string{}
			c = 0
		}
	}
	if c != 0 {
		mine(params, u, results)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
