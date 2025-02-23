package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	_ "image/png"
)

const (
	ImageNone      resource.ImageID = iota // ID = 0 (не используется)
	ImageMouse1                            // ID = 1
	ImageMouse2                            // ID = 2
	ImageKeyboard1                         // ID = 3
	ImageKeyboard2                         // ID = 4
)

func registerImageResources(loader *resource.Loader) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageMouse1:    {Path: "mouse1.png"},
		ImageMouse2:    {Path: "mouse2.png"},
		ImageKeyboard1: {Path: "keyboard1.png"},
		ImageKeyboard2: {Path: "keyboard2.png"}, // Указываем только имя файла
	}

	for id, res := range imageResources {
		loader.ImageRegistry.Set(id, res)
	}
}
