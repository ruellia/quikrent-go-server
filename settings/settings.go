package settings

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/ruellia/quikrent-go-server/util"

	uuid "github.com/satori/go.uuid"
)

// Settings are the settings used to query the database.
type Settings struct {
	BotSettings  BotSettings
	AbsolutePath string
}

// BotSettings are the settings that will be passed into the Slackbot.
type BotSettings struct {
	MinPrice                 float64                `json:"min_price"`
	MaxPrice                 float64                `json:"max_price"`
	SlackToken               string                 `json:"slack_token"`
	Bedrooms                 float64                `json:"bed"`
	Bathrooms                float64                `json:"bath"`
	CraigslistSite           string                 `json:"craigslist_site"`
	CraigslistHousingSection string                 `json:"craigslist_housing_section"`
	Areas                    []string               `json:"areas"`
	MaxTransitDistance       float64                `json:"max_transit_distance"`
	Neighborhoods            []string               `json:"neighborhoods"`
	TransitStations          map[string][]float64   `json:"transit_stations"`
	Boxes                    map[string][][]float64 `json:"boxes"`
}

// CreateJSONFile creates a JSON file based on the BotSettings struct.
func CreateJSONFile(settings *Settings) error {
	marshaled, err := json.Marshal(settings.BotSettings)
	if err != nil {
		return err
	}
	fileName := uuid.NewV4().String() + ".json"
	fileName, err = filepath.Abs(util.UserHomeDir() + "/" + fileName)
	if err != nil {
		return err
	}
	settings.AbsolutePath = fileName
	if err := ioutil.WriteFile(fileName, marshaled, 0644); err != nil {
		return err
	}
	return nil
}

// UpdateJSONFile updates a JSON file based on the BotSettings struct.
func UpdateJSONFile(settings Settings) error {
	marshaled, err := json.Marshal(settings.BotSettings)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(settings.AbsolutePath, marshaled, 0644); err != nil {
		return err
	}
	return nil
}

// ConvertJSONRequest converts a JSON HTTP request into our Settings struct.
func ConvertJSONRequest(r *http.Request) (Settings, error) {
	decoder := json.NewDecoder(r.Body)
	var converted BotSettings
	if err := decoder.Decode(&converted); err != nil {
		return Settings{}, err
	}
	return Settings{BotSettings: converted}, nil
}
