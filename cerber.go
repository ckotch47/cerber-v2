package main

import (
	"fmt"

	"cerber/internal/command"
	"cerber/internal/style"
)


    

func main() {
	fmt.Println(style.TitleStyle.Render(style.Logo))
	
	command.Execute()
}