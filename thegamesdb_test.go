package thegamesdb

import (
	"testing"
)

func TestGetGamesList(t *testing.T) {
	games, err := GetGamesList("Crash Bandicoot", "Sony PlayStation", "Adventure")

	if err != nil {
		t.Error(err)
	}

	if len(*games) != 3 {
		t.Errorf("Expected %v games, got %v", 5, len(*games))
	}
}

func TestGetGame(t *testing.T) {

}

func TestGetArt(t *testing.T) {

}

func TestGetPlatformsList(t *testing.T) {

}

func TestGetPlatform(t *testing.T) {

}

func TestGetPlatformGames(t *testing.T) {

}

func TestPlatformGames(t *testing.T) {

}

func TestUpdates(t *testing.T) {

}

func TestUserRating(t *testing.T) {

}

func TestUserFavorites(t *testing.T) {

}
