package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sync"
)

var sm sync.Map

// goroutine to handle output
func writer(results chan string) {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for res := range results {
		fmt.Fprintln(w, res)
	}
}

func isUnique(url string) bool {
	_, present := sm.Load(url)
	if present {
		return false
	}
	sm.Store(url, true)
	return true
}

func main() {
	// set up flags
	threads := flag.Int("t", 8, "Number of threads to use.")
	nparams := flag.Int("s", 64, "Number of params per request.")
	wordlist := flag.String("w", "", "Wordlist to mine.")
	flag.Parse()

	// check for wordlist
	if *wordlist == "" {
		fmt.Fprintln(os.Stderr, "No wordlist detected, use `echo $url | mine-params -w $wordlist`")
		os.Exit(1)
	}

	// check for stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "No input detected, use `echo $url | mine-params -w $wordlist`")
		os.Exit(1)
	}

	// set up concurrency
	sem := make(chan struct{}, *threads)
	var wg sync.WaitGroup

	// open chans
	results := make(chan string)

	// start pushing input
	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			u := ""
			line := s.Text()
			parsed, err := url.Parse(line)
			if err != nil {
				u = line
			} else {
				u = fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, parsed.Path)
			}

			if isUnique(u) {
				// start another goroutine if not too many
				select {
				case sem <- struct{}{}:
					wg.Add(1)
					go func() {
						poet(u, *wordlist, *nparams, results)
						<-sem
						wg.Done()
					}()
				default:
					poet(u, *wordlist, *nparams, results)
				}
			}
		}

		// close reults chan when all miners done, ending the program
		wg.Wait()
		close(results)
	}()

	// call writer, which closes after all workers are done
	writer(results)
}
