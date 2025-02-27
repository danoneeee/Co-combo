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
	ImageMenu
)
const (
	Roboto resource.FontID = iota // Обычный шрифт
)

func registerImageResources(loader *resource.Loader) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageBackground: {Path: "_data/images/background.png"},
		ImageMouse1:     {Path: "_data/images/mouse1.png"},
		ImageMouse2:     {Path: "_data/images/mouse2.png"},
		ImageMouse3:     {Path: "_data/images/mouse3.png"},
		ImageMouse4:     {Path: "_data/images/mouse4.png"},
		ImageMouse5:     {Path: "_data/images/mouse5.png"},
		ImageKeyboard1:  {Path: "_data/images/keyboard1.png"},
		ImageKeyboard2:  {Path: "_data/images/keyboard2.png"},
		ImageKeyboard3:  {Path: "_data/images/keyboard3.png"},
		ImageKeyboard4:  {Path: "_data/images/keyboard4.png"},
		ImageKeyboard5:  {Path: "_data/images/keyboard5.png"},
		ImageMenu:       {Path: "_data/images/menu.png"},
	}

	for id, res := range imageResources {
		loader.ImageRegistry.Set(id, res)
	}
}

func RegisterResource(loader *resource.Loader) {
	// Загрузка шрифта
	loader.FontRegistry.Set(Roboto, resource.FontInfo{
		Path: "_data/images/fonts/Pixelizer/PixelizerBold.ttf", // Путь к шрифту
		Size: 16,                                               // Размер шрифта
	})
}
