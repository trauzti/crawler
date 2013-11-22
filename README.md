wuddlypums
==============
An awesome web crawler written in Go

Features
==============
* Opens URLs in parrallel
* Canonicalizes URLs
* Visits URLs only once
* Allows a parameter for the maximum number of pages to crawl
* Is very polite :) waits 300 milliseconds between HTTP GET requests
* Follows the rules from robots.txt

Authors
==============
Petur Orri Ragnarsson <peturor@gmail.com>

Trausti Saemundsson <trauzti@gmail.com>

Build
==============
```bash
$ go build
```

Run
==============
```bash
$ ./crawler <URL> <TOPIC> <QUERYWORDS> [N]
```

Where N is the maximum number of pages to crawl

Note that the querywords must be within quotation marks (see example.sh)
