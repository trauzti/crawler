package gocrawler

import (
    "testing"
)

func TestCanonicalizer(t *testing.T) {
    var url1, url2, curl1, curl2 string
    url1 = "http://www.mbl.is/mm?trausti=10"
    url2 = "http://www.mbl.is:80/mm"
    curl1 = canonicalizeUrl(url1)
    curl2 = canonicalizeUrl(url2)
    if curl1 != curl2  {
        t.Error(curl1, "vs.", curl2)
    }
    url1 = "http://www.mbl.is/mm"
    url2 = "http://mbl.is:80/mm"
    curl1 = canonicalizeUrl(url1)
    curl2 = canonicalizeUrl(url2)
    if curl1 != curl2  {
        t.Error(curl1, "vs.", curl2)
    }
    url1 = "http://www.mbl.is"
    url2 = "http://mbl.is:80/"
    curl1 = canonicalizeUrl(url1)
    curl2 = canonicalizeUrl(url2)
    if curl1 != curl2  {
        t.Error(curl1, "vs.", curl2)
    }
   //assert canonical(mbl.is:80 == mbl.is)
}

func TestMakeAbsoluteUrl(t *testing.T) {
    var add1 = "/mm/frettir"
    var base = "http://mbl.is"
    result := "http://mbl.is/mm/frettir"
    if makeAbsoluteUrl(base, add1) != result  {
        t.Error(add1, base, "vs.", result)
    }
}

func TestMakeRelativeUrl(t *testing.T) {
    mbl := "http://mbl.is/mm/frettir"
    rel := "/mm/frettir"
    res := makeRelativeUrl(mbl)
    if res != rel  {
        t.Error(res, "vs.", rel)
    }
}
