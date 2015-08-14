package climenu

import (
	"fmt"

	"github.com/buger/goterm"
	"github.com/pkg/term"
)

const (
	ActionItem = iota
	CheckboxItem
)

type MenuItem struct {
	Text    string
	ID      string
	Type    int
	SubMenu *Menu
}

type Menu struct {
	Heading        string
	Question       string
	CursorPos      int
	MenuItems      []*MenuItem
}

func NewMenu(heading string, question string) *Menu {
	return &Menu{
		MenuItems:      make([]*MenuItem, 0),
		Heading:        heading,
		Question:       question,
	}
}

func (m *Menu) AddMenuItem(text string, id string, itemType int) *MenuItem {
	menuItem := &MenuItem{
		Text: text,
		ID:   id,
		Type: itemType,
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

func getChar() (ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 1)
	_, err = t.Read(bytes)
	if err != nil {
		return
	}
	if bytes[0] == 27 {
		_, err = t.Read(bytes)
		if err != nil {
			return
		}
		if bytes[0] == 91 {
			_, err = t.Read(bytes)
			if err != nil {
				return
			}

			if bytes[0] == 65 {
				keyCode = 38
			} else if bytes[0] == 66 {
				keyCode = 40
			} else if bytes[0] == 67 {
				keyCode = 39
			} else if bytes[0] == 68 {
				keyCode = 37
			}
		}
	} else {
		ascii = int(bytes[0])
	}
	t.Restore()
	t.Close()
	return
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

		t := fmt.Sprintf("%d: %s", index+1, menuItem.Text)
		if index == m.CursorPos {
			fmt.Printf("\r%s %s%s", goterm.Color(">", goterm.GREEN), t, newline)
		} else {
			fmt.Printf("\r  %s%s", t, newline)
		}
	}

	//	fmt.Printf("\r")
}

func (m *Menu) Render() {
	if m.Heading != "" {
		fmt.Println(m.Heading)
	}

	fmt.Printf("[%s] %s:\n", goterm.Color("?", goterm.GREEN), m.Question)

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

func (m *Menu) Run() (result string, escape bool) {
	defer func() {
		// Show cursor.
		fmt.Printf("\033[?25h")
	}()

	m.Render()

	for {
		ascii, keyCode, err := getChar()

		if ascii == 3 || err != nil {
			fmt.Println()
			return "", true
		}

		if ascii == 13 {
			fmt.Println()

			menuItem := m.MenuItems[m.CursorPos]

			return menuItem.ID, false
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
				return menuItem.ID, false
			}
		}

		if keyCode == 40 {
			m.CursorDown()
		} else if keyCode == 38 {
			m.CursorUp()
		}
	}
}
