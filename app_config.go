package main

type AppConfig struct {
	Hosts  []string `json:"hosts"`
	Influx struct {
		URL    string `json:"url"`
		Secret string `json:"secret"`
		Org    string `json:"org"`
		Bucket string `json:"bucket"`
	} `json:"influx"`
}
