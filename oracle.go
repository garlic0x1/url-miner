package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

var (
	Filter sync.Map
)

func filter(u string, n int) bool {
	parsed, _ := url.Parse(u)
	buildurl := parsed.Scheme + "://" + parsed.Host + parsed.Path
	str := fmt.Sprintf("%d%s", n, buildurl)
	_, present := Filter.Load(str)
	if present {
		return false
	}
	Filter.Store(str, true)
	return true
}

// goroutine to handle output
func writer() {
	for res := range Results {
		reflected := 0
		str := res.URL

		if strings.Contains(res.URL, "?") {
			str += "&"
		} else {
			str += "?"
		}

		for i, param := range res.Parameters {
			hash := fmt.Sprintf("zzx%dy", i)
			if strings.Contains(res.Response, hash) {
				reflected++
				str += param + "=" + hash + "&"
			}
		}

		// remove trailing '&'
		str = str[:len(str)-1]

		if reflected > 0 && filter(str, reflected) {
			fmt.Println(str)
		}
	}
}
