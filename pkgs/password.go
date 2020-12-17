package pkgs

import (
	"errors"
	"unicode"
)

// Password validates a plain password against the rules that are defined as follows:
//
// upp: at least one upper case letter
// low: at least one lower case letter
// num: at least one digit
// sym: at least one special character
// tot: at least eight characters long
// No empty string or whitespaces
func Password(pass string) (bool, error) {
	var (
		upp, low, num, sym bool
		tot                uint8
	)

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			return false, errors.New("No whitespaces/empty inputs allowed")
		}
	}

	switch {
	case !upp:
		return false, errors.New("No Uppercase")
	case !low:
		return false, errors.New("No Lowercase")
	case !num:
		return false, errors.New("No Number")
	case !sym:
		return false, errors.New("No Symbol")
	case tot<8:
		return false, errors.New("Lesser than 8 characters")
	}

	return true, nil
}
