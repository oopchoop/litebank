package config

import "errors"

func NewEnvError(name string) error {
	return errors.New("Не заполнена переменная окружения: " + name)
}
