package torrent

import "strconv"

var /*Types*/ (
	UNKNOWN = 0
	BSTR 	= 1
	BINT 	= 2
	BLIST 	= 3
	BDIR 	= 4
	BEND	= 5
)

func isNumeric(s uint8) bool {
	_, err := strconv.ParseFloat(string(s), 64)
	return err == nil
}

func GetType(c uint8) int {
	if isNumeric(c) {
		return BSTR
	} else if c == 'i' {
		return BINT
	} else if c == 'l' {
		return BLIST
	} else if c == 'd' {
		return BDIR
	} else if c == 'e' {
		return BEND
	}
	return UNKNOWN
}