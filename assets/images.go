package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	_ "image/png"
)

const (
	ImageBackground resource.ImageID = iota
	ImageMouse1
	ImageMouse2
	ImageMouse3
	ImageMouse4
	ImageMouse5
	ImageKeyboard1
	ImageKeyboard2
	ImageKeyboard3
	ImageKeyboard4
	ImageKeyboard5
)

func registerImageResources(loader *resource.Loader) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageBackground: {Path: "background.png"},
		ImageMouse1:     {Path: "mouse1.png"},
		ImageMouse2:     {Path: "mouse2.png"},
		ImageMouse3:     {Path: "mouse3.png"},
		ImageMouse4:     {Path: "mouse4.png"},
		ImageMouse5:     {Path: "mouse5.png"},
		ImageKeyboard1:  {Path: "keyboard1.png"},
		ImageKeyboard2:  {Path: "keyboard2.png"},
		ImageKeyboard3:  {Path: "keyboard3.png"},
		ImageKeyboard4:  {Path: "keyboard4.png"},
		ImageKeyboard5:  {Path: "keyboard5.png"},
	}

	for id, res := range imageResources {
		loader.ImageRegistry.Set(id, res)
	}
}
