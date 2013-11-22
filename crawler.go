package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "net/http"
    "strings"
    "strconv"
    "code.google.com/p/go.net/html"
    "container/heap"
)

// Toggle extra output.
var print_info = false

var frontier = PriorityQueue{}
var url, topic = "", ""
var maxCrawl int
var querywords []string
var visitedUrls = make(map[string]bool)
var totalUrlsFound = 0

func makeAbsoluteUrl(base, rest string) string {
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
func isAcceptableProtocol(url string) bool {
    return !(strings.Contains(url, ":") && url[:5] != "http:")
}

func extractBasePath(url string) string {
    count := 0
    findThirdSlash := func(c rune) bool {
        if c == '/' {
            count += 1
            return count == 3
        }
        return false
    }

    i := strings.IndexFunc(url, findThirdSlash)
    if i > 0 {
        return url[:i]
    }

    return url
}

func canonicalizeUrl(url string) string {
    if i := strings.IndexAny(url, "?#"); i >= 0 {
        return url[:i]
    }

    return url
}

func parseRobots(url string) []string {
    return []string{}
}

func getBody(url string) string {
    resp, err := http.Get(url)
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

func Crawl(starturl string)(pagesCrawled int, foundCount int) {
    heap.Init(&frontier)
    addToFrontier(starturl, 1)

    foundCount = 0
    pagesCrawled = 0
    for ; frontier.Len() > 0 && pagesCrawled < maxCrawl; pagesCrawled++ {
        currentItem := heap.Pop(&frontier).(*Item)
        currentUrl := currentItem.value

        if print_info {
            fmt.Println("Visiting ", currentUrl)
            fmt.Println("priority ", currentItem.priority)
            fmt.Println("secondaryPriority ", currentItem.secondaryPriority)
        }

        body := getBody(currentUrl)
        parseRobots(currentUrl)
        extractLinks(currentUrl, body)
        if findQuery(body) {
            fmt.Println("Query found in page:", currentUrl)
            foundCount += 1
        }
    }

    return
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
func addToFrontier(url string, priority int) {
    if _, in := visitedUrls[url]; !in {
        orderPriority -= 1
        heap.Push(&frontier, &Item{value:url, priority:priority, secondaryPriority:orderPriority})
        visitedUrls[url] = true
        totalUrlsFound += 1
    }
}

func extractLinks(url, body string) {
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
                nexturl = makeAbsoluteUrl(url, nexturl)
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


func main(){
    if len(os.Args) < 4 || len(os.Args) > 5 {
        printusage()
        return
    }

    url = os.Args[1]
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
    fmt.Println("Starting crawl, seed:", url)
    fmt.Println("Topic:", topic)
    fmt.Println("Query string:", strings.Join(querywords, " "))
    fmt.Println("Maximum number of pages to visit:", maxCrawl)
    fmt.Println("--------------------------------------------------------")

    // Do the work!
    pagesCrawled, foundCount := Crawl(url)

    fmt.Println("--------------------------------------------------------")
    fmt.Printf("Search complete. %d pages crawled\n", pagesCrawled)
    fmt.Printf("Search query \"%s\" found in %d pages\n", strings.Join(querywords, " "), foundCount)
    fmt.Printf("Total distinctive urls found: %d\n", totalUrlsFound)
    fmt.Println("--------------------------------------------------------")
}

