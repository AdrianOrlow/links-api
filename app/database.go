package app

import (
	"fmt"
	"github.com/AdrianOrlow/links-api/app/model"
	"github.com/AdrianOrlow/links-api/config"
	"github.com/jinzhu/gorm"
)

func (a *App) InitializeDatabase(config *config.Config) error {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		return err
	}

	a.DB = model.DBMigrate(db)
	return nil
}