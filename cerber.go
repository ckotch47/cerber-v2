package main

import (
	"cerber/internal/command"
	"cerber/internal/style"
	"fmt"
)

func main() {
	fmt.Println(style.TitleStyle.Render(style.Logo))
	command.Execute()
}
