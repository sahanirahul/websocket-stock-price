package utility

import (
	"fmt"
)

func GetInstrumentKey(token int64) string {
	return fmt.Sprint(token)
}

func GetTokenKey(itoken string, itype string) string {
	return fmt.Sprintf("%s:%s", itype, itoken)
}
