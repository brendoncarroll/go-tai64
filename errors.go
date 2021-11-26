package tai64

import "fmt"

func errWrongLength(ty string, have, want int) error {
	return fmt.Errorf("wrong length to be %s. HAVE: %d WANT: %d", ty, have, want)
}
