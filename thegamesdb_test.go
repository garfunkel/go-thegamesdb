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
	game, err := GetGame(2, "Crysis", "Crysis", "PC")

	if err != nil {
		t.Error(err)
	}

	if game.Name != "Crysis" {
		t.Errorf("Expected game 'Crysis', got '%s'", game.Name)
	}
}

func TestGetArt(t *testing.T) {
	images, err := GetArt(2)

	if err != nil {
		t.Error(err)
	}

	if len((*images).FanArtImages) == 0 {
		t.Error("Expected images, found 0")
	}
}

func TestGetPlatformsList(t *testing.T) {
	platforms, err := GetPlatformsList()

	if err != nil {
		t.Error(err)
	}

	if len(*platforms) == 0 {
		t.Error("Expected platforms, found 0")
	}
}

func TestGetPlatform(t *testing.T) {
	platform, err := GetPlatform(15)

	if err != nil {
		t.Error(err)
	}

	if platform.Name != "Microsoft Xbox 360" {
		t.Errorf("Expected name 'Microsoft Xbox 360', got '%v'", platform.Name)
	}
}

func TestGetPlatformGames(t *testing.T) {
	games, err := GetPlatformGames(1)

	if err != nil {
		t.Error(err)
	}

	if len(*games) == 0 {
		t.Error("Expected games, found 0")
	}
}

func TestGetUpdates(t *testing.T) {
	games, err := GetUpdates(2000000)

	if err != nil {
		t.Error(err)
	}

	if len(*games) == 0 {
		t.Error("Expected games, found 0")
	}
}

func TestGetUserRating(t *testing.T) {
	rating, err := GetUserRating(1, "58536D31278176DA")

	if err != nil {
		t.Error(err)
	}

	if rating == 0 {
		t.Error("Expected positive rating, got 0")
	}
}

func TestGetUserFavourites(t *testing.T) {
	_, err := GetUserFavourites("58536D31278176DA")

	if err != nil {
		t.Error(err)
	}
}
