package repositories

import "errors"

var (
	ErrBookNotFound = errors.New("kitap bulunamadı")
	ErrUserNotFound = errors.New("kullanıcı bulunamadı")
)
