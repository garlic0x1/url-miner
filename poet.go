package main

import (
	"bufio"
	"context"
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
		hash := fmt.Sprintf("zzx%dy", i)
		str += s + "=" + hash + "&"
	}
	return str
}

func chromeMine(params []string, u string, tab context.Context) {
	text := chromeRequest(buildPayload(params, u), Timeout, tab)

	Results <- Result{
		URL:        u,
		Parameters: params,
		Response:   text,
	}
}

func mine(params []string, u string) {
	text := request(buildPayload(params, u), Timeout)

	Results <- Result{
		URL:        u,
		Parameters: params,
		Response:   text,
	}
}

// send keys that effect the response to results
func poet(u string, wordlist string, nparams int, tab context.Context) {
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
			if Chrome {
				chromeMine(params, u, tab)
			} else {
				mine(params, u)
			}

			// reset
			params = []string{}
			c = 0
		}
	}
	if c != 0 {
		if Chrome {
			chromeMine(params, u, tab)
		} else {
			mine(params, u)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
