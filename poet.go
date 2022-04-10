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

// send keys that effect the response to results
func mine(u string, wordlist string, nparams int, results chan string) {
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
			text := request(buildPayload(params, u))

			for i, param := range params {
				hash := fmt.Sprintf("zzxy%d", i)
				if strings.Contains(text, hash) {
					fmt.Println("[reflected]", param+"="+hash)
				}
			}

			// reset
			params = []string{}
			c = 0
		}
	}
	if c != 0 {

		text := request(buildPayload(params, u))

		for i, param := range params {
			hash := fmt.Sprintf("zzxy%d", i)
			if strings.Contains(text, hash) {
				results <- fmt.Sprintf("[reflected] %s=%s", param, hash)
			}
		}

		// reset
		params = []string{}
		c = 0
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
