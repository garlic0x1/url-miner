package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func chromeRequest(u string, timeout int, ctx context.Context) string {
	c1 := make(chan string, 1)

	go func() {
		var document string
		err := chromedp.Run(ctx,
			chromedp.Navigate(u),
			chromedp.Evaluate("document.documentElement.innerHTML", &document),
		)
		if err != nil {
			log.Println(err, u)
			return
		}

		c1 <- document
	}()

	// listen to timer and response, whichever happens first
	select {
	case document := <-c1:
		return document
	case <-time.After(time.Duration(timeout) * time.Second):
		return ""
	}
}

func request(u string, timeout int) string {
	proxyURL, _ := url.Parse(os.Getenv("PROXY"))

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println("Error creating request", err)
		return ""
	}
	req.Close = true

	// apply custom headers
	if Header.name != "" && Header.value != "" {
		req.Header.Set(Header.name, Header.value)
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: Insecure}
	if UseProxy {
		t.Proxy = http.ProxyURL(proxyURL)
	}

	client := http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: t,
	}

	// perform request
	resp, err := client.Do(req)
	if err != nil {
		//log.Println("Error performing request", err)
		return ""
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return ""
	}
	return string(bodyBytes)
}
