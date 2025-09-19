package main

import (
	"flag"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/supersonic-app/go-subsonic/subsonic"
)

var (
	subsonicServer = flag.String("subsonic-server", "http://127.0.0.1:1080", "Subsonic server base address")
	subsonicUser   = flag.String("subsonic-user", "admin", "Subsonic user name")
	subsonicPass   = flag.String("subsonic-pass", "admin", "Subsonic password")

	playlistName = flag.String("playlist-name", "Daily Mix", "Playlist to create/update")
)

const (
	songBatchSize = 1000
)

func connect() (*subsonic.Client, error) {
	client := &subsonic.Client{
		Client:     http.DefaultClient,
		BaseUrl:    *subsonicServer,
		User:       *subsonicUser,
		ClientName: "BASS",
	}

	if *subsonicPass != "" {
		if err := client.Authenticate(*subsonicPass); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func songs(client *subsonic.Client) ([]*subsonic.Child, error) {
	slog.Info("Retrieving all songs from subsonic server")

	offset := 0
	params := map[string]string{
		"artistCount": "0",
		"albumCount":  "0",
		"songCount":   strconv.Itoa(songBatchSize),
		"songOffset":  strconv.Itoa(offset),
	}

	var children []*subsonic.Child
	for {
		res, err := client.Search3("", params)
		if err != nil {
			return nil, err
		}

		children = append(children, res.Song...)
		slog.Debug("Retrieved batch of songs", "size", len(res.Song), "total", len(children), "offset", offset)
		if len(res.Song) == songBatchSize {
			offset += songBatchSize
			params["songOffset"] = strconv.Itoa(offset)
		} else {
			break
		}
	}

	return children, nil
}

func updatePlaylist(client *subsonic.Client, songs []*subsonic.Child) error {
	playlists, err := client.GetPlaylists(nil)
	if err != nil {
		return err
	}

	for i := range playlists {
		if playlists[i].Name == *playlistName {
			slog.Info("Updating playlist", "name", playlists[i].Name, "id", playlists[i].ID)
			return updatePlaylistSongs(client, playlists[i], songs)
		}
	}

	slog.Info("Creating new playlist", "name", *playlistName)
	playlist, err := client.CreatePlaylist(map[string]string{"name": *playlistName})
	if err != nil {
		return err
	}

	return updatePlaylistSongs(client, playlist, songs)
}

func updatePlaylistSongs(client *subsonic.Client, playlist *subsonic.Playlist, songs []*subsonic.Child) error {
	var indicesToRemove []int
	for i := 0; i < playlist.SongCount; i++ {
		indicesToRemove = append(indicesToRemove, i)
	}
	var songsToAdd []string
	for i := range songs {
		songsToAdd = append(songsToAdd, songs[i].ID)
	}
	return client.UpdatePlaylistTracks(playlist.ID, songsToAdd, indicesToRemove)
}
