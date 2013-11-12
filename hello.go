package main

import (
    "fmt"
    "strings"
    "io/ioutil"
    "net/http"
    "code.google.com/p/go.net/html"
)

func main(){
  resp, err := http.Get("http://www.google.com")
  if err != nil {
      // handle error
  }
  body, err := ioutil.ReadAll(resp.Body)

  s := ""
  for i:= 0; i < len(body); i++ {
    s += fmt.Sprintf("c", body[i])
  }
//  fmt.Printf(s)
  doc, err := html.Parse(strings.NewReader(s))
  if err != nil {
      // ...
  }
  var f func(*html.Node)
  f = func(n *html.Node) {
      if n.Type == html.ElementNode {
          fmt.Println(n)
        // do something?
      } else if n.Data == "a" {
          fmt.Println(n)
          // Do something with n...
      }
      for c := n.FirstChild; c != nil; c = c.NextSibling {
          f(c)
      }
  }
  f(doc)
}

