package climenu

import "github.com/pkg/term"

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
