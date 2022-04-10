package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	sm       sync.Map
	Timeout  int
	Insecure bool
	UseProxy = false
	Header   header
	Queue    chan string
	Results  chan string
)

type header struct {
	name  string
	value string
}

// goroutine to handle output
func writer() {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for res := range Results {
		//fmt.Fprintln(w, res)
		fmt.Println(res)
	}
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
func spawnWorkers(n int, wordlist *string, nparams *int) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			// pop from queue
			for line := range Queue {
				u := ""
				parsed, err := url.Parse(line)
				if err != nil {
					u = line
				} else {
					u = fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, parsed.Path)
				}

				if isUnique(u) {
					poet(u, *wordlist, *nparams)
				}
			}
			wg.Done()
		}()
	}
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
	proxy := flag.String(("proxy"), "", "Proxy URL. Example: -proxy http://127.0.0.1:8080")
	timeout := flag.Int("timeout", 20, "Request timeout.")
	flag.Parse()
	Insecure = *insecure
	Timeout = *timeout

	// set custom header
	if *customheader != "" {
		hname, hvalue := parseHeader(*customheader)
		Header = header{hname, hvalue}
	}

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
	Results = make(chan string)
	Queue = make(chan string, 1)

	// start pushing input
	go reader()
	go spawnWorkers(*threads, wordlist, nparams)
	writer()
}
