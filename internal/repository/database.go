package repository

import "gorm.io/gorm"

//go:generate mockgen -source=database.go -destination=database_mock.go -package repository
type Database interface {
	Create(value interface{}) (tx *gorm.DB)
	Rollback() *gorm.DB
	Commit() *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, where ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	First(dest interface{}, where ...interface{}) *gorm.DB
}
