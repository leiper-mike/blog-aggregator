package config

import (
	"encoding/json"
	"fmt"
	"os"
)
type Config struct{
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}


func Read()  (Config, error){
	var config Config
	path, err := os.UserHomeDir()
	if err != nil{
		return config, fmt.Errorf("error creating path: %v", err)
	}
	path = path + "/gatorconfig.json"
	data, err := os.ReadFile(path)
	if err != nil{
		return config, fmt.Errorf("error reading file: %v", err)
	}
	
	jsonErr := json.Unmarshal(data, &config)
	if jsonErr != nil{
		return config, fmt.Errorf("error parsing json: %v", jsonErr)
	}
	return config, nil
}
func (c Config) SetUser(user string) error{
	if user == ""{
		return fmt.Errorf("user must be non empty")
	}
	c.CurrentUserName = user
	data, err := json.Marshal(c)
	if err != nil{
		return fmt.Errorf("error encoding config: %v", err)
	}
	path, err := os.UserHomeDir()
	if err != nil{
		return fmt.Errorf("error creating path: %v", err)
	}
	path = path + "/gatorconfig.json"
	writeError := os.WriteFile(path, data, os.ModePerm)
	if writeError != nil{
		return fmt.Errorf("error writing data to file: %v", writeError)
	}
	return nil
}