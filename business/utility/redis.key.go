package utility

import (
	"fmt"
)

func GetInstrumentKey(token int64) string {
	return fmt.Sprint(token)
}

func GetTokenKey(isymbol, itype string) string {
	return fmt.Sprintf("%s:%s", itype, isymbol)
}
