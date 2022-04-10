# url-miner
Finds hidden GET parameters by testing for reflection  

# Usage
Takes target urls from stdin, and a wordlist using the `-w` flag  
Example:
```
$ echo https://www.google.com/search | url-miner  -w wordlist
[reflected] q=zzxy0
[reflected] search=zzxy1
[reflected] video=zzxy2
[reflected] query=zzxy7
[reflected] hq=zzxy14
[reflected] action=zzxy24
```

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
