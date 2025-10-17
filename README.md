# BASS - Batch Audio Selection for Subsonic

BASS is a tool to automatically generate playlists for a Subsonic server,
such as [Navidrome](https://www.navidrome.org/).

It selects tracks randomly, but weighted according to their rating, play
history, and so on.

You configure BASS using flags or environment vars:

| Flag                     | Env var                  | Description                                                                                                                         |
|--------------------------|--------------------------|-------------------------------------------------------------------------------------------------------------------------------------|
| subsonic-server          | SUBSONIC_SERVER          | Base URL for the server, e.g. https://music.example.com/                                                                            |
| subsonic-user            | SUBSONIC_USER            | Username for the subsonic server, may be blank if auth is not required                                                              |
| subsonic-pass            | SUBSONIC_PASS            | Password for the subsonic server, may be blank if auth is not required                                                              |
| run-at                   | RUN_AT                   | Time to run playlist generation. If not specified, runs immediately. Must be a full time including timezone (e.g. `03:00:00+02:00`) |
| playlist-name            | PLAYLIST_NAME            | Name of the playlist to create/update (default: "Daily Mix")                                                                        |
| target-length            | TARGET_LENGTH            | Total length of tracks to add to the playlist (default: 18 hours)                                                                   |
| weight-starred           | WEIGHT_STARRED           | How to weight starred (loved) tracks (default: 20)                                                                                  |
| weight-rated-1           | WEIGHT_RATED_1           | How to weight tracks rated 1/5 (default: 0)                                                                                         |
| weight-rated-2           | WEIGHT_RATED_2           | How to weight tracks rated 2/5 (default: 0.5)                                                                                       |
| weight-rated-3           | WEIGHT_RATED_3           | How to weight tracks rated 3/5 (default: 1)                                                                                         |
| weight-rated-4           | WEIGHT_RATED_4           | How to weight tracks rated 4/5 (default: 5)                                                                                         |
| weight-rated-5           | WEIGHT_RATED_5           | How to weight tracks rated 5/5 (default: 10)                                                                                        |
| weight-never-played      | WEIGHT_NEVER_PLAYED      | How to weight tracks with a play count of zero (default: 5)                                                                         |
| weight-rarely-played     | WEIGHT_RARELY_PLAYED     | How to weight tracks with a play count under half the average (default: 3)                                                          |
| weight-frequently-played | WEIGHT_FREQUENTLY_PLAYED | How to weight tracks with a play count over double the average (default: 1)                                                         |
| weight-same-album        | WEIGHT_SAME_ALBUM        | How to weight tracks from the same album as the previous song (default: 20)                                                         |
| weight-same-artist       | WEIGHT_SAME_ARTIST       | How to weight tracks from the same artist as the previous song (default: 15)                                                        |
| weight-same-genre        | WEIGHT_SAME_GENRE        | How to weight tracks from the same genre as the previous song (default: 10)                                                         |
| weight-early-track       | WEIGHT_EARLY_TRACK       | How to weight tracks at positions 1-3 on an album (default: 10)                                                                     |
| weight-middle-track      | WEIGHT_MIDDLE_TRACK      | How to weight tracks at positions 4-6 on an album (default: 5)                                                                      |
| weight-late-track        | WEIGHT_LATE_TRACK        | How to weight tracks at positions 7-10 on an album (default: 1)                                                                     |
| weight-extended-track    | WEIGHT_EXTENDED_TRACK    | How to weight tracks at positions 11+ on an album (default: 0.75)                                                                   |
| weight-duplicate         | WEIGHT_DUPLICATE         | How to weight tracks that have already been picked (default: 0)                                                                     |
| log.level                | LOG_LEVEL                | Minimum level of logs which should be displayed (default: INFO; alternatives: WARN, ERROR, DEBUG)                                   |
| log.format               | LOG_FORMAT               | Format for logs (default: text; alternatives: json)                                                                                 |

## Weighting algorithm

Each track starts with a weight of `1.0`, which is then multiplied by five
components:

- The score from the track rating (starred, rated 1-5)
- The score from the play count (never played, rarely played, frequently played)
- The score from matching the last song (album, artist, genre)
- The score from being a duplicate
- The score from the track position on the album (early, middle, late, extended)

The first listed weight wins within each section (e.g., if a song is both
starred and rated 5, its "track rating" score will be the starred weight).

A weight of `1.0` makes no change, a weight of `0.0` excludes the track
entirely, and higher weights make it more likely to be picked.