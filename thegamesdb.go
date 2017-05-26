package thegamesdb

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type GameDate time.Time

type Game struct {
	ID          uint      `xml:"id"`
	GameTitle   string    `xml:"GameTitle"`
	ReleaseData GameDate  `xml:"ReleaseDate"`
	Platform    string    `xml:"Platform"`
}

type Games []Game

type getGamesListData struct {
	Games Games `xml:"Game"`
}

func (gameDate GameDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	s := ""

	err = d.DecodeElement(&s, &start)

	if err != nil {
		return
	}

	if s == "" {
		return
	}

	parsedDate, err := time.Parse("01/02/2006", s)

	if err != nil {
		return
	}

	gameDate = GameDate(parsedDate)

	return
}

func GetGamesList(name, platform, genre string) (games *Games, err error) {
	u := fmt.Sprintf("http://thegamesdb.net/api/GetGamesList.php?name=%v", url.QueryEscape(name))

	if platform != "" {
		u += fmt.Sprintf("&platform=%v", url.QueryEscape(platform))
	}

	if genre != "" {
		u += fmt.Sprintf("&genre=%v", url.QueryEscape(genre))
	}

	resp, err := http.Get(u)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	data := getGamesListData{}

	err = xml.Unmarshal(body, &data)

	games = &data.Games

	return
}

func GetGame() {

}

func GetArt() {

}

func GetPlatformsList() {

}

func GetPlatform() {

}

func GetPlatformGames() {

}

func PlatformGames() {

}

func Updates() {

}

func UserRating() {

}

func UserFavorites() {

}
