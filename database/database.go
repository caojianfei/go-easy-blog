package database

import (
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/siddontang/go-log/log"
	"sync"
)

type DB struct {
	*sql.DB
}

type DbConfig struct {
	DbHost     string
	DbPort     string
	DbDatabase string
	DbUsername string
	DbPassword string
}

var instance *DB
var once sync.Once

func New() *DB {
	once.Do(func() {
		db, err := getInstance()
		if err != nil {
			log.Fatal(err)
		}
		instance = db
	})

	return instance
}

func getInstance() (*DB, error) {
	db := &DB{}
	conf, err := getDbConfig()
	if err != nil {
		return db, err
	}

	mysql, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.DbUsername, conf.DbPassword, conf.DbHost, conf.DbPort, conf.DbDatabase))

	if err != nil {
		return db, err
	}

	db.DB = mysql
	return db, nil
}

func getDbConfig() (DbConfig, error) {
	var conf DbConfig
	var path = "./conf.toml"
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}
