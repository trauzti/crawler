package gocrawler

import (
    "code.google.com/p/go.net/html"
    "container/heap"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
    "sync"
    "sync/atomic"
    "time"
)

// Toggle extra output.
var print_info = false

var frontier struct {
    sync.Mutex
    pq PriorityQueue
}
var URL, topic = "", ""
var maxCrawl int
var querywords []string
var visitedUrls = make(map[string]bool)
var totalUrlsFound = 0
var wg sync.WaitGroup
var foundCount int64 = 0
var politeness = 300 * time.Millisecond


func makeAbsoluteUrl(base, rest string) string {
    // TODO: use url.isAbs
    if len(rest) >= 7 && rest[:4+3] == "http://" {
        return rest
    }
    base = extractBasePath(base)
    if rest[:1] != "/" {
        rest = "/" + rest
    }
    return base + rest
}

// Returns true iff the url is for an acceptable protocol or
// it is a relative url.
func isAcceptableProtocol(_url string) bool {
    return !(strings.Contains(_url, ":") && _url[:5] != "http:")
}

func extractBasePath(_url string) string {
    count := 0
    findThirdSlash := func(c rune) bool {
        if c == '/' {
            count += 1
            return count == 3
        }
        return false
    }

    i := strings.IndexFunc(_url, findThirdSlash)
    if i > 0 {
        return _url[:i]
    }

    return _url
}

// Does: remove port 80, querystrings (like mbl.is/?yeah) and www. from the beginning
// TODO: 
//     1) remove . and .. loops from the end
//     2) Add trailing /
func canonicalizeUrl(_url string) string {
    x, err := url.Parse(_url)
    if err != nil {
        fmt.Println("url.Parse", err)
    }
    host, path := x.Host, x.Path
    if len(host) > 3 && host[len(host)-3:] == ":80" {
        host = host[:len(host)-3]
    }
    var res string
    if len(host) > 4 && host[:4] == "www." {
        res = "http://" + host[4:] + path
    } else {
        res = "http://" + host + path
    }
    if res[len(res)-1:] != "/" {
        res = res + "/"
    }
    return res
}

func parseRobots(_url string) []string {
    return []string{}
}

func getBody(_url string) string {
    resp, err := http.Get(_url)
    if err != nil {
        fmt.Println("http.Get", err)
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("ioutil.ReadAll", err)
    }

    if err != nil {
        fmt.Println("html.Parse", err)
    }
    s := string(body)

    return s
}


func Crawl(starturl string)(pagesCrawled int) {
    addToFrontier(starturl, 1)

    pagesCrawled = 0
    for ; frontier.pq.Len() > 0 && pagesCrawled < maxCrawl; pagesCrawled++ {
        frontier.Lock()
        currentItem := heap.Pop(&frontier.pq).(*Item)
        frontier.Unlock()
        currentUrl := currentItem.value

        if print_info {
            fmt.Println("Visiting ", currentUrl)
            fmt.Println("priority ", currentItem.priority)
            fmt.Println("secondaryPriority ", currentItem.secondaryPriority)
        }
        wg.Add(1)
        go handleUrl(currentUrl) // go => runs in a different thread
        time.Sleep(politeness)
    }

    return
}
func handleUrl(currentUrl string) {
    body := getBody(currentUrl)
    parseRobots(currentUrl)
    extractLinks(currentUrl, body)
    if findQuery(body) {
        fmt.Println("Query found in page:", currentUrl)
        atomic.AddInt64(&foundCount, 1)
    }
    wg.Done()
}


func findQuery(body string) bool {
    // This is supposed to be a phrase query so the whole splitting
    // thing doesn't make much sense...
    wholeQuery := strings.Join(querywords, " ")
    lowerQuery := strings.ToLower(wholeQuery)
    lowerBody := strings.ToLower(body)

    return strings.Contains(lowerBody, lowerQuery)
}

// Add an url to the frontier if it hasn't been visited before.
// Give decreasing priority to new links which are not topical
// to enforce a breadth first search.
var orderPriority = 0
func addToFrontier(_url string, priority int) {
    frontier.Lock()
    if _, in := visitedUrls[_url]; !in {
        orderPriority -= 1
        heap.Push(&frontier.pq, &Item{value:_url, priority:priority, secondaryPriority:orderPriority})
        visitedUrls[_url] = true
        totalUrlsFound += 1
    }
    frontier.Unlock()
}

func extractLinks(_url, body string) {
    doc, err := html.Parse(strings.NewReader(body))
    if err != nil {
        fmt.Println("html.Parse", err)
    }

    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for i := 0; i < len(n.Attr) ; i++ {
                // Some pages put more attributes in their <a> tags.
                if n.Attr[i].Key != "href" {
                    continue
                }
                nexturl := n.Attr[i].Val
                // Ignore other protocol links, eg mailto, ftp or https
                // XXX: Do we want to ignore https? Should we fetch these pages anyway?
                if !isAcceptableProtocol(nexturl) {
                    continue
                }
                nexturl = makeAbsoluteUrl(_url, nexturl)
                nexturl = canonicalizeUrl(nexturl)

                var priority = 0 // we use priority 0 if the anchor text isn't found in the page
                if n.FirstChild != nil {
                    // Maybe there's a false assumption here on where the anchor text is...
                    var anchorText = strings.ToLower(n.FirstChild.Data)
                    if strings.Contains(anchorText, topic) {
                        if print_info {
                            fmt.Println("Found topic in link anchor text", n.FirstChild.Data)
                        }
                        priority = 1
                    }
                }
                addToFrontier(nexturl, priority)
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)
}

func printusage() {
    fmt.Fprintf(os.Stderr, "Usage: %s <URL> <TOPIC> <QUERYWORDS> [MAX_CRAWL]\n",  os.Args[0])
}


func Docrawl(){
    frontier.pq = PriorityQueue{}
    heap.Init(&frontier.pq)

    if len(os.Args) < 4 || len(os.Args) > 5 {
        printusage()
        return
    }

    URL = os.Args[1]
    topic = strings.ToLower(os.Args[2])
    querywords = strings.Fields(os.Args[3])

    // Maximum number of pages to crawl is an optional param with default value.
    if len(os.Args) == 5 {
        var err error
        if maxCrawl, err = strconv.Atoi(os.Args[4]); err != nil {

        //if err != nil {
            printusage()
            return
        }
    } else {
        maxCrawl = 500
    }

    fmt.Println("--------------------------------------------------------")
    fmt.Println("Starting crawl, seed:", URL)
    fmt.Println("Topic:", topic)
    fmt.Println("Query string:", strings.Join(querywords, " "))
    fmt.Println("Maximum number of pages to visit:", maxCrawl)
    fmt.Println("--------------------------------------------------------")

    // Do the work!
    pagesCrawled := Crawl(URL)

    fmt.Println("Waiting for all threads to terminate")

    wg.Wait()

    fmt.Println("--------------------------------------------------------")
    fmt.Printf("Search complete. %d pages crawled\n", pagesCrawled)
    fmt.Printf("Search query \"%s\" found in %d pages\n", strings.Join(querywords, " "), foundCount)
    fmt.Printf("Total distinctive urls found: %d\n", totalUrlsFound)
    fmt.Println("--------------------------------------------------------")
}
