package config

import (
	"encoding/json"
	"os"
)

type sendEmail struct {
	Server         string `json:"server" binding:"required"`
	Port           string `json:"port" binding:"required"`
	Sender         string `json:"sender" binding:"required"`
	PasswordSender string `json:"password_sender" binding:"required"`
}

type dataBase struct {
	NameDriver string `json:"name_driver" binding:"required"`
	Path       string `json:"path" binding:"required"`
}

type config struct {
	ModeStart string    `json:"mode_start" binding:"required"`
	Port      string    `json:"port" binding:"required"`
	SendEmail sendEmail `json:"send_email" binding:"required"`
	DataBase  dataBase  `json:"data_base" binding:"required"`
}

var AppConfig *config

func InitConfig(pathConfigFile string) error {

	configFile, err := os.Open(pathConfigFile)
	defer configFile.Close()
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&AppConfig)
	if err != nil {
		return err
	}

	return nil

}
