package alert

import (
	"encoding/json"
	"github.com/lordralex/absol/database"
	"github.com/lordralex/absol/logger"
	"strings"
)

func ImportFromDatabase() {
	db, err := database.Get()
	if err != nil {
		logger.Err().Printf("Error connecting to database: %s\n", err.Error())
		return
	}

	var data []log
	err = db.Table("sites_timed_out").Find(&data).Error
	if err != nil {
		logger.Err().Printf("Error connecting to database: %s\n", err.Error())

		//try it again
		err = db.Table("sites_timed_out").Find(&data).Error
		if err != nil {
			logger.Err().Printf("Error connecting to database: %s\n", err.Error())
			return
		}
	}

	for _, d := range data {
		r := strings.NewReader(d.Log)

		var m map[string]interface{}

		err = json.NewDecoder(r).Decode(&m)
		if err != nil {
			logger.Err().Printf("Error decoding: %s\n", err.Error())
			continue
		}

		err = submitToElastic(d.Id, m)
		if err != nil {
			logger.Err().Printf("Error sending to ES: %s\n", err.Error())
		}
	}
}

type log struct {
	Id string `gorm:"identifier"`
	Log string `gorm:"log"`
}
