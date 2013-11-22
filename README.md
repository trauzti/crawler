=== wuddlypums ===
An awesome web crawler written in Go

=== Features ===
1) Opens URLs in parrallel
2) Canonicalizes URLs
3) Visits URLs only once
4) Allows a parameter for the maximum number of pages to crawl
5) Is very polite :) waits 300 milliseconds between HTTP GET requests
6) Follows the rules from robots.txt

=== Authors ===
Petur Orri Ragnarsson <peturor@gmail.com>
Trausti Saemundsson <trauzti@gmail.com>

=== Build ===
$ go build

=== Run ===
./crawler <URL> <TOPIC> <QUERYWORDS> [N]
Where N is the maximum number of pages to crawl
Note that the querywords must be within quotation marks (see example.sh)
