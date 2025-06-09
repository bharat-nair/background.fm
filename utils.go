package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/user"
	"strings"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

type Config struct {
	LastfmApiKey       string `json:"lastfm_api_key"`
	LastfmSharedSecret string `json:"lastfm_shared_secret"`
	LastFmUsername     string `json:"lastfm_username"`
	DownloadDir        string `json:"download_dir"`
}

func GetConfigFile() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("%s/.config/background.fm/config.json", usr.HomeDir)
	slog.Info("reading configuration file at", "path", filePath)

	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetImageUrl(image []Image, size string) string {
	var url string

	for _, img := range image {
		if img.Size == size {
			url = img.Url
		}
	}

	return url
}

func DownloadImages(urls []string, downloadsDir string) []string {
	var filePaths []string

	for _, url := range urls {
		res, err := http.Get(url)
		if err != nil {
			slog.Error("error downloading image", "error", err)
			continue
		}
		defer res.Body.Close()

		s := strings.Split(url, "/")
		filePath := fmt.Sprintf("%s/%s", downloadsDir, s[len(s)-1])

		f, err := os.Create(filePath)
		if err != nil {
			slog.Error("error saving image", "error", err)
			continue
		}

		io.Copy(f, res.Body)
		filePaths = append(filePaths, filePath)
	}

	return filePaths
}

func GetColorPalette(filePath string, n int) ([][]int, error) {
	var palette [][]int

	f, err := os.Open(filePath)
	if err != nil {
		slog.Error("error reading image", "error", err)
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		slog.Error("error decoding image", "error", err)
		return nil, err
	}

	bounds := img.Bounds()
	var imgc [][]uint32

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			imgc = append(imgc, []uint32{r, g, b, a})
		}
	}

	km := kmeans.New()
	var dataset clusters.Observations
	for _, rgb := range imgc {
		dataset = append(dataset, clusters.Coordinates{
			float64(rgb[0]),
			float64(rgb[1]),
			float64(rgb[2]),
		})
	}
	clusters, err := km.Partition(dataset, n)
	if err != nil {
		slog.Error("error getting color palette", "error", err)
		return nil, err
	}

	for _, c := range clusters {
		palette = append(palette, []int{int(c.Center[0] / 256), int(c.Center[1] / 256), int(c.Center[2] / 256)})
	}

	return palette, nil
}
