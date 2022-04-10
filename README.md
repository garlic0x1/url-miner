# url-miner
Finds hidden GET parameters by testing for reflection  

# Usage
Takes target urls from stdin, and a wordlist using the `-w` flag  
Example:
```
$ echo https://ac4f1f281ee763c2c025718500780061.web-security-academy.net/ | url-miner  -w wordlist.txt -s 2000
[reflected] number=zzxy3
[reflected] search=zzxy39
```

# Help
```
$ url-miner -h
Usage of url-miner:
  -head string
    	Custom header. Example: -head 'Hello: world'
  -insecure
    	Disable TLS verification.
  -proxy string
    	Proxy URL. Example: -proxy http://127.0.0.1:8080
  -s int
    	Number of params per request. (default 64)
  -t int
    	Number of threads to use. (default 8)
  -timeout int
    	Request timeout. (default 20)
  -w string
    	Wordlist to mine.

```
