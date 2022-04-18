package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type Output struct {
	URL  string
	Keys []string
}

var (
	Filter sync.Map
)

func filter(u string, n int) bool {
	if n < 4 {
		return true
	}
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
		keys := []string{}

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
				keys = append(keys, param)
			}
		}

		// remove trailing '&'
		str = str[:len(str)-1]

		if reflected > 0 && filter(str, reflected) {
			if Json {
				b, err := json.Marshal(Output{
					URL:  res.URL,
					Keys: keys,
				})
				if err != nil {
					log.Println("Error:", err)
					continue
				}
				fmt.Println(string(b))
			} else if Yaml {
				b, err := yaml.Marshal(Output{
					URL:  res.URL,
					Keys: keys,
				})
				if err != nil {
					log.Println("Error:", err)
					continue
				}
				fmt.Println(string(b))
			} else {
				fmt.Println(str)
			}
		}
	}
}
