package game

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	gridSize    = 8
	cellSize    = 50
	fieldWidth  = 450
	fieldHeight = 450
)

type Grid struct {
	id       int   // Уникальный ID клетки
	centerX  int   // Координата X центра клетки
	centerY  int   // Координата Y центра клетки
	occupied bool  // Занята ли клетка
	item     *Item // Предмет, лежащий в клетке (если есть)
}

type GridSave struct {
	ID       int  `json:"id"`       // Уникальный ID клетки
	CenterX  int  `json:"centerX"`  // Координата X центра клетки
	CenterY  int  `json:"centerY"`  // Координата Y центра клетки
	Occupied bool `json:"occupied"` // Занята ли клетка
	ItemID   int  `json:"itemID"`   // ID предмета, лежащего в клетке (если есть)
}
type SaveData struct {
	Images []Item     `json:"images"` // Данные о предметах
	Grid   []GridSave `json:"grid"`   // Данные о клетках
}

func CreateGrid() []Grid {
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

func (g *Game) UpdateGridOccupancy() {
	// Сначала сбрасываем все клетки в состояние "не занято"
	for i := range g.Grid {
		g.Grid[i].occupied = false
		g.Grid[i].item = nil
	}

	// Затем обновляем клетки, на которых находятся объекты
	for _, img := range g.Images {
		for i := range g.Grid {
			cell := &g.Grid[i]
			if img.X >= cell.centerX && img.X <= cell.centerX+cellSize/2 &&
				img.Y >= cell.centerY && img.Y <= cell.centerY+cellSize/2 {
				cell.occupied = true
				cell.item = img
				break
			}
		}
	}
}

func (g *Game) SaveGame(filename string) error {
	// Создаем структуру для сохранения
	saveData := SaveData{
		Images: make([]Item, len(g.Images)),
		Grid:   make([]GridSave, len(g.Grid)),
	}

	// Заполняем данные о изображениях
	for i, img := range g.Images {
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
	for i, cell := range g.Grid {
		saveData.Grid[i] = GridSave{
			ID:       cell.id,
			CenterX:  cell.centerX,
			CenterY:  cell.centerY,
			Occupied: cell.occupied,
			ItemID:   -1, // По умолчанию предмета нет
		}
		if cell.item != nil {
			// Находим ID предмета
			for j, img := range g.Images {
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

func (g *Game) LoadGame(filename string) error {
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
	g.Images = make([]*Item, len(saveData.Images))
	for i, imgSave := range saveData.Images {
		g.Images[i] = &Item{
			X:          imgSave.X,
			Y:          imgSave.Y,
			Dragging:   imgSave.Dragging,
			OffsetX:    imgSave.OffsetX,
			OffsetY:    imgSave.OffsetY,
			TypeObject: imgSave.TypeObject,
		}
	}

	// Восстанавливаем клетки
	g.Grid = make([]Grid, len(saveData.Grid))
	for i, cellSave := range saveData.Grid {
		g.Grid[i] = Grid{
			id:       cellSave.ID,
			centerX:  cellSave.CenterX,
			centerY:  cellSave.CenterY,
			occupied: cellSave.Occupied,
			item:     nil, // По умолчанию предмета нет
		}
		if cellSave.ItemID != -1 {
			// Привязываем предмет к клетке
			g.Grid[i].item = g.Images[cellSave.ItemID]
		}
	}

	return nil
}
