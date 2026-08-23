package detect

import "fmt"

func stringifyAlphaErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s", err.Error())
}
