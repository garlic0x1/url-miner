# url-miner
Finds hidden GET parameters by testing for reflection  
By default, 64 parameters are tested per request, this can be increased in some cases to 64,000 to significantly speed up your scan  

# Installation
Go install:  
```
go install github.com/garlic0x1/url-miner@main
```  
Docker install:  
```
git clone https://github.com/garlic0x1/url-miner
cd url-miner
sudo docker build -t "garlic0x1/url-miner" .
echo http://testphp.vulnweb.com/listproducts.php | sudo docker run --rm -i garlic0x1/url-miner
```

# Usage
Takes target urls from stdin, and a wordlist using the `-w` flag  
Examples:
```
$ echo http://testphp.vulnweb.com/listproducts.php | url-miner -w wordlist.txt 
[reflected] http://testphp.vulnweb.com/listproducts.php?cat=zzx54y
[reflected] http://testphp.vulnweb.com/listproducts.php?artist=zzx60y
```
```
$ echo http://testphp.vulnweb.com/listproducts.php | hakrawler -u | grep -e "vulnweb.com" | url-miner -w wordlist.txt
[reflected] http://testphp.vulnweb.com/listproducts.php?cat=zzx54y
[reflected] http://testphp.vulnweb.com/listproducts.php?artist=zzx60y
[reflected] http://testphp.vulnweb.com/hpp/?pp=zzx25y
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
