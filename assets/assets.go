package assets

import (
	"embed"
	resource "github.com/quasilyte/ebitengine-resource"
	"io"
)

//go:embed _data/images/*
var gameAssets embed.FS

func OpenAsset(path string) io.ReadCloser {
	f, err := gameAssets.Open(path) // Указываем полный путь
	if err != nil {
		panic("cant open asset: " + err.Error()) // Добавляем информацию об ошибке
	}
	return f
}
func RegisterResources(loader *resource.Loader) {
	registerImageResources(loader)
	RegisterResource(loader)
}
