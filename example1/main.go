package main

import (
	"fmt"
	"os"

	"github.com/paulrademacher/climenu"
)

func callback(id string) {
	fmt.Println("Chose item:", id)

}

func main() {
	menu := climenu.NewButtonMenu("Welcome", "Choose an action")
	menu.AddMenuItem("Create entry", "create")
	menu.AddMenuItem("Edit entry", "edit")

	action, escaped := menu.Run()
	if escaped {
		os.Exit(0)
	}

	fmt.Println("action >", action)

	checkbox := climenu.NewCheckboxMenu("Let's try some checkboxes",
		"Select options", "OK", "Cancel")
	checkbox.AddMenuItem("Apples", "apples")
	checkbox.AddMenuItem("Oranges", "oranges")
	checkbox.AddMenuItem("Bananas", "bananas")

	selection, escaped := checkbox.Run()
	if escaped {
		os.Exit(0)
	}

	fmt.Println("selected >", selection)

	response := climenu.GetText("Say something interesting", "hi")
	if escaped {
		os.Exit(0)
	}

	fmt.Printf("text > \"%s\"\n", response)
}
