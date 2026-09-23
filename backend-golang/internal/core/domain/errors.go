package domain

import "errors"

var ErrStateNotFound = errors.New("estado nao encontrado")
var ErrCityNotFound = errors.New("municipio nao encontrado")

var ErrInvalidOrdenar = errors.New("ordenar invalido")
var ErrInvalidOrdem = errors.New("ordem invalida")
