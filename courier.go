package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func request(u string, timeout int) string {
	proxyURL, _ := url.Parse(os.Getenv("PROXY"))

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println("Error creating request", u)
		return ""
	}

	// apply custom headers
	if Header.name != "" && Header.value != "" {
		req.Header.Set(Header.name, Header.value)
	}

	var client http.Client
	if UseProxy {
		client = http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: Insecure,
				},
			},
		}
	} else {
		client = http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: Insecure,
				},
			},
		}
	}

	// perform request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error performing request", u)
		return ""
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println("Error reading response:", err)
		return ""
	}
	//return bodyBytes
	return string(bodyBytes)
}
