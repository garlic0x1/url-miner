package main

import (
	"fmt"
	"strings"
)

type Result struct {
	URL        string
	Parameters []string
	Response   string
}

// goroutine to handle output
func writer() {
	for res := range Results {
		for i, param := range res.Parameters {
			hash := fmt.Sprintf("zzx%dy", i)
			if strings.Contains(res.Response, hash) {
				if strings.Contains(res.URL, "?") {
					fmt.Printf("[reflected] %s&%s=%s\n", res.URL, param, hash)
				} else {
					fmt.Printf("[reflected] %s?%s=%s\n", res.URL, param, hash)
				}
			}
		}
	}
}
