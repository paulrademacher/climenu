package climenu

import (
	"bufio"
	"fmt"
	"os"

	"github.com/buger/goterm"
)

type MenuItem struct {
	Text     string
	ID       string
	SubMenu  *Menu
	Selected bool  // For checkboxes.
}

const (
	ButtonType = iota
	CheckboxType
)

type Menu struct {
	Type      int
	Heading   string
	Question  string
	CursorPos int
	MenuItems []*MenuItem
}

type CheckboxMenu struct {
	Menu
	Yes string
	No  string
}

type ButtonMenu struct {
	Menu
}

func NewMenu(heading string, question string, menuType int) *Menu {
	return &Menu{
		MenuItems: make([]*MenuItem, 0),
		Heading:   heading,
		Question:  question,
		Type:      menuType,
	}
}

func NewButtonMenu(heading string, question string) *ButtonMenu {
	return &ButtonMenu{
		Menu: *NewMenu(heading, question, ButtonType),
	}
}

func NewCheckboxMenu(heading string, question string, yes string, no string) *CheckboxMenu {
	return &CheckboxMenu{
		Menu: *NewMenu(heading, question, CheckboxType),
		Yes:  yes,
		No:   no,
	}
}

func (m *Menu) AddMenuItem(text string, id string) *MenuItem {
	menuItem := &MenuItem{
		Text: text,
		ID:   id,
	}

	m.MenuItems = append(m.MenuItems, menuItem)

	return menuItem
}

func (m *Menu) CursorDown() {
	m.CursorPos = (m.CursorPos + 1) % len(m.MenuItems)
	m.DrawMenuItems(true)
}

func (m *Menu) CursorUp() {
	m.CursorPos = (m.CursorPos + len(m.MenuItems) - 1) % len(m.MenuItems)
	m.DrawMenuItems(true)
}

func (m *Menu) ToggleSelection() {
	menuItem := m.MenuItems[m.CursorPos]
	if menuItem.Selected {
		menuItem.Selected = false
	} else {
		menuItem.Selected = true
	}
	m.DrawMenuItems(true)
}

func (mi *MenuItem) SetSubMenu(menu *Menu) {
	mi.SubMenu = menu
}

func (m *Menu) Dump() {
	m.DumpIndent(0)
}

func printf(format string, s ...interface{}) {
	s = append(s, "\r\f")
	fmt.Printf(format+"%s", s...)
}

func (m *Menu) DrawMenuItems(redraw bool) {
	if redraw {
		// Move cursor up by number of menu items.  Assumes each menu item is ONE line.
		fmt.Printf("\033[%dA", len(m.MenuItems)-1)
	}

	for index, menuItem := range m.MenuItems {
		var newline = "\n"
		if index == len(m.MenuItems)-1 {
			newline = ""
		}

		prefix := "  "
		if m.Type == ButtonType {
			prefix = fmt.Sprintf("%d: ", index + 1)
		} else if m.Type == CheckboxType {
			if menuItem.Selected {
				prefix = "\u25c9 "
			} else {
				prefix = "\u25ef "
			}
		}

		if index == m.CursorPos {
			cursor := goterm.Color("> ", goterm.CYAN)
			fmt.Printf("\r%s%s %s%s", cursor, prefix, menuItem.Text, newline)
		} else {
			fmt.Printf("\r%s%s %s%s", "  ", prefix, menuItem.Text, newline)
		}
	}
}

func (m *Menu) Render() {
	if m.Heading != "" {
		fmt.Println(m.Heading)
	}

	fmt.Printf("%s\n", goterm.Color(goterm.Bold(m.Question) + ":", goterm.GREEN))

	m.DrawMenuItems(false)

	// Hide cursor.
	fmt.Printf("\033[?25l")
}

func printIndent(indent int) {
	for i := 0; i < indent; i++ {
		fmt.Printf("  ")
	}
}

func (m *Menu) DumpIndent(indent int) {
	printIndent(indent)
	fmt.Println("Menu: ", m.Question)

	for index, menuItem := range m.MenuItems {
		printIndent(indent)
		fmt.Printf("      %d: %s\n", index+1, menuItem.Text)
	}
}

func (m *ButtonMenu) Run() (string, bool) {
	results, escape := m.RunInternal()
	return results[0], escape
}

func (m *CheckboxMenu) Run() ([]string, bool) {
	results, escape := m.RunInternal()
	return results, escape
}

func (m *Menu) RunInternal() (results []string, escape bool) {
	defer func() {
		// Show cursor.
		fmt.Printf("\033[?25h")
	}()

	m.Render()

	for {
		ascii, keyCode, err := getChar()

		if (ascii == 3 || ascii == 27) || err != nil {
			fmt.Println()
			return []string{""}, true
		}

		if m.Type == ButtonType && ascii == 13 {
			fmt.Println()
			menuItem := m.MenuItems[m.CursorPos]
			return []string{menuItem.ID}, false
		}

		if m.Type == CheckboxType {
			if ascii == ' ' {
				m.ToggleSelection()
			} else if ascii == 13 {
				selections := make([]string, 0)
				for _, menuItem := range m.MenuItems {
					if menuItem.Selected {
						selections = append(selections, menuItem.ID)
					}
				}
				fmt.Println()
				return selections, false
			}
		}

		const one = 49
		const nine = 57

		if ascii >= one && ascii <= nine {
			num := ascii - one
			if num < len(m.MenuItems) {
				// Redraw items with new cursor selection.
				m.CursorPos = num
				m.DrawMenuItems(true)
				fmt.Println()

				menuItem := m.MenuItems[num]
				return []string{menuItem.ID}, false
			}
		}

		if keyCode == 40 {
			m.CursorDown()
		} else if keyCode == 38 {
			m.CursorUp()
		}
	}
}

func GetText(message string, defaultText string) string {
	fmt.Printf("%s", goterm.Color(goterm.Bold(message), goterm.GREEN))

	if defaultText != "" {
		fmt.Printf(" %s%s%s",
			goterm.Color(goterm.Bold("["), goterm.GREEN),
			goterm.Color(defaultText, goterm.YELLOW),
			goterm.Color(goterm.Bold("]"), goterm.GREEN))
	}

	fmt.Printf("%s ", goterm.Color(goterm.Bold(":"), goterm.GREEN))

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	if text[len(text)-1] == '\n' {
		text = text[:len(text)-1]
	}

	if text == "" {
		text = defaultText
	}

	return text
}
