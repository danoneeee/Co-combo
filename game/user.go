package game

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// User представляет данные пользователя
type User struct {
	ID    string `json:"id"`    // Уникальный идентификатор пользователя
	Name  string `json:"name"`  // Имя пользователя
	Coins int    `json:"coins"` // Количество монет
}

// GameData представляет данные игры для конкретного пользователя
type GameD struct {
	User   User    `json:"user"`   // Данные пользователя
	Images []*Item `json:"images"` // Объекты на поле
	Grid   []Grid  `json:"grid"`   // Состояние сетки
}

// SaveGame сохраняет данные игры в файл
func SaveG(userID string, gameData *GameD) error {
	file, err := json.Marshal(gameData)
	if err != nil {
		return err
	}
	return os.WriteFile("save_"+userID+".json", file, 0644)
}

func Load(userID string) (*GameD, error) {
	file, err := os.ReadFile("save_" + userID + ".json")
	if err != nil {
		return nil, err
	}
	var gameData GameD
	if err := json.Unmarshal(file, &gameData); err != nil {
		return nil, err
	}
	return &gameData, nil
}

// LoadUsers загружает список пользователей из файла
func LoadUsers() ([]User, error) {
	file, err := os.ReadFile("users.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Файл users.json не найден, создаем новый список пользователей")
			return []User{}, nil // Если файла нет, возвращаем пустой список
		}
		log.Println("Ошибка при чтении файла users.json:", err)
		return nil, err
	}

	log.Println("Содержимое файла users.json:", string(file))

	var users []User
	if err := json.Unmarshal(file, &users); err != nil {
		log.Printf("Ошибка при разборе JSON: %v", err)
		return nil, fmt.Errorf("ошибка при разборе JSON: %w", err)
	}

	if len(users) == 0 {
		log.Println("Файл users.json пуст или не содержит пользователей")
	}

	return users, nil
}

// SaveUsers сохраняет список пользователей в файл
func SaveUsers(users []User) error {
	file, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return os.WriteFile("users.json", file, 0644)
}
