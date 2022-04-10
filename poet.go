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
	str := u
	if strings.Contains(u, "?") {
		str += "&"
	} else {
		str += "?"
	}
	for i, s := range params {
		hash := fmt.Sprintf("zzxy%d", i)
		str += s + "=" + hash + "&"
	}
	return str
}

func mine(params []string, u string) {
	text := request(buildPayload(params, u), Timeout)

	for i, param := range params {
		hash := fmt.Sprintf("zzxy%d", i)
		if strings.Contains(text, hash) {
			if strings.Contains(u, "?") {
				Results <- fmt.Sprintf("[reflected] %s&%s=%s", u, param, hash)
			} else {
				Results <- fmt.Sprintf("[reflected] %s?%s=%s", u, param, hash)
			}
		}
	}
}

// send keys that effect the response to results
func poet(u string, wordlist string, nparams int) {
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
			mine(params, u)

			// reset
			params = []string{}
			c = 0
		}
	}
	if c != 0 {
		mine(params, u)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
