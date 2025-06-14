package main

import (
	"fmt"
	"os/exec"
)

func SetWallpaperKDE(filepath string, background string) (string, error) {
	cmd := exec.Command(
		"qdbus6",
		"org.kde.plasmashell",
		"/PlasmaShell",
		"org.kde.PlasmaShell.evaluateScript",
		fmt.Sprintf(`
		var allDesktops = desktops();
		// print(allDesktops);
		for (i=0;i<allDesktops.length;i++) {
			d = allDesktops[i];
			d.wallpaperPlugin = "org.kde.image";
			d.currentConfigGroup = Array("Wallpaper",
										"org.kde.image",
										"General");
			d.writeConfig("Image", "file://%s");
			// d.writeConfig("FillMode", "1");  // Centered
			d.writeConfig("Color", "%s");  // Set background color
			d.writeConfig("Background", "1");  // 0 = Blur, 1 = Solid Color
		}`, filepath, background),
	)

	stdout, err := cmd.Output()

	return string(stdout), err
}

func SetWallpaperSway(filepath string, background string) (string, error) {
	cmd := exec.Command(
		"swaybg",
		"--image",
		filepath,
		"--color",
		background,
		"--mode",
		"center",
	)

	stdout, err := cmd.Output()

	return string(stdout), err
}

func SetWallpaperGnome(filepath string, background string) (string, error) {

	cmd := exec.Command(
		"gsettings",
		"set",
		"org.gnome.desktop.background",
		"picture-uri",
		fmt.Sprintf("file://%s", filepath),
	)
	stdout, err := cmd.Output()
	if err != nil {
		return string(stdout), err
	}

	cmd = exec.Command(
		"gsettings",
		"set",
		"org.gnome.desktop.background",
		"picture-uri-dark",
		fmt.Sprintf("file://%s", filepath),
	)
	stdout, err = cmd.Output()
	if err != nil {
		return string(stdout), err
	}

	cmd = exec.Command(
		"gsettings",
		"set",
		"org.gnome.desktop.background",
		"picture-options",
		"centered",
	)
	stdout, err = cmd.Output()
	if err != nil {
		return string(stdout), err
	}

	cmd = exec.Command(
		"gsettings",
		"set",
		"org.gnome.desktop.background",
		"primary-color",
		background,
	)
	stdout, err = cmd.Output()
	if err != nil {
		return string(stdout), err
	}

	return string(stdout), err
}
