package brawlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type BrawlerJSON struct {
	Name string `json:"name"`
}

func GetMatchUps() map[string]float32 {
	matchUpList, err := ioutil.ReadFile("./data/1v1.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload map[string]float32
	err = json.Unmarshal(matchUpList, &payload)
	if err != nil {
		log.Fatal("Error when unmarshalling file: ", err)
	}

	return payload
}

func GetMapData() map[string]map[string]float32 {
	mapList, err := ioutil.ReadFile("./data/map.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payload map[string]map[string]float32
	err = json.Unmarshal(mapList, &payload)
	if err != nil {
		log.Fatal("Error when unmarshalling file: ", err)
	}
	return payload
}

func GetAllBrawlers() []string {
	allBrawlers := []string{}
	brawlerList, err := ioutil.ReadFile("./data/BrawlerList.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payload []BrawlerJSON
	err = json.Unmarshal(brawlerList, &payload)
	if err != nil {
		log.Fatal("Error when unmarshalling file: ", err)
	}
	for _, brawler := range payload {
		allBrawlers = append(allBrawlers, brawler.Name)
	}
	return allBrawlers
}
