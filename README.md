# background.fm
Set your desktop background as your now playing. Uses the Last.fm API.

# How Does It Work?
Last.fm tracks the music you listen to via a handy feature called 'scrobbling'.
Many [services](https://www.last.fm/about/trackmymusic) are supported.
Your scrobbled tracks can then be interacted in various ways via their [API](https://www.last.fm/api).
This app fetches the album art of your most recently played tracks, creates an image of specified size and sets it as the wallpaper.
Currently, setting wallpaper functionality is only supported in Linux; specifically on KDE, GNOME, and Sway, but adding your own desktop environment/window manager is trivial by adding a custom function in `wallpaper.go`.
The image generation will work regardless of the platform, and is tested on Windows 10.

# Install Instructions (Linux)
**Arch Linux** users can directly install from the [AUR](https://aur.archlinux.org/packages/background.fm)

# Build Instructions (Cross-platform)
1. Clone the repository
2. Install Go using your preferred package manager
3. Depending on your target platform, build the app:

   ```
   GOOS=linux GOARCH=amd64 go build -o background.fm main.go lastfm.go wallpaper.go utils.go # linux x64
   ```
   ```
   GOOS=windows GOARCH=amd64 go build -o background.fm main.go lastfm.go wallpaper.go utils.go # windows x64
   ```
   ```
   GOOS=darwin GOARCH=arm64 go build -o background.fm main.go lastfm.go wallpaper.go utils.go # macos arm
   ```
   ```
   GOOS=linux GOARCH=arm GOARM=7 go build -o background.fm main.go lastfm.go wallpaper.go utils.go # raspberrypi
   ```
   
5. Create configuration file at `$HOME/.config/background.fm/config.json`
6. Set wallpapers with the command: `./background.fm --width=1920 --height=1080 --desktop_environment=kde --image_size=extralarge`

# Sample Configuration File
```
{
    "lastfm_api_key": "<last.fm api key>",
    "lastfm_shared_secret": "<last.fm shared secret>",
    "lastfm_username": "<last.fm username>",
    "download_dir": "/tmp"
}
```

# Screenshots
`./background.fm --width=1920 --height=1080 --desktop_environment=kde --image_size=extralarge`

![Screenshot_20250609_232214](https://github.com/user-attachments/assets/3f3c0eda-690f-4e5f-ba19-f341585a30a4)


`./background.fm --width=640 --height=480 --desktop_environment=kde --image_size=extralarge --wal_i=1`

![Screenshot_20250609_232606](https://github.com/user-attachments/assets/70ea8d92-c693-439b-b351-5b35d40be5f8)


`./background.fm --width=1280 --height=720 --desktop_environment=kde --image_size=extralarge --wal_i=2`

![Screenshot_20250609_233519](https://github.com/user-attachments/assets/9f1b4884-ab68-48b8-a2d8-7cb96dc11fac)


`./background.fm --width=300 --height=300 --desktop_environment=kde --image_size=medium --wal_i=4`

![Screenshot_20250609_235912](https://github.com/user-attachments/assets/3295af55-14e7-48b6-ad20-173429a847e2)

