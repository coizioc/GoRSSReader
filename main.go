package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Channel is a struct representing the <channel> tag in an RSS feed.
type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

// Item is a struct representing the <item> tag in an RSS feed.
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"guid"`
	Description string `xml:"description"`
	Date        string `xml:"pubDate"`
}

func main() {
	url := GetURL()

	data, err := GetDataFromURL(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	xmlData, err := ReadXMLData(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	PrintXMLData(xmlData)
}

// GetURL gets a url string from a command-line argument, or if none provided,
// asks the user to provide one.
func GetURL() string {
	var url string
	if len(os.Args) > 1 {
		url = os.Args[1]
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter URL of RSS Feed: ")
		url, _ = reader.ReadString('\n')
		// Remove trailing \n.
		url = url[:len(url)-1]
		// Remove trailing \r if exists (Windows support)
		if url[len(url)-1] == '\r' {
			url = url[:len(url)-1]
		}
	}
	return url
}

// GetDataFromURL requests data from a given URL and returns an array of bytes
// that make up the body of the response, if successful.
func GetDataFromURL(url string) ([]byte, error) {
	// Request the contents from the URL.
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	// Check if the request was successful.
	if res.StatusCode != http.StatusOK {
		statusErr := fmt.Errorf("Status error: %v\n", res.StatusCode)
		return []byte{}, statusErr
	}

	// Read the body of the result into an array of bytes.
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

// ReadXMLData takes in an array of bytes and tries to parse those values into a Channel struct.
func ReadXMLData(data []byte) (Channel, error) {
	// Declare anonymous struct rss representing the <rss> tag to wrap around Channel.
	var rss struct {
		Channel Channel `xml:"channel"`
	}

	// Retrieve the data from the byte array and put the result into rss.
	err := xml.Unmarshal(data, &rss)
	if err != nil {
		return Channel{}, err
	}

	return rss.Channel, nil
}

// PrintXMLData takes in a Channel struct, formats it, and prints out its members.
func PrintXMLData(xmlData Channel) {
	fmt.Println(xmlData.Title)
	for i, item := range xmlData.Items {
		fmt.Printf("    %2d. %s\n", i+1, item.Title)
		fmt.Printf("        %s\n", item.Description)
		fmt.Printf("        %s\n", item.Link)
		fmt.Printf("        %s\n", item.Date)
	}
}
