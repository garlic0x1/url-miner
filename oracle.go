package main

import (
	"fmt"
	"strings"
)

// goroutine to handle output
func writer() {
	for res := range Results {
		reflected := false
		str := res.URL

		if strings.Contains(res.URL, "?") {
			str += "&"
		} else {
			str += "?"
		}

		for i, param := range res.Parameters {
			hash := fmt.Sprintf("zzx%dy", i)
			if strings.Contains(res.Response, hash) {
				reflected = true
				str += param + "=" + hash + "&"
			}
		}

		if reflected {
			// remove trailing '&'
			str = str[:len(str)-1]
			fmt.Println(str)
		}
	}
}
