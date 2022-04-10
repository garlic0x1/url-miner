package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
)

// goroutine to handle output
func writer(results chan string) {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for res := range results {
		fmt.Fprintln(w, res)
	}
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

	// set up concurrency limit
	sem := make(chan struct{}, *threads)

	// Set up async
	var wg sync.WaitGroup

	// open chans
	results := make(chan string)

	// start pushing input
	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			u := s.Text()

			// start another goroutine if not too many
			select {
			case sem <- struct{}{}:
				wg.Add(1)
				go func() {
					mine(u, *wordlist, *nparams, results)
					<-sem
					wg.Done()
				}()
			default:
				mine(u, *wordlist, *nparams, results)
			}
		}

		// close reults chan when all miners done, ending the program
		wg.Wait()
		close(results)
	}()

	// call writer, which closes after all workers are done
	writer(results)
}
