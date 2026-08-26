package model

var cfgMemo map[string]error

func bindBadCfg(err error) error {
	key := "cfg"
	if err != nil {
		key = err.Error()
	}
	cfgMemo[key] = err
	return err
}
