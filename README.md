# url-miner
Finds hidden GET parameters by testing for reflection  
By default, 64 parameters are tested per request, this can be increased in some cases to 64,000 to significantly speed up your scan  

# Installation
`$ go install github.com/garlic0x1/url-miner@main`  

# Usage
Takes target urls from stdin, and a wordlist using the `-w` flag  
Example:
```
$ echo http://testphp.vulnweb.com/listproducts.php | url-miner -w wordlist.txt 
[reflected] http://testphp.vulnweb.com/listproducts.php?comment=zzxy5
[reflected] http://testphp.vulnweb.com/listproducts.php?cat=zzxy52
[reflected] http://testphp.vulnweb.com/listproducts.php?newpassword=zzxy5
[reflected] http://testphp.vulnweb.com/listproducts.php?artist=zzxy58
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
