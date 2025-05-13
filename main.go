package main

import (
	"fmt"
	"main/internal/config"
)

func main(){
	c, err := config.Read()
	if err != nil{
		fmt.Printf("%v", err)
	}
	c.SetUser("ivy")
	c, err = config.Read()
	if err != nil{
		fmt.Printf("%v", err)
	}
	fmt.Println(c.CurrentUserName)
	fmt.Println(c.DbUrl)
}