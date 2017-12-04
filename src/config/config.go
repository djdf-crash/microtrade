package config

import (
	"encoding/json"
	"os"
)

type sendEmail struct {
	Server         string `json:"server,required"`
	Port           string `json:"port,required"`
	Sender         string `json:"sender,required"`
	PasswordSender string `json:"password_sender,required"`
}

type dataBase struct {
	NameDriver string `json:"name_driver,required"`
	Path       string `json:"path,required"`
}

type config struct {
	ModeStart string    `json:"mode_start,required"`
	Port      string    `json:"port,required"`
	SendEmail sendEmail `json:"send_email,required"`
	DataBase  dataBase  `json:"data_base,required"`
}

var AppConfig *config

func InitConfig(pathConfigFile string) error {

	configFile, err := os.Open(pathConfigFile)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(configFile)
	err = dec.Decode(&AppConfig)
	if err != nil {
		return err
	}

	return nil

}
