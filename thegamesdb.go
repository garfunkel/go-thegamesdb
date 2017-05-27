package thegamesdb

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// APIBaseURL is the base API URL for TheGamesDB.
const APIBaseURL string = "http://thegamesdb.net/api/"

// GameDate represents a date.
type GameDate time.Time

// Game is the data structure for game data.
type Game struct {
	ID           uint     `xml:"id"`
	Name         string   `xml:"GameTitle"`
	PlatformID   string   `xml:"PlatformId"`
	Platform     string   `xml:"Platform"`
	ReleaseDate  GameDate `xml:"ReleaseDate"`
	Overview     string   `xml:"Overview"`
	ESRB         string   `xml:"ESRB"`
	Genres       []string `xml:"Genres>genre"`
	Players      string   `xml:"Players"`
	CoOp         string   `xml:"Co-op"`
	YouTube      string   `xml:"Youtube"`
	Publisher    string   `xml:"Publisher"`
	Developer    string   `xml:"Developer"`
	Rating       float32  `xml:"Rating"`
	SimilarGames Games    `xml:"Similar>Game"`
	Images       Images   `xml:"Images"`
}

// Images is the data structure for a group of image assets.
type Images struct {
	FanArtImages FanArtImages `xml:"fanart"`
	BoxArtImages BoxArtImages `xml:"boxart"`
	Banners      Banners      `xml:"banner"`
	Screenshots  Screenshots  `xml:"screenshot"`
	ClearLogos   ClearLogos   `xml:"clearlogo"`
}

// FanArtImage is the data structure for fan art.
type FanArtImage struct {
	URL       string
	Width     int
	Height    int
	Thumbnail string
}

// BoxArtImage is the data structure for box art.
type BoxArtImage struct {
	URL       string `xml:".chardata"`
	Width     int    `xml:"width,attr"`
	Height    int    `xml:"height,attr"`
	Thumbnail string `xml:"thumb,attr"`
	Side      string `xml:"side,attr"`
}

// Banner is the data structure for banner images.
type Banner struct {
	URL    string `xml:",chardata"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// Screenshot is the data structure for screenshot images.
type Screenshot struct {
	URL       string
	Width     int
	Height    int
	Thumbnail string
}

// ClearLogo is the data structure for clear logo images.
type ClearLogo struct {
	URL    string `xml:",chardata"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// Platform is the data structure for platform data.
type Platform struct {
	ID             uint    `xml:"id"`
	Name           string  `xml:"name"`
	Alias          string  `xml:"alias"`
	Console        string  `xml:"console"`
	Controller     string  `xml:"controller"`
	Overview       string  `xml:"overview"`
	Developer      string  `xml:"developer"`
	Manufacturer   string  `xml:"manufacturer"`
	CPU            string  `xml:"cpu"`
	Memory         string  `xml:"memory"`
	Graphics       string  `xml:"graphics"`
	Sound          string  `xml:"sound"`
	Display        string  `xml:"display"`
	Media          string  `xml:"media"`
	MaxControllers int     `xml:"maxcontrollers"`
	Rating         float32 `xml:"Rating"`
	Images         Images  `xml:"Images"`
}

// Games is an array of games.
type Games []Game

// FanArtImages is an array of fan art images.
type FanArtImages []FanArtImage

// BoxArtImages is an array of box art images.
type BoxArtImages []BoxArtImage

// Banners is an array of banner images.
type Banners []Banner

// Screenshots is an array of screenshot images.
type Screenshots []Screenshot

// ClearLogos is an array of clear logo images.
type ClearLogos []ClearLogo

// Platforms is an array of platforms.
type Platforms []Platform

type getGamesListData struct {
	Games Games `xml:"Game"`
}

type getGameData struct {
	BaseImgURL string `xml:"baseImgUrl"`
	Game       Game   `xml:"Game"`
}

type getArtData struct {
	BaseImgURL string `xml:"baseImgUrl"`
	Images     Images `xml:"Images"`
}

type getPlatformsListData struct {
	BasePlatformURL string    `xml:"basePlatformUrl"`
	Platforms       Platforms `xml:"Platforms>Platform"`
}

type getPlatformData struct {
	BaseImgURL string   `xml:"baseImgUrl"`
	Platform   platform `xml:"Platform"`
}

type getPlatformGamesData struct {
	Games Games `xml:"Game"`
}

type getUpdatesData struct {
	Time  int    `xml:"Time"`
	Games []uint `xml:"Game"`
}

type getUserRatingData struct {
	Rating float32 `xml:"game>Rating"`
}

type getUserFavouritesData struct {
	Games []uint `xml:"Game"`
}

type fanArt struct {
	Original struct {
		URL    string `xml:",chardata"`
		Width  int    `xml:"width,attr"`
		Height int    `xml:"height,attr"`
	} `xml:"original"`
	Thumbnail string `xml:"thumb"`
}

type screenshot struct {
	Original struct {
		URL    string `xml:",chardata"`
		Width  int    `xml:"width,attr"`
		Height int    `xml:"height,attr"`
	} `xml:"original"`
	Thumbnail string `xml:"thumb"`
}

type platform struct {
	Platform
	Name string `xml:"Platform"`
}

// UnmarshalXML unmarshals XML data for a game date.
func (gameDate *GameDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
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
		parsedDate, err = time.Parse("2006", s)

		if err != nil {
			return
		}
	}

	*gameDate = GameDate(parsedDate)

	return
}

// UnmarshalXML unmarshals XML data for a fan art image.
func (fanArtImage *FanArtImage) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	f := fanArt{}
	err = d.DecodeElement(&f, &start)

	if err != nil {
		return
	}

	fanArtImage.URL = f.Original.URL
	fanArtImage.Width = f.Original.Width
	fanArtImage.Height = f.Original.Height
	fanArtImage.Thumbnail = f.Thumbnail

	return
}

// UnmarshalXML unmarshals XML data for a screenshot image.
func (ss *Screenshot) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	s := screenshot{}
	err = d.DecodeElement(&s, &start)

	if err != nil {
		return
	}

	ss.URL = s.Original.URL
	ss.Width = s.Original.Width
	ss.Height = s.Original.Height
	ss.Thumbnail = s.Thumbnail

	return
}

func (images *Images) applyBaseImgURL(baseImgURL string) {
	for i, x := range images.FanArtImages {
		if images.FanArtImages[i].URL != "" {
			images.FanArtImages[i].URL = baseImgURL + x.URL
		}

		if images.FanArtImages[i].Thumbnail != "" {
			images.FanArtImages[i].Thumbnail = baseImgURL + x.Thumbnail
		}
	}

	for i, x := range images.BoxArtImages {
		if images.BoxArtImages[i].URL != "" {
			images.BoxArtImages[i].URL = baseImgURL + x.URL
		}

		if images.BoxArtImages[i].Thumbnail != "" {
			images.BoxArtImages[i].Thumbnail = baseImgURL + x.Thumbnail
		}
	}

	for i, x := range images.Banners {
		if images.Banners[i].URL != "" {
			images.Banners[i].URL = baseImgURL + x.URL
		}
	}

	for i, x := range images.Screenshots {
		if images.Screenshots[i].URL != "" {
			images.Screenshots[i].URL = baseImgURL + x.URL
		}

		if images.Screenshots[i].Thumbnail != "" {
			images.Screenshots[i].Thumbnail = baseImgURL + x.Thumbnail
		}
	}

	for i, x := range images.ClearLogos {
		if images.ClearLogos[i].URL != "" {
			images.ClearLogos[i].URL = baseImgURL + x.URL
		}
	}
}

// GetGamesList returns a list of games for a given platform.
func GetGamesList(name, platform, genre string) (games *Games, err error) {
	if name == "" && platform == "" && genre == "" {
		err = errors.New("No valid parameters specified")

		return
	}

	u := APIBaseURL + "GetGamesList.php?"

	if name != "" {
		u = strings.Join([]string{u, fmt.Sprintf("name=%v", url.QueryEscape(name))}, "&")
	}

	if platform != "" {
		u = strings.Join([]string{u, fmt.Sprintf("platform=%v", url.QueryEscape(platform))}, "&")
	}

	if genre != "" {
		u = strings.Join([]string{u, fmt.Sprintf("genre=%v", url.QueryEscape(genre))}, "&")
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

	if err != nil {
		return
	}

	games = &data.Games

	return
}

// GetGame returns all information about a game.
func GetGame(id int, name, exactName, platform string) (game *Game, err error) {
	if id == 0 && name == "" {
		err = errors.New("No valid id or name parameter specified")

		return
	}

	u := APIBaseURL + "GetGame.php?"

	if id != 0 {
		u = strings.Join([]string{u, fmt.Sprintf("id=%v", id)}, "&")
	}

	if name != "" {
		u = strings.Join([]string{u, fmt.Sprintf("name=%v", url.QueryEscape(name))}, "&")
	}

	if exactName != "" {
		u = strings.Join([]string{u, fmt.Sprintf("exactname=%v", url.QueryEscape(exactName))}, "&")
	}

	if platform != "" {
		u = strings.Join([]string{u, fmt.Sprintf("platform=%v", url.QueryEscape(platform))}, "&")
	}

	resp, err := http.Get(u)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	data := getGameData{}
	err = xml.Unmarshal(body, &data)

	if err != nil {
		return
	}

	game = &data.Game

	game.Images.applyBaseImgURL(data.BaseImgURL)

	return
}

// GetArt returns all information about a game's artwork.
func GetArt(id int) (images *Images, err error) {
	if id == 0 {
		err = errors.New("No valid id parameter specified")

		return
	}

	u := APIBaseURL + fmt.Sprintf("GetArt.php?id=%v", id)
	resp, err := http.Get(u)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	data := getArtData{}
	err = xml.Unmarshal(body, &data)

	if err != nil {
		return
	}

	images = &data.Images

	images.applyBaseImgURL(data.BaseImgURL)

	return
}

// GetPlatformsList returns a list of available platforms.
func GetPlatformsList() (platforms *Platforms, err error) {
	u := APIBaseURL + "GetPlatformsList.php"
	resp, err := http.Get(u)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	data := getPlatformsListData{}
	err = xml.Unmarshal(body, &data)

	if err != nil {
		return
	}

	platforms = &data.Platforms

	return
}

// GetPlatform returns all information about a platform.
func GetPlatform(id int) (platform *Platform, err error) {
	if id == 0 {
		err = errors.New("No valid id parameter specified")

		return
	}

	u := APIBaseURL + fmt.Sprintf("GetPlatform.php?id=%v", id)
	resp, err := http.Get(u)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	data := getPlatformData{}
	err = xml.Unmarshal(body, &data)

	if err != nil {
		return
	}

	platform = &data.Platform.Platform
	platform.Name = data.Platform.Name

	return
}

// GetPlatformGames returns a list of games for a given platform.
func GetPlatformGames(id int) (games *Games, err error) {
	if id == 0 {
		err = errors.New("No valid id parameter specified")

		return
	}

	u := APIBaseURL + fmt.Sprintf("GetPlatformGames.php?platform=%v", id)
	resp, err := http.Get(u)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	data := getPlatformGamesData{}
	err = xml.Unmarshal(body, &data)

	if err != nil {
		return
	}

	games = &data.Games

	return
}

// GetUpdates returns all TheGamesDB updates since a number of seconds ago.
func GetUpdates(since uint) (games *Games, err error) {
	u := APIBaseURL + fmt.Sprintf("Updates.php?time=%v", since)
	resp, err := http.Get(u)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	data := getUpdatesData{}
	err = xml.Unmarshal(body, &data)

	if err != nil {
		return
	}

	games = new(Games)
	*games = make(Games, len(data.Games))

	for i, x := range data.Games {
		(*games)[i].ID = x
	}

	return
}

// GetUserRating returns the user's rating for a game.
func GetUserRating(gameID int, apiID string) (rating float32, err error) {
	if gameID == 0 {
		err = errors.New("No valid gameID parameter specified")

		return
	}

	if apiID == "" {
		err = errors.New("No valid apiID parameter specified")

		return
	}

	u := APIBaseURL + fmt.Sprintf("User_Rating.php?itemid=%v&accountid=%v", gameID, apiID)
	resp, err := http.Get(u)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	data := getUserRatingData{}
	err = xml.Unmarshal(body, &data)

	if err != nil {
		return
	}

	rating = data.Rating

	return
}

// GetUserFavourites returns the user's favourite games.
func GetUserFavourites(apiID string) (games *Games, err error) {
	if apiID == "" {
		err = errors.New("No valid apiID parameter specified")

		return
	}

	u := APIBaseURL + fmt.Sprintf("User_Favorites.php?accountid=%v", apiID)
	resp, err := http.Get(u)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	data := getUserFavouritesData{}
	err = xml.Unmarshal(body, &data)

	if err != nil {
		return
	}

	games = new(Games)
	*games = make(Games, len(data.Games))

	for i, x := range data.Games {
		(*games)[i].ID = x
	}

	return
}
