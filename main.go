package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/csmith/envflag/v2"
	"github.com/csmith/slogflags"
	"github.com/supersonic-app/go-subsonic/subsonic"
)

var (
	runAt = flag.String("run-at", "", "Time at which to generate the playlist, instead of immediately.")
)

func main() {
	envflag.Parse()
	_ = slogflags.Logger(slogflags.WithSetDefault(true))

	c, err := connect()
	if err != nil {
		slog.Error("Failed to connect to server", "error", err)
		os.Exit(1)
	}

	if *runAt == "" {
		slog.Debug("No time specified, doing one-shot run")
		run(c)
		return
	}

	day := time.Now()

	for {
		target, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT%s", day.Format("2006-01-02"), *runAt))
		if err != nil {
			slog.Error("Failed to parse time", "error", err, "runAt", *runAt)
			os.Exit(1)
		}

		day = day.Add(24 * time.Hour)

		if target.After(time.Now()) {
			until := time.Until(target)
			slog.Info("Sleeping until next run", "duration", until, "target", target, "runAt", *runAt)
			time.Sleep(until)
			run(c)
		}
	}
}

func run(c *subsonic.Client) {
	s, err := songs(c)
	if err != nil {
		slog.Error("Failed to list songs", "error", err)
		os.Exit(1)
	}

	slog.Info("Found songs", "count", len(s))

	playlist := generate(s)
	slog.Info("Playlist generated", "count", len(playlist))

	if err := updatePlaylist(c, playlist); err != nil {
		slog.Error("Failed to update playlist", "error", err)
	}
}
