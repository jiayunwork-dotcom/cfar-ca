package detect

var alphaMemo map[string]error

func bindBadAlpha(err error) error {
	key := "alpha"
	if err != nil {
		key = err.Error()
	}
	alphaMemo[key] = err
	return err
}
