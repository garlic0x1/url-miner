package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/chromedp/chromedp"
)

var (
	sm         sync.Map
	Chrome     bool
	ChromeCtx  context.Context
	Timeout    int
	ScriptWait int
	Insecure   bool
	UseProxy   = false
	Header     header
	Queue      chan string
	Results    chan Result
)

type Result struct {
	URL        string
	Parameters []string
	Response   string
}

type header struct {
	name  string
	value string
}

// goroutine to handle input
func reader() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		Queue <- s.Text()
	}
	close(Queue)
}

// worker routines
func spawnWorkers(n int, wordlist *string, nparams *int, includeVals *bool) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			// spin up a chrome tab to use as well if needed
			var tab context.Context
			if Chrome {
				tctx, cancel := chromedp.NewContext(ChromeCtx)
				tab = tctx
				defer cancel()
			}
			// pop from queue
			for line := range Queue {
				u := ""
				parsed, err := url.Parse(line)
				if err != nil {
					u = line
				} else {
					if *includeVals {
						u = line
					} else {
						u = fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, parsed.Path)
					}
				}

				if isUnique(u) {
					poet(u, *wordlist, *nparams, tab)
				}
			}
			wg.Done()
		}()
	}

	// wait for all jobs to be finished before ending
	wg.Wait()
	close(Results)
}

func isUnique(url string) bool {
	_, present := sm.Load(url)
	if present {
		return false
	}
	sm.Store(url, true)
	return true
}

func parseHeader(h string) (string, string) {
	split1 := strings.Split(h, ":")
	name := strings.TrimSpace(split1[0])
	value := strings.TrimSpace(split1[1])
	return name, value
}

func main() {
	// set up flags
	threads := flag.Int("t", 8, "Number of threads to use.")
	nparams := flag.Int("s", 64, "Number of params per request.")
	wordlist := flag.String("w", "", "Wordlist to mine.")
	customheader := flag.String("head", "", "Custom header. Example: -head 'Hello: world'")
	insecure := flag.Bool("insecure", false, "Disable TLS verification.")
	chrome := flag.Bool("chrome", false, "Use headless browser to evaluate DOM.")
	includeVals := flag.Bool("d", false, "Include default GET values from input.")
	proxy := flag.String(("proxy"), "", "Proxy URL. Example: -proxy http://127.0.0.1:8080")
	timeout := flag.Int("timeout", 20, "Request timeout.")
	swait := flag.Int("wait", 0, "Seconds to wait on page after loading in chrome mode. (Use to wait for AJAX reqs)")
	flag.Parse()
	ScriptWait = *swait
	Insecure = *insecure
	Timeout = *timeout
	Chrome = *chrome

	// set up chrome ctx
	if *chrome {
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ProxyServer(*proxy),
			// block all images
			chromedp.Flag("blink-settings", "imagesEnabled=false"),
			chromedp.Flag("headless", true))...)

		ChromeCtx, cancel = chromedp.NewContext(ctx)
		defer cancel()
	}
	// set custom header
	if *customheader != "" {
		hname, hvalue := parseHeader(*customheader)
		Header = header{hname, hvalue}
	}

	// set proxy
	if *proxy != "" {
		os.Setenv("PROXY", *proxy)
		UseProxy = true
	}

	// check for wordlist
	if *wordlist == "" {
		fmt.Fprintln(os.Stderr, "No wordlist detected, use `cat urls.txt | url-miner -w wordlist.txt`")
		os.Exit(1)
	}

	// check for stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "No input detected, use `cat urls.txt | url-miner -w wordlist.txt`")
		os.Exit(1)
	}

	// open chans
	Results = make(chan Result)
	Queue = make(chan string, 1)

	// start pushing input, when done, close Queue
	go reader()
	// start *threads workers
	// ended by Queue closing, when done, close Results
	go spawnWorkers(*threads, wordlist, nparams, includeVals)
	// start writing output
	// ended by Results closing
	writer()
}
