package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	resource "github.com/quasilyte/ebitengine-resource"
	"image/color"
)

type Game struct {
	cursorX, cursorY  int
	leftButtonPressed bool
	Loader            *resource.Loader
	Images            []*Item
	DraggingIndex     int
	Grid              []Grid
	Menu              *Menu
	User              *User
	BackgroundImage   *ebiten.Image // Фоновое изображение
}

func (g *Game) Update() error {
	if g.Menu.Active {
		g.Menu.Update()
		return nil
	}
	// Логика обновления игры
	g.cursorX, g.cursorY = ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.DraggingIndex == -1 {
			for i := range g.Images {
				img := g.Images[i]
				if g.cursorX >= img.X && g.cursorX <= img.X+cellSize &&
					g.cursorY >= img.Y && g.cursorY <= img.Y+cellSize {
					img.OffsetX = float64(g.cursorX - img.X)
					img.OffsetY = float64(g.cursorY - img.Y)
					img.Dragging = true
					g.DraggingIndex = i

					for j := range g.Grid {
						if g.Grid[j].item == g.Images[i] {
							g.Grid[j].occupied = false
							g.Grid[j].item = nil
							break
						}
					}
					break
				}
			}
		} else {
			img := g.Images[g.DraggingIndex]
			img.X = int(float64(g.cursorX) - img.OffsetX)
			img.Y = int(float64(g.cursorY) - img.OffsetY)
		}
	} else {
		if g.DraggingIndex != -1 {
			img := g.Images[g.DraggingIndex]

			for j := range g.Grid {
				cell := &g.Grid[j]
				if g.cursorX >= cell.centerX-cellSize/2 && g.cursorX <= cell.centerX+cellSize/2 &&
					g.cursorY >= cell.centerY-cellSize/2 && g.cursorY <= cell.centerY+cellSize/2 {
					if !cell.occupied {
						img.X = cell.centerX - cellSize/2
						img.Y = cell.centerY - cellSize/2
						cell.occupied = true
						cell.item = img
					} else {
						if cell.item.TypeObject == img.TypeObject {
							if cell.item.TypeObject%5 != 0 {
								cell.item.TypeObject += 1
								g.Images = append(g.Images[:g.DraggingIndex], g.Images[g.DraggingIndex+1:]...)
								g.DraggingIndex = -1
								break
							} else {
								for k := 0; k < len(g.Grid); k++ {
									if !g.Grid[k].occupied {
										img.X = g.Grid[k].centerX - cellSize/2
										img.Y = g.Grid[k].centerY - cellSize/2
										g.Grid[k].occupied = true
										g.Grid[k].item = img
										break
									}
								}
							}
						} else {
							for k := 0; k < len(g.Grid); k++ {
								if !g.Grid[k].occupied {
									img.X = g.Grid[k].centerX - cellSize/2
									img.Y = g.Grid[k].centerY - cellSize/2
									g.Grid[k].occupied = true
									g.Grid[k].item = img
									break
								}
							}
						}
					}
					break
				}
			}
			img.Dragging = false
			g.DraggingIndex = -1
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.Menu.Active {
		g.Menu.Draw(screen)
		return
	}
	// screen.Fill(color.RGBA{R: 130, G: 155, B: 175, A: 255})
	if g.BackgroundImage != nil {
		screen.DrawImage(g.BackgroundImage, &ebiten.DrawImageOptions{})
	}

	for i := 1; i <= gridSize+1; i++ {
		ebitenutil.DrawLine(screen, float64(i*cellSize), 50, float64(i*cellSize), fieldHeight, color.Black)
		ebitenutil.DrawLine(screen, 50, float64(i*cellSize), fieldWidth, float64(i*cellSize), color.Black)
	}

	for _, img := range g.Images {
		image := g.Loader.LoadImage(img.TypeObject)
		var options ebiten.DrawImageOptions
		options.GeoM.Translate(float64(img.X), float64(img.Y))
		screen.DrawImage(image.Data, &options)
	}

	debugText := fmt.Sprintf(
		"Cursor X: %d\nCursor Y: %d\n",
		g.cursorX, g.cursorY)
	ebitenutil.DebugPrint(screen, debugText)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
