package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	FirebaseConfig `json:"firebaseConfig"`
	FirebaseSDKServiceAccount
}

type FirebaseConfig struct {
	FIREBASE_API_KEY             string `json:"apiKey"`
	FIREBASE_AUTH_DOMAIN         string `json:"authDomain"`
	FIREBASE_PROJECT_ID          string `json:"projectId"`
	FIREBASE_STORAGE_BUCKET      string `json:"storageBucket"`
	FIREBASE_MESSAGING_SENDER_ID string `json:"messagingSenderId"`
	FIREBASE_APP_ID              string `json:"appId"`
	FIREBASE_MEASUREMENT_ID      string `json:"measurementId"`
}

type FirebaseSDKServiceAccount struct {
	FIREBASE_SERVICE_ACCOUNT_LOCATION string //location of the sdk service account.json file
}

func GetConfig() *Config {
	rtn := Config{}
	rtn.FIREBASE_API_KEY = os.Getenv("FIREBASE_API_KEY")
	rtn.FIREBASE_AUTH_DOMAIN = os.Getenv("FIREBASE_AUTH_DOMAIN")
	rtn.FIREBASE_PROJECT_ID = os.Getenv("FIREBASE_PROJECT_ID")
	rtn.FIREBASE_STORAGE_BUCKET = os.Getenv("FIREBASE_STORAGE_BUCKET")
	rtn.FIREBASE_MESSAGING_SENDER_ID = os.Getenv("FIREBASE_MESSAGING_SENDER_ID")
	rtn.FIREBASE_APP_ID = os.Getenv("FIREBASE_APP_ID")
	rtn.FIREBASE_MEASUREMENT_ID = os.Getenv("FIREBASE_MEASUREMENT_ID")
	rtn.FIREBASE_SERVICE_ACCOUNT_LOCATION = os.Getenv("FIREBASE_SERVICE_ACCOUNT_LOCATION")
	return &rtn
}

func (c *Config) FirebaseConfigAsJsonBytes() []byte {
	jsonData, err := json.Marshal(c.FirebaseConfig)
	if err != nil {
		fmt.Println("unable to convert firebase config to json")
	}
	return jsonData
}

func (c *Config) FirebaseConfigAsJsonString() string {
	jsonData := c.FirebaseConfigAsJsonBytes()
	return string(jsonData)
}
