package game

import (
	"Cocombo/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	resource "github.com/quasilyte/ebitengine-resource"
	"golang.org/x/image/font"
	"image/color"
)

type Menu struct {
	Active       bool
	NameInput    string
	Font         font.Face
	Background   *ebiten.Image
	Loader       *resource.Loader
	User         *User
	screenWidth  int // Ширина экрана
	screenHeight int // Высота экрана
}

func NewMenu(loader *resource.Loader, user *User, screenWidth, screenHeight int) *Menu {

	return &Menu{
		Active:       true,
		NameInput:    "",
		Font:         loader.LoadFont(assets.Roboto).Face,
		Background:   loader.LoadImage(assets.ImageMenu).Data,
		Loader:       loader,
		User:         user,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

func (m *Menu) Update() {
	if m.Active {
		// Обработка ввода текста
		for _, r := range ebiten.AppendInputChars(nil) {
			if r == '\b' && len(m.NameInput) > 0 { // Backspace
				m.NameInput = m.NameInput[:len(m.NameInput)-1]
			} else if r >= ' ' && r <= '~' || r >= 'а' && r <= 'я' || r >= 'А' && r <= 'Я' {
				m.NameInput += string(r)
			}
		}

		// Завершение ввода по нажатию Enter
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			m.Active = false
		}
	}
}

func (m *Menu) Draw(screen *ebiten.Image) {
	if m.Active {
		// Масштабируем фон под размер экрана
		bgWidth, bgHeight := m.Background.Size()
		scaleX := float64(m.screenWidth) / float64(bgWidth)
		scaleY := float64(m.screenHeight) / float64(bgHeight)

		// Создаем options для отрисовки с масштабированием
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Scale(scaleX, scaleY)

		// Отрисовка фона
		screen.DrawImage(m.Background, options)

		// Отрисовка текста
		msg := m.NameInput
		text.Draw(screen, msg, m.Font, 570, 325, color.Black)
	}
}
