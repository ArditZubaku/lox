package main

import "errors"

var ErrUnterminatedString = errors.New("Unterminated string.")
var ErrInvalidNumberLiteral = errors.New("Invalid number literal.")
var ErrUnexpectedCharacter = errors.New("Unexpected character.")
var ErrUnterminatedComment = errors.New("Multi-line comment met EOF.")
