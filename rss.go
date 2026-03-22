package main

import (
	"fmt"
	"html"
	"strings"
	"time"
)

func createRSSXML(blogItems []blogItemStruct) string {
	itemsXMLBuilder := strings.Builder{}

	for _, blogItem := range blogItems {
		template := `<item>
    <title>%s</title>
    <link>%s</link>
    <pubDate>%s</pubDate>
	<guid isPermaLink="false">%s</guid>
</item>`
		itemLink := fmt.Sprintf("https://pilcrowonpaper.com/blog/%d", blogItem.id)
		location, _ := time.LoadLocation("Asia/Tokyo")
		publishedDate := time.Date(blogItem.publishedDate.year, blogItem.publishedDate.month, blogItem.publishedDate.day, 12, 0, 0, 0, location)
		guid := fmt.Sprintf("pilcrowonpaper.com:%d", blogItem.id)
		itemXML := fmt.Sprintf(template, html.EscapeString(blogItem.title), html.EscapeString(itemLink), publishedDate.Format(time.RFC1123Z), html.EscapeString(guid))

		itemsXMLBuilder.WriteString(itemXML)
	}

	template := `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0" xmlns:atom="https://www.w3.org/2005/Atom">
<channel>
	<title>Pilcrow</title>
	<link>https://www.pilcrowonpaper.com</link>
	<description>Pilcrow's personal blog.</description>
	<language>en-us</language>
	<atom:link href="https://pilcrowonpaper.com/rss.xml" rel="self" type="application/rss+xml" />
	%s
</channel>
</rss>`

	rssXML := fmt.Sprintf(template, itemsXMLBuilder.String())

	return rssXML
}
