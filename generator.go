package main

import (
	"flag"
	"log/slog"
	"maps"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/supersonic-app/go-subsonic/subsonic"
)

var (
	targetLength = flag.Duration("target-length", 18*time.Hour, "Target length of the playlist")

	weightStarred    = flag.Float64("weight-starred", 20, "Weight of starred tracks")
	weightRatedOne   = flag.Float64("weight-rated-1", 0, "Weight of tracks rated 1/5")
	weightRatedTwo   = flag.Float64("weight-rated-2", 0.5, "Weight of tracks rated 2/5")
	weightRatedThree = flag.Float64("weight-rated-3", 1, "Weight of tracks rated 3/5")
	weightRatedFour  = flag.Float64("weight-rated-4", 5, "Weight of tracks rated 4/5")
	weightRatedFive  = flag.Float64("weight-rated-5", 10, "Weight of tracks rated 5/5")

	weightNeverPlayed      = flag.Float64("weight-never-played", 5, "Weight of tracks never played")
	weightRarelyPlayed     = flag.Float64("weight-rarely-played", 3, "Weight of tracks played a lot less often than average")
	weightFrequentlyPlayed = flag.Float64("weight-frequently-played", 1, "Weight of tracks played a lot more often than average")

	weightSameAlbum  = flag.Float64("weight-same-album", 20, "Weight of tracks in the same album to the last song")
	weightSameArtist = flag.Float64("weight-same-artist", 15, "Weight of tracks by the same artist as the last song")
	weightSameGenre  = flag.Float64("weight-same-genre", 10, "Weight of tracks in the same genre as the last song")

	weightDuplicate = flag.Float64("weight-duplicate", 0, "Weight of tracks that are already in the playlist")
)

func generate(songs []*subsonic.Child) []*subsonic.Child {
	var output []*subsonic.Child
	var length = time.Duration(0)
	var lastSong *subsonic.Child
	var averagePlayCount = averagePlays(songs)
	slog.Debug("Prepared to generate", "averagePlayCount", averagePlayCount)

	var weights = make([]float64, len(songs))
	for length < *targetLength {
		var total = float64(0)
		for i := range songs {
			weights[i] = weight(output, averagePlayCount, lastSong, songs[i])
			total += weights[i]
		}

		selected := rand.Float64() * total
		for i := range songs {
			selected -= weights[i]
			if selected <= 0 {
				lastSong = songs[i]
				output = append(output, songs[i])
				length += time.Duration(songs[i].Duration) * time.Second
				slog.Debug("Picked song", "index", i, "duration", songs[i].Duration, "total length", length, "target length", *targetLength, "title", songs[i].Title, "artist", songs[i].Artist, "genre", songs[i].Genre, "weight", weights[i])
				break
			}
		}
	}

	return output
}

func randomGenre(songs []*subsonic.Child) string {
	genres := make(map[string]int)
	for i := range songs {
		genres[songs[i].Genre]++
	}
	values := slices.Collect(maps.Keys(genres))
	return values[rand.IntN(len(values))]
}

func randomArtist(songs []*subsonic.Child) string {
	artists := make(map[string]int)
	for i := range songs {
		artists[songs[i].ArtistID]++
	}
	values := slices.Collect(maps.Keys(artists))
	return values[rand.IntN(len(values))]
}

func averagePlays(songs []*subsonic.Child) float64 {
	total := int64(0)
	for i := range songs {
		total += songs[i].PlayCount
	}
	return float64(total) / float64(len(songs))
}

func weight(selected []*subsonic.Child, averagePlayCount float64, lastSong, song *subsonic.Child) float64 {
	return ratingWeight(song) *
		frequencyWeight(averagePlayCount, song) *
		followerWeight(lastSong, song) *
		duplicateWeight(selected, song)
}

func ratingWeight(song *subsonic.Child) float64 {
	if !song.Starred.IsZero() {
		return *weightStarred
	} else if song.UserRating == 1 {
		return *weightRatedOne
	} else if song.UserRating == 2 {
		return *weightRatedTwo
	} else if song.UserRating == 3 {
		return *weightRatedThree
	} else if song.UserRating == 4 {
		return *weightRatedFour
	} else if song.UserRating == 5 {
		return *weightRatedFive
	} else {
		return 1
	}
}

func frequencyWeight(averagePlayCount float64, song *subsonic.Child) float64 {
	if song.PlayCount == 0 {
		return *weightNeverPlayed
	} else if float64(song.PlayCount) >= 2*averagePlayCount {
		return *weightFrequentlyPlayed
	} else if float64(song.PlayCount) <= 0.5*averagePlayCount {
		return *weightRarelyPlayed
	} else {
		return 1
	}
}

func followerWeight(lastSong, song *subsonic.Child) float64 {
	if lastSong != nil && lastSong.AlbumID == song.AlbumID {
		return *weightSameAlbum
	} else if lastSong != nil && lastSong.ArtistID == song.ArtistID {
		return *weightSameArtist
	} else if lastSong != nil && lastSong.Genre == song.Genre {
		return *weightSameGenre
	} else {
		return 1
	}
}

func duplicateWeight(selected []*subsonic.Child, song *subsonic.Child) float64 {
	for i := range selected {
		if selected[i].ID == song.ID {
			return *weightDuplicate
		}
	}
	return 1
}
