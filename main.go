package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

// # => id
// . => class
// x => tag
// x[y] => attr y in tag x, e.g. a[href] => the link in a tag

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	// only way to traverse children in 1 place?
	c.OnHTML(".lister-list", func(e *colly.HTMLElement) {
		e.DOM.Children().Each(func(i int, s *goquery.Selection) {
			posterSpans := s.Find(".posterColumn span[data-value]")
			posterSpans.Each(func(_ int, ss *goquery.Selection) {
				// looking for a more fucking elegant way to get it
				if v, _ := ss.Attr("name"); v == "rk" {
					v, _ = ss.Attr("data-value")
					fmt.Println("rank: ", v)
				}
			})
			fmt.Println("title: ", s.Find("a[title]").Text())
			fmt.Println("year: ", s.Find(".secondaryInfo").Text())
			rating := s.Find(".imdbRating strong")
			fmt.Println("rating: ", rating.Text())
			val, _ := rating.Attr("title")
			fmt.Println("rating: ", val)
			// kindergarten level seperator
			fmt.Println("-------------------------------")
		})
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.imdb.com/chart/top/")
}

// some helper functions for primitive Node traversal
func exploreNode(node *html.Node, level int) {
	builder := strings.Builder{}
	for i := 0; i < level; i++ {
		builder.WriteByte(' ')
	}
	builder.WriteString(fmt.Sprintf("(%d)[%s]%s", level, getAttr(node.Attr, "class"), node.Data))
	fmt.Println(builder.String())
	if node.FirstChild != nil {
		exploreNode(node.FirstChild, level+1)
	}
	if node.NextSibling != nil {
		exploreNode(node.NextSibling, level)
	}
}

func getAttr(attributes []html.Attribute, key string) string {
	for _, attr := range attributes {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
