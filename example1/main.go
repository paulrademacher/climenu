package main

import (
	"fmt"

	"github.com/paulrademacher/climenu"
)

func callback(id string) {
	fmt.Println("Chose item:", id)

}

func main() {
	menu := climenu.NewMenu("Welcome", "Choose an action")
	menu.AddMenuItem("Create entry", "create", climenu.ActionItem)
	menu.AddMenuItem("Edit entry", "edit", climenu.ActionItem)

	action, escaped := menu.Run()
	if !escaped {
		fmt.Println("action>", action)
	}
}
