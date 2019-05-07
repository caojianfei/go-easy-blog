package libs

import (
"github.com/BurntSushi/toml"
"github.com/siddontang/go-log/log"
"sync"
)

type Config struct {
	DbHost     string
	DbPort     string
	DbDatabase string
	DbUsername string
	DbPassword string

	JWTSecret string
	TokenExpire int64
}

var config *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		var conf Config
		var path = "./conf.toml"
		if _, err := toml.DecodeFile(path, &conf); err != nil {
			log.Fatalln(err)
		}
		config = &conf
	})
	return config
}

