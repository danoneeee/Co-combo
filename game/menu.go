package game

import (
	"Cocombo/assets"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	resource "github.com/quasilyte/ebitengine-resource"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"time"
)

type Menu struct {
	Active       bool
	NameInput    string
	Font         font.Face
	Users        []User // Список пользователей
	SelectedUser int    // Индекс выбранного пользователя
	Background   *ebiten.Image
	Loader       *resource.Loader
	User         *User
	screenWidth  int // Ширина экрана
	screenHeight int // Высота экрана
}

func NewMenu(loader *resource.Loader, screenWidth, screenHeight int) *Menu {
	users, err := LoadUsers()
	if err != nil {
		log.Println("Ошибка при загрузке пользователей:", err)
		users = []User{} // Используем пустой список, если не удалось загрузить
	}

	return &Menu{
		Active:       true,
		Users:        users,
		SelectedUser: 0,
		Font:         loader.LoadFont(assets.Roboto).Face,
		Background:   loader.LoadImage(assets.ImageMenu).Data,
		Loader:       loader,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

func (m *Menu) Update() *User {
	if m.Active {
		// Обработка ввода текста
		for _, r := range ebiten.AppendInputChars(nil) {
			if r == '\b' && len(m.NameInput) > 0 { // Backspace
				m.NameInput = m.NameInput[:len(m.NameInput)-1]
			} else if r >= ' ' && r <= '~' || r >= 'а' && r <= 'я' || r >= 'А' && r <= 'Я' {
				m.NameInput += string(r)
			}
		}

		// Переключение между пользователями
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			m.SelectedUser = (m.SelectedUser + 1) % len(m.Users)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			m.SelectedUser = (m.SelectedUser - 1 + len(m.Users)) % len(m.Users)
		}

		// Выбор пользователя
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			m.Active = false
			if m.SelectedUser < len(m.Users) {
				return &m.Users[m.SelectedUser] // Возвращаем выбранного пользователя
			}
		}

		// Создание нового пользователя
		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			if len(m.NameInput) > 0 {
				newUser := User{
					ID:    generateUserID(), // Генерация уникального ID
					Name:  m.NameInput,
					Coins: 0,
				}
				m.Users = append(m.Users, newUser)
				m.NameInput = "" // Сбрасываем ввод имени
				if err := SaveUsers(m.Users); err != nil {
					log.Println("Ошибка при сохранении пользователей:", err)
				}
			}
		}
	}
	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) {
	if m.Active {
		// Масштабируем фон под размер экрана
		bgWidth, bgHeight := m.Background.Size()
		scaleX := float64(m.screenWidth) / float64(bgWidth)
		scaleY := float64(m.screenHeight) / float64(bgHeight)

		options := &ebiten.DrawImageOptions{}
		options.GeoM.Scale(scaleX, scaleY)
		screen.DrawImage(m.Background, options)

		// Отрисовка списка пользователей
		for i, user := range m.Users {
			msg := user.Name
			if i == m.SelectedUser {
				msg = "> " + msg // Выделяем выбранного пользователя
			}
			text.Draw(screen, msg, m.Font, 100, 100+i*30, color.White)
		}
		// Отрисовка ввода имени
		text.Draw(screen, "Введите имя: "+m.NameInput, m.Font, 100, m.screenHeight-100, color.White)
		// Подсказка
		text.Draw(screen, "Нажмите Enter для выбора, N для создания нового пользователя", m.Font, 100, m.screenHeight-50, color.White)
	}
}

// generateUserID генерирует уникальный ID для пользователя
func generateUserID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
