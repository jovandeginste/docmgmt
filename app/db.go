package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func (a *App) InitDB() error {
	a.Logf(LOG_DEBUG, "Opening database, using '%s'", a.Configuration.MetadataDB)
	db, err := gorm.Open("sqlite3", a.Configuration.MetadataDB)
	if err != nil {
		return err
	}

	a.Logf(LOG_DEBUG, "Initializing database schema")
	// Migrate the schema
	db.AutoMigrate(&Info{})
	db.AutoMigrate(&Body{})
	db.AutoMigrate(&Metadata{})
	db.AutoMigrate(&FileMetadata{})

	db = db.Set("gorm:auto_preload", true)

	a.DB = db
	return db.Error
}
