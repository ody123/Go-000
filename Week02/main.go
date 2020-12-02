package main

import (
	"fmt"
	"week02/service"
)

func main() {
	s := service.NewService()
	user, err := s.GetUserById(2)
	if err != nil {
		fmt.Printf("error info: %+v", err)
		return
	}
	fmt.Printf("user info: %+v", user)
}
