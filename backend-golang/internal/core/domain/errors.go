package domain

import "errors"

var ErrStateNotFound = errors.New("estado nao encontrado")
var ErrCityNotFound = errors.New("municipio nao encontrado")
var ErrAgeGroupNotFound = errors.New("faixa etária não encontrada")

var ErrInvalidOrdenar = errors.New("ordenar invalido")
var ErrInvalidOrdem = errors.New("ordem invalida")
