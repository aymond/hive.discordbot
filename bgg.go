package main

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"

	"net/http"
	"os"
	"time"

	"golang.org/x/net/html"
)

const (
	// BASEURL of boardgamegeek XMLAPIv2 https://boardgamegeek.com/wiki/page/BGG_XML_API2
	BASEURL = "https://boardgamegeek.com/xmlapi2/"
)

// BGGHTMLMeta is the meta info of the BGG Item webpage.
// <meta property="og:title" content="Brass: Birmingham" />
// <meta property="og:image" content="https://cf.geekdo-images.com/x3zxjr-Vw5iU4yDPg70Jgw__opengraph_left/img/lYAj3vj2GtibZtG_62bVHD5Xy8c=/fit-in/445x445/filters:strip_icc()/pic3490053.jpg" />
// <meta property="og:url" content="https://boardgamegeek.com/boardgame/224517/brass-birmingham" />
// <meta property="og:site_name" content="BoardGameGeek" />
// <meta property="og:description" content="Build networks, grow industries, and navigate the world of the Industrial Revolution." />
type BGGHTMLMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	SiteName    string `json:"site_name"`
	URL         string `json:"url"`
}

// Items struct
// https://boardgamegeek.com/xmlapi2/search?query=settlers
// <items total="190" termsofuse="https://boardgamegeek.com/xmlapi/termsofuse">
// 		<item type="boardgame" id="17419">
// 			<name type="alternate" value="10th Anniversary Treasure Chest Set: Settlers of Catan"/>
// 			<yearpublished value="2005"/>
// 		</item>
// 		<item type="boardgame" id="12004">
// 			<name type="primary" value="Candamir: The First Settlers"/>
// 			<yearpublished value="2004"/>
// 		</item>
// </items>
type Items struct {
	XMLName    xml.Name `xml:"items"`
	Text       string   `xml:",chardata"`
	Total      string   `xml:"total,attr"`
	Termsofuse string   `xml:"termsofuse,attr"`
	Items      []item   `xml:"item"`
}

type item struct {
	Text          string          `xml:",chardata"`
	Type          string          `xml:"type,attr"`
	ID            string          `xml:"id,attr"`
	Names         []name          `xml:"name"`
	YearPublished []yearpublished `xml:"yearpublished"`
}

type name struct {
	Text  string `xml:",chardata"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

type yearpublished struct {
	Text  string `xml:",chardata"`
	Value string `xml:"value,attr"`
}

// getSearchXML Private function that queries the API prepared in SearchItems.
func getSearchXML(url string) Items {
	log.Println(url)
	response, err := http.Get(url)
	var v Items
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		data, err := ioutil.ReadAll(response.Body)
		err = xml.Unmarshal([]byte(data), &v)
		if err != nil {
			log.Fatal(err)
		}
	}
	return v
}

// BGGSearchItems Used to search for items of given type. For strict searches exact should be true.
// Query: Returns all types of Items that match SEARCH_QUERY
// Item Types: rpgitem, videogame, boardgame, boardgameaccessory or boardgameexpansion
// Exact: Limit results to items that match the query exactly
func BGGSearchItems(query string, gametype string, exact bool) (games Items, search string) {
	search = BASEURL + "search?query=" + query
	if gametype != "" {
		search = search + "&type=" + gametype
	}
	if exact == true {
		search = search + "&exact=1"
	}
	games = getSearchXML(search)
	return games, search
}

// BGGGetItemPage is used to request the item page. Will need to parse the page to get any interesting content.
// <meta property="og:image" content="https://cf.geekdo-images.com/x3zxjr-Vw5iU4yDPg70Jgw__opengraph_left/img/lYAj3vj2GtibZtG_62bVHD5Xy8c=/fit-in/445x445/filters:strip_icc()/pic3490053.jpg" />
// <meta property="og:url" content="https://boardgamegeek.com/boardgame/224517/brass-birmingham" />
// <meta property="og:description" content="Build networks, grow industries, and navigate the world of the Industrial Revolution." />
func BGGGetItemPage(requestURL string) (meta *BGGHTMLMeta) {

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", "Not Firefox")
	// Make request
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	// get the HTML Metadata
	meta = extract(response.Body)

	defer response.Body.Close()

	// Copy data from the response to standard output
	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return meta
}

// BGGGetThumbnail Get the thumbnail
func BGGGetThumbnail(requestURL string) []byte {

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", "Not Firefox")
	// Make request
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	// Copy data from the response to standard output
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Trouble reading reesponse body!")
	}

	return contents
}

func extract(resp io.Reader) *BGGHTMLMeta {
	z := html.NewTokenizer(resp)

	titleFound := false

	hm := new(BGGHTMLMeta)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return hm
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == `body` {
				return hm
			}
			if t.Data == "title" {
				titleFound = true
			}
			if t.Data == "meta" {
				desc, ok := extractMetaProperty(t, "description")
				if ok {
					hm.Description = desc
				}

				ogTitle, ok := extractMetaProperty(t, "og:title")
				if ok {
					hm.Title = ogTitle
				}

				ogDesc, ok := extractMetaProperty(t, "og:description")
				if ok {
					hm.Description = ogDesc
				}

				ogImage, ok := extractMetaProperty(t, "og:image")
				if ok {
					hm.Image = ogImage
				}

				ogSiteName, ok := extractMetaProperty(t, "og:site_name")
				if ok {
					hm.SiteName = ogSiteName
				}
			}
		case html.TextToken:
			if titleFound {
				t := z.Token()
				hm.Title = t.Data
				titleFound = false
			}
		}
	}
	return hm
}

// extractMetaProperty gets the meta property out of the html document.
func extractMetaProperty(t html.Token, prop string) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == "property" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return
}
