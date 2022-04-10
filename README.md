# url-miner
Finds hidden GET parameters by testing for reflection  

# Usage
Takes target urls from stdin, and a wordlist using the `-w` flag  
Example: `$ cat endpoints.txt | url-miner -w wordlist.txt`  

# Help
```
$ url-miner -h
Usage of url-miner:
  -insecure
    	Disable TLS verification.
  -proxy string
    	Proxy URL. E.g.: -proxy http://127.0.0.1:8080
  -s int
    	Number of params per request. (default 64)
  -t int
    	Number of threads to use. (default 8)
  -timeout int
    	Request timeout. (default 20)
  -w string
    	Wordlist to mine.

```
