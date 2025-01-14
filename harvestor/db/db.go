package db

import (
	"github.com/onrik/gorm-logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"harvestor/config"
	l "harvestor/logger"
	"time"
)

var dbHarvester *gorm.DB

func Init() {
	initHarvester()
	migrateHarvester()
}

func initHarvester() {
	// Create new logger
	var logger = l.NewLogger()
	// getting config for our db
	conf := config.GetConf()
	// trying to open sqlite DB
	d, err := gorm.Open(sqlite.Open(conf.Database.DBFile()), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	if err != nil {
		logger.Fatal("Can NOT open SQLite DB: ", conf.Database.DBFile(), " | ", err)
	}
	// assign to global
	dbHarvester = d
	// details
	db, err := d.DB()
	if err != nil {
		logger.Fatal("Can NOT get Stats for SQLite DB:", conf.Database.DBFile(), " | ", err)
	}

	// Controlling open connections
	db.SetMaxOpenConns(conf.Database.MaxOpenConnections())
	// Max Idle connections
	db.SetMaxIdleConns(conf.Database.MaxIdleConnections())
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Duration(conf.Database.MaxConnectionLifeTime()) * time.Minute)
	// All good here
	logger.Info("Harvester Database connected !!!")
}

func migrateHarvester() {
	// Create new logger
	var logger = l.NewLogger()
	// just a plcae holder
	logger.Info("Harvester Database checking for migration ... ")
	db := GetHarvesterDB()
	// Migrate the schema for file table
	err := db.AutoMigrate(&File{})
	if err != nil {
		logger.Fatal("Can NOT AutoMigrate `File` for SQLite DB:", err)
	}
	// Migrate the schema for Sidecar table
	err = db.AutoMigrate(&Sidecar{})
	if err != nil {
		logger.Fatal("Can NOT AutoMigrate `Sidecar` for SQLite DB:", err)
	}
	// Migrate the schema for upload table
	err = db.AutoMigrate(&Upload{})
	if err != nil {
		logger.Fatal("Can NOT AutoMigrate `Upload` for SQLite DB:", err)
	}
	// Migrate the schema for meta table
	err = db.AutoMigrate(&Meta{})
	if err != nil {
		logger.Fatal("Can NOT AutoMigrate `Meta` for SQLite DB:", err)
	}
	// forcing foreign keys in SQLite
	db.Exec("PRAGMA foreign_keys = ON;")
	// Done
	logger.Info("Harvester Database AutoMigrate is Done !!!")
}

// Get our db for sql operations
func GetHarvesterDB() *gorm.DB {
	return dbHarvester
}
