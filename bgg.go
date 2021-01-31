package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// BASEURL of boardgamegeek XMLAPIv2 https://boardgamegeek.com/wiki/page/BGG_XML_API2
	BASEURL = "https://boardgamegeek.com/xmlapi2/"
)

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
	Items      []Item   `xml:"item"`
}

type Item struct {
	Text          string          `xml:",chardata"`
	Type          string          `xml:"type,attr"`
	ID            string          `xml:"id,attr"`
	Names         []Name          `xml:"name"`
	YearPublished []Yearpublished `xml:"yearpublished"`
}

type Name struct {
	Text  string `xml:",chardata"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

type Yearpublished struct {
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

// SearchItems Used to search for items of given type. For strict searches exact should be true.
// Query: Returns all types of Items that match SEARCH_QUERY
// Item Types: rpgitem, videogame, boardgame, boardgameaccessory or boardgameexpansion
// Exact: Limit results to items that match the query exactly
func SearchItems(query string, gametype string, exact bool) Items {
	search := BASEURL + "search?query=" + query
	if gametype != "" {
		search = search + "&type=" + gametype
	}
	if exact == true {
		search = search + "&exact=1"
	}
	games := getSearchXML(search)
	return games
}
