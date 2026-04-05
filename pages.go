package main

import (
	_ "embed"
	"fmt"
	"html"
	"strings"
)

func createHomePageHTML() string {
	bodyHTML := `<h1>Pilcrow</h1>
<p>Hi! I'm a software developer in Tokyo interested in auth and application security. I love thinking through tiny details and creating things from scratch. Some of my other hobbies are cooking, taking photos, and traveling.</p>
<p>Here are some projects I've worked on:</p>
<ul>
    <li><a href="https://basic-example.auth.pilcrowonpaper.com">basic-example.auth.pilcrowonpaper.com</a>: An example website that implements email address and password authentication following best practices.</li>
    <li><a href="https://lucia-auth.com">Lucia</a>: Former NPM package and current JavaScript/TypScript learning resource on session management.</li>
    <li><a href="https://arcticjs.dev">Arctic</a>: An OAuth 2.0 client NPM package.</li>
    <li><a href="https://oslojs.dev">Oslo</a>: A collection of auth-related NPM packages including packages for cryptographic-operations and the web authentication API (passkeys). </li>
    <li><a href="https://thecopenhagenbook.com">The Copenhagen Book</a>: A general guideline on implementing auth in web applications.</li>
</ul>
<p>Here are some links:</p>
<ul>
	<li>Blog: <a href="/blog">pilcrowonpaper.com/blog</a></li>
    <li>GitHub and GitHub Sponsors: <a href="https://github.com/pilcrowonpaper">github.com/pilcrowonpaper</a></li>
    <li>Twitter: <a href="https://x.com/pilcrowonpaper">x.com/pilcrowonpaper</a></li>
    <li>Bluesky: <a href="https://bsky.app/profile/pilcrowonpaper.com">bsky.app/profile/pilcrowonpaper.com</a></li>
    <li>Email: <a href="mailto:pilcrow@pilcrowonpaper.com">pilcrow@pilcrowonpaper.com</a></li>
	<li>RSS feed: <a href="/rss.xml">pilcrowonpaper.com/rss.xml</a></li>
</ul>
<p>The source code for this site is hosted at <a href="https://github.com/pilcrowonpaper/pilcrowonpaper.com">pilcrowonpaper/pilcrowonpaper.com</a> GitHub repository.</p>`
	pageHTML := createPageHTML("Pilcrow", "https://pilcrowonpaper.com", bodyHTML)

	return pageHTML
}

func createBlogPageHTML(blogItems []blogItemStruct) string {
	blogPostsListHTMLBuilder := strings.Builder{}
	blogPostsListHTMLBuilder.WriteString("<ul>")
	for _, blogItem := range blogItems {
		postLink := fmt.Sprintf("/blog/%d", blogItem.id)
		formattedPublishedDate := fmt.Sprintf("%s %d, %d", blogItem.publishedDate.month.String(), blogItem.publishedDate.day, blogItem.publishedDate.year)
		blogPostsListItemHTMLTemplate := `<li><a href="%s">%s (%s)</a></li>`
		blogPostsListItemHTML := fmt.Sprintf(blogPostsListItemHTMLTemplate, html.EscapeString(postLink), html.EscapeString(blogItem.title), html.EscapeString(formattedPublishedDate))
		blogPostsListHTMLBuilder.WriteString(blogPostsListItemHTML)
	}
	blogPostsListHTMLBuilder.WriteString("</ul>")
	blogPostsListHTML := blogPostsListHTMLBuilder.String()

	bodyHTMLTemplate := `<h1>Blog</h1>
<p>The RSS feed is published at <a href="/rss.xml">pilcrowonpaper.com/rss.xml</a>.</p>
%s`

	bodyHTML := fmt.Sprintf(bodyHTMLTemplate, blogPostsListHTML)

	pageHTML := createPageHTML("Blog", "https://pilcrowonpaper.com/blog", bodyHTML)

	return pageHTML
}

func createNotFoundPageHTML() string {
	bodyHTML := `<h1>Page not found</h1>
<p>The page you're looking for doesn't exist.</p>`
	pageHTML := createPageHTML("Page not found", "", bodyHTML)

	return pageHTML
}

//go:embed assets/stylesheet.css
var stylesheetCSS string

func createPageHTML(title string, link string, mainHTML string) string {
	canonicalLinkHTML := ""
	if link != "" {
		canonicalLinkHTML = fmt.Sprintf(`<link rel="canonical" href="%s" />`, html.EscapeString(link))
	}

	pageHTMLTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
	<title>%s</title>
	<meta name="description" content="Pilcrow's personal website." />

	<meta charset="utf-8" />
    <meta name="viewport" content="width=device-width" />

	<meta property="og:title" content="%s" />
	<meta property="og:type" content="website" />
	<meta property="og:locale" content="en_US" />
	<meta property="og:site_name" content="Pilcrow" />
	<meta property="og:description" content="Pilcrow's personal website." />
	<meta property="og:url" content="https://pilcrowonpaper.com" />
	<meta property="og:image" content="https://pilcrowonpaper.com/pilcrow.jpeg" />

	<meta name="twitter:card" content="summary">
    <meta name="twitter:site" content="@pilcrowonpaper">

	%s

	<style>%s</style>
</head>

<body>
	<header>
		<a id="pilcrow-link" href="/"><img id="pilcrow-icon" src="/pilcrow.jpeg" alt="Pilcrow"></a>
	</header>
	<main>%s</main>
</body>
</html>`

	pageHTML := fmt.Sprintf(pageHTMLTemplate, html.EscapeString(title), html.EscapeString(title), canonicalLinkHTML, stylesheetCSS, mainHTML)

	return pageHTML
}
