package outletwebserver

import (
	"io/ioutil"
)

// WebPage represents a web page to display
type WebPage struct {
	Title   string
	Content []byte
}

// LoadPage to load a web page by its title
func LoadPage(title string) WebPage {
	filename := title + ".html"

	var page WebPage
	page.Title = title

	body, error := ioutil.ReadFile(filename)
	if error == nil {
		header := loadHeader()
		styles := loadStyles()

		content := header + "<style>" + styles + "</style>" + string(body)
		page.Content = []byte(content)
	}

	return page
}

func loadHeader() string {
	content, error := ioutil.ReadFile("header.html")
	if error == nil {
		return string(content)
	}
	return ""
}

func loadStyles() string {
	content, error := ioutil.ReadFile("styles.css")
	if error == nil {
		return string(content)
	}
	return ""
}
