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