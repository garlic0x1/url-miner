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
				chromeRequest(u, params, Timeout, tab)
			} else {
				request(u, params, Timeout)
			}

			// reset
			params = []string{}
			c = 0
		}
	}
	if c != 0 {
		if Chrome {
			chromeRequest(u, params, Timeout, tab)
		} else {
			request(u, params, Timeout)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
