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

var priorityqueue = PriorityQueue{}
var url, topic = "", ""
var N = 500
var querywords []string

func makeAbsoluteUrl(base, rest string) string {
    if len(rest) >= 7 && rest[:4+3] == "http://" {
        return rest
    }
    return base + "/" + rest
}

func canonicalizeUrl(url string) string {
    // TODO: writeme
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

func Crawl(starturl string) {
    heap.Init(&priorityqueue)
    priorityqueue.Push(&Item{value:starturl, priority:1})

    for priorityqueue.Len() > 0 {
        current_url := priorityqueue.Pop().(*Item).value

        body := getBody(current_url)
        parseRobots(current_url)
        extractLinks(current_url, body)

        break
    }
}

func AppendString(slice []string, data ...string) []string {
    m := len(slice)
    n := m + len(data)
    if n > cap(slice) { // if necessary, reallocate
        // allocate double what's needed, for future growth.
        newSlice := make([]string, (n+1)*2)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0:n]
    copy(slice[m:n], data)
    return slice
}


func extractLinks(url, body string) {
    doc, err := html.Parse(strings.NewReader(body))
    if err != nil {
        fmt.Println("html.Parse", err)
    }
    // TODO: check if querywords is actually found in this page =)
    fmt.Println("Query found in page:", url)
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for i := 0; i < len(n.Attr) ; i++ {
                nexturl := n.Attr[i].Val
                nexturl = makeAbsoluteUrl(url, nexturl)
                nexturl = canonicalizeUrl(nexturl)
                fmt.Println(nexturl)

                var priority = 0
                // Maybe there's a false assumption here on where the anchor text is...
                var anchorText = strings.ToLower(n.FirstChild.Data)
                if strings.Contains(anchorText, topic) {
                    fmt.Println("Found topic in link anchor text", n.FirstChild.Data)
                    priority = 1
                }
                priorityqueue.Push(&Item{value:nexturl, priority:priority})
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)
}

func printusage() {
    fmt.Fprintf(os.Stderr, "Usage: %s <URL> <TOPIC> <QUERYWORDS> <N>\n",  os.Args[0])
}


func main(){
    switch len(os.Args) {
        case 2:
            url = os.Args[1]
        case 5:
            url = os.Args[1]
            topic = strings.ToLower(os.Args[2])
            querywords = strings.Fields(os.Args[3])

            var err error
            N, err = strconv.Atoi(os.Args[4])
            if err != nil {
                printusage()
                return
            }

            fmt.Println(url, topic, querywords, N)
        default:
            printusage()
            return
    }
    fmt.Println("--------------------------------------------------------")
    fmt.Println("Starting crawl, seed:", url)
    fmt.Println("Topic:", topic)
    fmt.Println("Query string:", querywords)
    fmt.Println("Maximum number of pages to visit:", N)
    fmt.Println("--------------------------------------------------------")
    // Get the page
    Crawl(url)
    fmt.Println(priorityqueue)
}

