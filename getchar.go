package climenu

import "fmt"
import "github.com/pkg/term"

func getChar() (ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		if bytes[2] == 65 {
			keyCode = 38
		} else if bytes[2] == 66 {
			keyCode = 40
		} else if bytes[2] == 67 {
			keyCode = 39
		} else if bytes[2] == 68 {
			keyCode = 37
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
	} else {
		// Two characters read??
	}
	t.Restore()
	t.Close()
	return
}
