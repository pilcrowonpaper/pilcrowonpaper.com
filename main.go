package main

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"slices"
	"time"

	"github.com/pilcrowonpaper/go-json"
)

//go:embed assets/pilcrow.jpg
var pilcrowJPEGImage []byte

//go:embed blog.json
var blogJSON string

//go:embed assets/_redirects
var redirectsText string

func main() {
	blogJSONArray, err := json.ParseArray(blogJSON)
	if err != nil {
		fmt.Printf("Failed to parse blog.json file as a JSON array: %s\n", err.Error())
		os.Exit(0)
	}
	blogItems, err := mapJSONArrayToBlogItems(blogJSONArray)
	if err != nil {
		fmt.Printf("Failed to map json array to blog items: %s\n", err.Error())
		os.Exit(0)
	}

	err = os.RemoveAll("dist")
	if err != nil {
		fmt.Printf("Failed to remove 'dist' directory and its content: %s\n", err.Error())
		os.Exit(0)
	}

	err = os.MkdirAll("dist", os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create directory 'dist': %s\n", err.Error())
		os.Exit(0)
	}

	homePageHTML := createHomePageHTML()
	err = os.WriteFile("dist/index.html", []byte(homePageHTML), os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to write file 'dist/index.html': %s\n", err.Error())
		os.Exit(0)
	}

	notFoundPageHTML := createNotFoundPageHTML()
	err = os.WriteFile("dist/404.html", []byte(notFoundPageHTML), os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to write file 'dist/404.html': %s\n", err.Error())
		os.Exit(0)
	}

	err = os.WriteFile("dist/pilcrow.jpeg", pilcrowJPEGImage, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to write file 'pilcrow.jpeg': %s\n", err.Error())
		os.Exit(0)
	}

	err = os.Mkdir("dist/blog", os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create directory 'dist/blog': %s\n", err.Error())
		os.Exit(0)
	}

	slices.SortFunc(blogItems, func(a, b blogItemStruct) int {
		return b.id - a.id
	})

	blogPageHTML := createBlogPageHTML(blogItems)
	err = os.WriteFile("dist/blog/index.html", []byte(blogPageHTML), os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to write file 'dist/blog/index.html': %s\n", err.Error())
		os.Exit(0)
	}

	for _, blogItem := range blogItems {
		blogContentHTMLFileName := fmt.Sprintf("%d.html", blogItem.id)
		blogContentHTMLFilePath := path.Join("blog", blogContentHTMLFileName)
		blogContentHTMLBytes, err := os.ReadFile(blogContentHTMLFilePath)
		if err != nil {
			fmt.Printf("Failed to read file '%s': %s\n", blogContentHTMLFilePath, err.Error())
			os.Exit(0)
		}

		blogLink := fmt.Sprintf("https://pilcrowonpaper.com/blog/%d", blogItem.id)

		blogPageHTML := createPageHTML(blogItem.title, blogLink, string(blogContentHTMLBytes))

		distFilePath := path.Join("dist", "blog", blogContentHTMLFileName)
		err = os.WriteFile(distFilePath, []byte(blogPageHTML), os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to write file '%s': %s\n", distFilePath, err.Error())
			os.Exit(0)
		}
	}

	imagesDirectoryFS := os.DirFS("images")

	err = os.CopyFS("dist/images", imagesDirectoryFS)
	if err != nil {
		fmt.Printf("Failed to copy images directory to 'dist': %s\n", err.Error())
		os.Exit(0)
	}

	rssXML := createRSSXML(blogItems)
	err = os.WriteFile("dist/rss.xml", []byte(rssXML), os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to write file 'rss.xml': %s\n", err.Error())
		os.Exit(0)
	}

	err = os.WriteFile("dist/_redirects", []byte(redirectsText), os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to write file '_redirects': %s\n", err.Error())
		os.Exit(0)
	}
}

type blogItemStruct struct {
	id            int
	title         string
	publishedDate dateStruct
}

func mapJSONArrayToBlogItems(blogItemsJSONArray json.ArrayStruct) ([]blogItemStruct, error) {
	blogItems := []blogItemStruct{}
	for i := range blogItemsJSONArray.Length() {
		blogItemJSONObject, err := blogItemsJSONArray.GetJSONObject(i)
		if err != nil {
			return nil, fmt.Errorf("failed to get json object at index %d: %s", i, err.Error())
		}
		blogItem, err := mapJSONObjectToBlogItem(blogItemJSONObject)
		if err != nil {
			return nil, fmt.Errorf("failed to map json object to blog item: %s", err.Error())
		}
		blogItems = append(blogItems, blogItem)
	}
	return blogItems, nil
}

func mapJSONObjectToBlogItem(blogItemJSONObject json.ObjectStruct) (blogItemStruct, error) {
	id, err := blogItemJSONObject.GetInt("id")
	if err != nil {
		return blogItemStruct{}, fmt.Errorf("failed to get 'id' integer value: %s", err.Error())
	}
	title, err := blogItemJSONObject.GetString("title")
	if err != nil {
		return blogItemStruct{}, fmt.Errorf("failed to get 'title' string value: %s", err.Error())
	}
	publishedDateString, err := blogItemJSONObject.GetString("published_date")
	if err != nil {
		return blogItemStruct{}, fmt.Errorf("failed to get 'published_date' string value: %s", err.Error())
	}
	publishedDate, err := parseDateString(publishedDateString)
	if err != nil {
		return blogItemStruct{}, fmt.Errorf("failed to parse 'published_date' value as a date: %s", err.Error())
	}

	blogItem := blogItemStruct{
		id:            id,
		title:         title,
		publishedDate: publishedDate,
	}
	return blogItem, nil
}

type dateStruct struct {
	year  int
	month time.Month
	day   int
}

func parseDateString(dateString string) (dateStruct, error) {
	chars := []rune(dateString)
	if len(chars) != 10 {
		return dateStruct{}, fmt.Errorf("invalid length")
	}

	if !isDigit(chars[0]) || !isDigit(chars[1]) || !isDigit(chars[2]) || !isDigit(chars[3]) {
		return dateStruct{}, fmt.Errorf("invalid character")
	}
	year := int(chars[0]-'0')*1000 + int(chars[1]-'0')*100 + int(chars[2]-'0')*10 + int(chars[3]-'0')

	if chars[4] != '-' {
		return dateStruct{}, fmt.Errorf("invalid character")
	}

	if !isDigit(chars[5]) || !isDigit(chars[6]) {
		return dateStruct{}, fmt.Errorf("invalid character")
	}
	month := time.Month(int(chars[5]-'0')*100 + int(chars[6]-'0'))
	if month < 1 || month > 12 {
		return dateStruct{}, fmt.Errorf("invalid month")
	}

	if chars[7] != '-' {
		return dateStruct{}, fmt.Errorf("invalid character")
	}

	if !isDigit(chars[8]) || !isDigit(chars[9]) {
		return dateStruct{}, fmt.Errorf("invalid character")
	}
	day := int(chars[8]-'0')*10 + int(chars[9]-'0')
	if day < 1 || day > 31 {
		return dateStruct{}, fmt.Errorf("invalid day")
	}

	date := dateStruct{year, month, day}

	return date, nil
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
