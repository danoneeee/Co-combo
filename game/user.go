package game

import (
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	Name  string `json:"name"`
	Coins int    `json:"coins"`
}

func LoadUser() (*User, error) {
	file, err := os.ReadFile("game/user.json")
	if err != nil {
		return nil, err
	}
	var user User
	if err := json.Unmarshal(file, &user); err != nil {
		return nil, err
	}
	fmt.Println(user)
	return &user, nil
}

// Сохраняем данные пользователя в файл
func SaveUser(user *User) error {
	file, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return os.WriteFile("game/user.json", file, 0644)
}
