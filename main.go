package main

import (
	"Cocombo/assets"
	"Cocombo/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	//"github.com/hajimehoshi/ebiten/v2/inpututil"
	//"github.com/hajimehoshi/ebiten/v2/text"
	resource "github.com/quasilyte/ebitengine-resource"
	//"golang.org/x/image/font"
	//"golang.org/x/image/font/opentype"
	//"image/color"
	//"io"
	"log"
	"os"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var (
	user       *game.User
	menuActive bool = true
	nameInput  string
)

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

	menu := game.NewMenu(loader, screenWidth, screenHeight)
	log.Println(menu)
	// Основной цикл меню
	var currentUser *game.User
	currentUser = menu.Update()

	// После выхода из цикла меню
	if currentUser != nil {
		log.Println("Выбран пользователь:", currentUser.Name)
		// Логика работы с выбранным пользователем
	} else {
		log.Println("Пользователь не выбран.")
	}

	gameData, err := game.Load(currentUser.ID)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Сохранение не найдено, создаем новую игру")
			gameData = &game.GameD{
				User:   *currentUser,
				Images: []*game.Item{},
				Grid:   game.CreateGrid(),
			}
		} else {
			log.Fatal("Ошибка при загрузке игры:", err)
		}
	}
	menuActive = true
	// Загружаем игру для выбранного пользователя
	//if currentUser == nil {
	//	log.Println("gege")
	//}

	g := &game.Game{
		Loader:          loader,
		Data:            gameData,
		BackgroundImage: loader.LoadImage(assets.ImageBackground).Data,
		Menu:            menu,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Merge-2 Game")

	// Запускаем игру
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

	// Сохраняем данные при завершении программы
	if err := game.SaveG(currentUser.ID, g.Data); err != nil {
		log.Println("Ошибка при сохранении игры:", err)
	}
}

// App объединяет игру и меню
//type App struct {
//	*game.Game
//}
//
//func (a *App) Update() error {
//	if menuActive {
//		// Обработка ввода имени в меню
//		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
//			menuActive = false
//			user = &game.User{Name: nameInput, Coins: 0}
//			if err := game.SaveUser(user); err != nil {
//				log.Println("Ошибка при сохранении пользователя:", err)
//			}
//			a.Game.User = user
//		} else {
//			// Обработка ввода текста
//			for _, r := range ebiten.AppendInputChars(nil) {
//				if r == '\b' && len(nameInput) > 0 {
//					// Удаление последнего символа (backspace)
//					nameInput = nameInput[:len(nameInput)-1]
//				} else if r >= ' ' {
//					// Добавление символа, если это не управляющий символ
//					nameInput += string(r)
//				}
//			}
//		}
//		return nil
//	}
//	return a.Game.Update()
//}
//
//func (a *App) Draw(screen *ebiten.Image) {
//	if menuActive {
//		// Отрисовка меню
//		msg := "Введите ваше имя: " + nameInput
//		text.Draw(screen, msg, loadFont(), 100, 100, color.White)
//		return
//	}
//	a.Game.Draw(screen)
//}
//
//func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
//	return outsideWidth, outsideHeight
//}
