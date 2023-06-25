package utility

import "encoding/json"

func MapObjectToAnother(fromObj interface{}, toObj interface{}) error {
	b, err := json.Marshal(fromObj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, toObj)
	return err
}

const (
	TOKENFORALLUNDERLYING = "ALLUNDERYINGS"
	EQUITY                = "EQ"
	DERIVATIVES           = "DERIVATIVES"
	DataTypeQuote         = "quote"
	DataTypePing          = "ping"
	DataTypeError         = "error"
	MsgCommandSubscribe   = "subscribe"
	MsgCommandUnSubscribe = "unsubscribe"
)
