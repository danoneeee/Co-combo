package main

import (
	"Cocombo/assets"
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	resource "github.com/quasilyte/ebitengine-resource"
	"image/color"
	"log"
	"os"
)

const (
	fieldWidth   = 450
	fieldHeight  = 450
	screenWidth  = 800
	screenHeight = 600
	gridSize     = 8
	cellSize     = 50
)

type SaveData struct {
	Images []Item     `json:"images"` // Данные о предметах
	Grid   []GridSave `json:"grid"`   // Данные о клетках
}

type GridSave struct {
	ID       int  `json:"id"`       // Уникальный ID клетки
	CenterX  int  `json:"centerX"`  // Координата X центра клетки
	CenterY  int  `json:"centerY"`  // Координата Y центра клетки
	Occupied bool `json:"occupied"` // Занята ли клетка
	ItemID   int  `json:"itemID"`   // ID предмета, лежащего в клетке (если есть)
}
type Grid struct {
	id       int   // Уникальный ID клетки
	centerX  int   // Координата X центра клетки
	centerY  int   // Координата Y центра клетки
	occupied bool  // Занята ли клетка
	item     *Item // Предмет, лежащий в клетке (если есть)
}

type Item struct {
	X, Y             int     // Позиция изображения
	Dragging         bool    // Флаг, указывающий, что изображение перемещается
	OffsetX, OffsetY float64 // Смещение курсора относительно изображения
	TypeObject       resource.ImageID
}

type Game struct {
	cursorX, cursorY  int
	leftButtonPressed bool
	loader            *resource.Loader
	images            []*Item // Данные о нескольких изображениях
	draggingIndex     int     // Индекс изображения, которое перемещается
	grid              []Grid  // Сетка
}

func createGrid() []Grid {
	grid := make([]Grid, 0, gridSize*gridSize)
	for i := 1; i < gridSize+1; i++ {
		for j := 1; j < gridSize+1; j++ {
			centerX := 50 + j*cellSize - cellSize/2
			centerY := 50 + i*cellSize - cellSize/2
			grid = append(grid, Grid{
				id:       i*gridSize + j + 1, // Уникальный ID
				centerX:  centerX,
				centerY:  centerY,
				occupied: false,
				item:     nil,
			})
		}
	}
	fmt.Println(grid)
	return grid
}

func (g *Game) saveGame(filename string) error {
	// Создаем структуру для сохранения
	saveData := SaveData{
		Images: make([]Item, len(g.images)),
		Grid:   make([]GridSave, len(g.grid)),
	}

	// Заполняем данные о изображениях
	for i, img := range g.images {
		saveData.Images[i] = Item{
			X:          img.X,
			Y:          img.Y,
			Dragging:   img.Dragging,
			OffsetX:    img.OffsetX,
			OffsetY:    img.OffsetY,
			TypeObject: img.TypeObject,
		}
	}

	// Заполняем данные о клетках
	for i, cell := range g.grid {
		saveData.Grid[i] = GridSave{
			ID:       cell.id,
			CenterX:  cell.centerX,
			CenterY:  cell.centerY,
			Occupied: cell.occupied,
			ItemID:   -1, // По умолчанию предмета нет
		}
		if cell.item != nil {
			// Находим ID предмета
			for j, img := range g.images {
				if cell.item == img {
					saveData.Grid[i].ItemID = j
					break
				}
			}
		}
	}

	// Сериализуем данные в JSON
	data, err := json.Marshal(saveData)
	if err != nil {
		return err
	}

	// Сохраняем JSON в файл
	return os.WriteFile(filename, data, 0644)
}

func (g *Game) loadGame(filename string) error {
	// Читаем данные из файла
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Десериализуем JSON
	var saveData SaveData
	if err := json.Unmarshal(data, &saveData); err != nil {
		return err
	}

	// Восстанавливаем предметы
	g.images = make([]*Item, len(saveData.Images))
	for i, imgSave := range saveData.Images {
		g.images[i] = &Item{
			X:          imgSave.X,
			Y:          imgSave.Y,
			Dragging:   imgSave.Dragging,
			OffsetX:    imgSave.OffsetX,
			OffsetY:    imgSave.OffsetY,
			TypeObject: imgSave.TypeObject,
		}
	}

	// Восстанавливаем клетки
	g.grid = make([]Grid, len(saveData.Grid))
	for i, cellSave := range saveData.Grid {
		g.grid[i] = Grid{
			id:       cellSave.ID,
			centerX:  cellSave.CenterX,
			centerY:  cellSave.CenterY,
			occupied: cellSave.Occupied,
			item:     nil, // По умолчанию предмета нет
		}
		if cellSave.ItemID != -1 {
			// Привязываем предмет к клетке
			g.grid[i].item = g.images[cellSave.ItemID]
		}
	}

	return nil
}

func (g *Game) Update() error {
	// Получаем текущие координаты курсора
	g.cursorX, g.cursorY = ebiten.CursorPosition()

	// Проверяем, нажата ли левая кнопка мыши
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.draggingIndex == -1 {
			// Начало перемещения: ищем изображение под курсором
			for i := range g.images {
				img := g.images[i]
				// Проверяем, находится ли курсор над изображением
				if g.cursorX >= img.X && g.cursorX <= img.X+cellSize &&
					g.cursorY >= img.Y && g.cursorY <= img.Y+cellSize {
					// Запоминаем смещение курсора относительно изображения
					img.OffsetX = float64(g.cursorX - img.X)
					img.OffsetY = float64(g.cursorY - img.Y)
					img.Dragging = true
					g.draggingIndex = i // Запоминаем индекс перемещаемого изображения

					// Освобождаем клетку, если предмет был привязан к ней
					for j := range g.grid {
						if g.grid[j].item == g.images[i] {
							g.grid[j].occupied = false
							g.grid[j].item = nil
							break
						}
					}
					break
				}
			}
		} else {
			// Перемещение изображения вместе с курсором с учетом смещения
			img := g.images[g.draggingIndex]
			img.X = int(float64(g.cursorX) - img.OffsetX)
			img.Y = int(float64(g.cursorY) - img.OffsetY)
		}
	} else {
		// Кнопка мыши отпущена: привязываем изображение к сетке
		if g.draggingIndex != -1 {
			img := g.images[g.draggingIndex]

			// Ищем ближайшую клетку
			for j := range g.grid {
				cell := &g.grid[j]
				// Проверяем, находится ли предмет в пределах клетки
				if g.cursorX >= cell.centerX-cellSize/2 && g.cursorX <= cell.centerX+cellSize/2 &&
					g.cursorY >= cell.centerY-cellSize/2 && g.cursorY <= cell.centerY+cellSize/2 {
					fmt.Println(cell.centerX-cellSize/2, cell.centerY-cellSize/2)
					// Если клетка свободна, привязываем предмет к ней
					if !cell.occupied {
						img.X = cell.centerX - cellSize/2
						img.Y = cell.centerY - cellSize/2
						cell.occupied = true
						cell.item = img
						fmt.Println(cell.occupied, cell.item)
					} else {
						if cell.item.TypeObject == img.TypeObject {

							cell.item.TypeObject += 1
							g.images = append(g.images[:g.draggingIndex], g.images[g.draggingIndex+1:]...)

							// Сбрасываем индекс перемещаемого объекта
							g.draggingIndex = -1
							break
						}
					}
					break
				}
			}
			img.Dragging = false
			g.draggingIndex = -1
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Заливаем фон белым цветом
	screen.Fill(color.RGBA{R: 130, G: 155, B: 175, A: 255})

	// Рисуем сетку
	for i := 1; i <= gridSize+1; i++ {
		// Вертикальные линии
		ebitenutil.DrawLine(screen, float64(i*cellSize), 50, float64(i*cellSize), fieldHeight, color.Black)
		// Горизонтальные линии
		ebitenutil.DrawLine(screen, 50, float64(i*cellSize), fieldWidth, float64(i*cellSize), color.Black)
	}

	// Рисуем все изображения
	for _, img := range g.images {
		// Загружаем изображение по его ID
		image := g.loader.LoadImage(img.TypeObject)

		// Создаем объект DrawImageOptions для позиционирования изображения
		var options ebiten.DrawImageOptions
		options.GeoM.Translate(float64(img.X), float64(img.Y))

		// Рисуем изображение
		screen.DrawImage(image.Data, &options)
	}

	// Отображаем позицию курсора и состояние кнопки
	debugText := fmt.Sprintf(
		"Cursor X: %d\nCursor Y: %d\nLeft Button Pressed: %v\n",
		g.cursorX, g.cursorY, g.leftButtonPressed)
	ebitenutil.DebugPrint(screen, debugText)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func createLoader() *resource.Loader {
	sampleRate := 44100
	audioContext := audio.NewContext(sampleRate)
	loader := resource.NewLoader(audioContext)
	loader.OpenAssetFunc = assets.OpenAsset
	return loader
}

func main() {
	loader := createLoader()
	assets.RegisterResources(loader)

	g := &Game{
		loader:        loader,
		images:        []*Item{},
		draggingIndex: -1,
		grid:          createGrid(),
	}

	// Загружаем сохраненные данные (если файл существует)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Merge-2 Game")

	if err := g.loadGame("save.json"); err != nil {
		if os.IsNotExist(err) {
			log.Println("Файл сохранения не найден, создаем новую игру")
		} else {
			log.Println("Ошибка при загрузке игры:", err)
		}
		// Инициализируем начальные данные
		g.images = []*Item{
			{X: 100, Y: 100, TypeObject: 1},
			{X: 100, Y: 150, TypeObject: 1},
			{X: 200, Y: 200, TypeObject: 2},
			{X: 100, Y: 300, TypeObject: 3},
			{X: 200, Y: 250, TypeObject: 4},
			{X: 200, Y: 300, TypeObject: 4},
		}
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

	// Сохраняем данные при завершении программы
	if err := g.saveGame("save.json"); err != nil {
		log.Println("Ошибка при сохранении игры:", err)
	}
}
