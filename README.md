# url-miner
Finds hidden GET params  

# Usage
Takes target urls from stdin, and a wordlist using the `-w` flag  
Example: `$ cat endpoints.txt | url-miner -w wordlist.txt`  

# Help
```
$ ./url-miner -h
Usage of ./url-miner:
  -s int
    	Number of params per request. (default 64)
  -t int
    	Number of threads to use. (default 8)
  -w string
    	Wordlist to mine.
```
