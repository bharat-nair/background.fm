package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"time"
)

func main() {

	dimensions := map[string]int{"small": 34, "medium": 64, "large": 174, "extralarge": 300}
	imageSize := flag.String("image_size", "extralarge", "The size of album art to use to build the wallpaper. Available options: extralarge, large, medium, small")
	wallpaperWidth := flag.Int("width", 1920, "The width of the wallpaper.")
	wallpaperHeight := flag.Int("height", 1080, "The height of the wallpaper.")
	walI := flag.Int("wal_i", 0, "The index to choose from the array of album arts. This art will then be used to extract colors for applying as the final wallpaper background.")
	desktopEnvironment := flag.String("desktop_environment", "", "The desktop environment to set the wallpaper in. Supported environments: kde")

	flag.Parse()

	config, err := GetConfigFile()
	if err != nil {
		log.Fatalln(err)
	}

	slog.Debug("read", "config", config)

	wallpaperDimensions := []int{*wallpaperWidth, *wallpaperHeight}
	numberOfImages := wallpaperDimensions[0] * wallpaperDimensions[1] / (dimensions[*imageSize] * dimensions[*imageSize])
	downloadsDir := config.DownloadDir
	oneDay := int64(24 * 60 * 60)

	var imageUrls []string
	urlSet := map[string]struct{}{}

	remaining := numberOfImages
	fromTimestamp := time.Now().Unix()
	toTimestamp := fromTimestamp - oneDay

	lastfm := NewLastFm(config)
	for remaining > 0 {
		tracks, err := lastfm.GetRecentTracks(fromTimestamp, toTimestamp)
		if err != nil {
			log.Fatalln(err)
		}

		for _, track := range tracks {
			url := GetImageUrl(track.Images, *imageSize)
			if _, exists := urlSet[url]; !exists {
				imageUrls = append(imageUrls, url)
				urlSet[url] = struct{}{}
			}

			remaining = numberOfImages - len(imageUrls)
			if remaining <= 0 {
				break
			}
		}

		if remaining == numberOfImages-len(imageUrls) {
			fromTimestamp = toTimestamp
		}

		toTimestamp = toTimestamp - 24*60*60
	}

	// download all images
	filePaths := DownloadImages(imageUrls, downloadsDir)

	// start creating the wallpaper
	// set a background using the color palette of one of the album arts
	wallpaperImage := image.NewRGBA(image.Rect(0, 0, wallpaperDimensions[0], wallpaperDimensions[1]))

	palette, err := GetColorPalette(filePaths[*walI], 4)
	if err != nil {
		palette = [][]int{}
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	background := palette[r.Intn(len(palette))]

	draw.Draw(
		wallpaperImage,
		wallpaperImage.Bounds(),
		image.NewUniform(color.RGBA{uint8(background[0]), uint8(background[1]), uint8(background[2]), 255}),
		image.Point{},
		draw.Src,
	)

	// add all images to the wallpaperImage with offsets, so final image is centered
	var currentWidth, currentHeight int
	var offsetWidth, offsetHeight int

	offsetWidth = wallpaperImage.Bounds().Max.X % dimensions[*imageSize] / 2
	offsetHeight = wallpaperImage.Bounds().Max.Y % dimensions[*imageSize] / 2

	currentWidth, currentHeight = offsetWidth, offsetHeight

	width := wallpaperImage.Bounds().Max.X - currentWidth
	height := wallpaperImage.Bounds().Max.Y - currentHeight

	for _, filePath := range filePaths {

		f, err := os.Open(filePath)
		if err != nil {
			slog.Error("error reading file", "error", err)
			continue
		}

		img, _, err := image.Decode(f)
		if err != nil {
			slog.Warn("error decoding image", "filePath", filePath, "error", err)
			continue
		}
		f.Close()

		if currentWidth >= width && currentHeight < height {
			currentWidth = offsetWidth
			currentHeight += img.Bounds().Max.X
		}
		if currentHeight >= height {
			break
		}

		draw.Draw(
			wallpaperImage,
			image.Rect(currentWidth, currentHeight, currentWidth+dimensions[*imageSize], currentHeight+dimensions[*imageSize]),
			img,
			image.Point{},
			draw.Over,
		)

		currentWidth += img.Bounds().Max.Y
	}

	outputFile := fmt.Sprintf("%s/%d.png", downloadsDir, time.Now().UnixNano())
	output, err := os.Create(outputFile)
	if err != nil {
		slog.Error("error creating wallpaper", "error", err)
	}
	defer output.Close()

	jpeg.Encode(output, wallpaperImage, nil)

	switch *desktopEnvironment {
	case "kde":
		SetWallpaperKDE(
			outputFile,
			fmt.Sprintf("#%02x%02x%02x", uint8(background[0]), uint8(background[1]), uint8(background[2])),
		)
	case "sway":
		SetWallpaperSway(
			outputFile,
			fmt.Sprintf("#%02x%02x%02x", uint8(background[0]), uint8(background[1]), uint8(background[2])),
		)
	case "gnome":
		SetWallpaperGnome(
			outputFile,
			fmt.Sprintf("#%02x%02x%02x", uint8(background[0]), uint8(background[1]), uint8(background[2])),
		)
	}
}
