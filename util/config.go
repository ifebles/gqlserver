package util

import "os"

type config struct {
	Server struct {
		URI      string
		Database struct {
			Name       string
			Collection string
		}
	}
}

var Config = &config{}

func (c *config) Init() {
	c.Server.URI = os.Getenv("DB_URI")
	c.Server.Database.Name = os.Getenv("DB_NAME")
	c.Server.Database.Collection = os.Getenv("DB_COLLECTION")
}
