package database

type Database interface {
	Close() error
	// Добавьте методы для работы с данными
}
