package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// request_url = f"http://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user={USERNAME}&from={to_timestamp}&to={from_timestamp}&api_key={API_KEY}&format=json&limit=200"

type LastFm struct {
	baseUrl  string
	username string
	apiKey   string
	client   http.Client
}

type LastFmTrackImage struct {
	Size string `json:"size"`
	Url  string `json:"#text"`
}

type LastFmTrack struct {
	Artist struct {
		Name string `json:"#text"`
	} `json:"artist"`
	Album struct {
		Name string `json:"#text"`
	} `json:"album"`
	Name  string             `json:"name"`
	Image []LastFmTrackImage `json:"image"`
}

type RecentTracksResponse struct {
	RecentTracks struct {
		Track []LastFmTrack `json:"track"`
	} `json:"recenttracks"`
}

type Image struct {
	Size string
	Url  string
}

type RecentTrack struct {
	Album  string
	Artist string
	Track  string
	Images []Image
}

func NewLastFm(config *Config) *LastFm {
	return &LastFm{
		baseUrl:  "http://ws.audioscrobbler.com/2.0",
		username: config.LastFmUsername,
		apiKey:   config.LastfmApiKey,
	}
}

func (l *LastFm) GetRecentTracks(startTimestamp, endTimestamp int64) ([]*RecentTrack, error) {

	url := fmt.Sprintf("%s?method=user.getrecenttracks&user=%s&from=%d&to=%d&api_key=%s&format=json&limit=200",
		l.baseUrl,
		l.username,
		endTimestamp,
		startTimestamp,
		l.apiKey,
	)
	slog.Debug(url)

	res, err := http.Get(url)
	if err != nil {
		slog.Error("error in last.fm api", "error", err)
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("error reading last.fm api response", "error", err)
		return nil, err
	}

	var response map[string]interface{}
	json.Unmarshal(b, &response)

	var tracks []interface{}
	tracks, ok := response["recenttracks"].(map[string]interface{})["track"].([]interface{})
	if !ok {
		track := response["recenttracks"].(map[string]interface{})["track"].(map[string]interface{})
		tracks = append(tracks, track)
	}

	var recentTracks []*RecentTrack
	for _, track := range tracks {
		album := track.(map[string]interface{})["album"].(map[string]interface{})["#text"].(string)
		artist := track.(map[string]interface{})["artist"].(map[string]interface{})["#text"].(string)
		tname := track.(map[string]interface{})["name"].(string)
		imgs := track.(map[string]interface{})["image"].([]interface{})

		var images []Image
		for _, img := range imgs {
			images = append(images, Image{
				Size: img.(map[string]interface{})["size"].(string),
				Url:  img.(map[string]interface{})["#text"].(string),
			})
		}

		recentTracks = append(recentTracks, &RecentTrack{
			Album:  album,
			Artist: artist,
			Track:  tname,
			Images: images,
		})

	}

	return recentTracks, nil
}
