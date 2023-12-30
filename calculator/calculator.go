package calculator

import (
	"brawler-manager/brawlers"
	"sort"
	"sync"
)

const poolFactor = 20
const mapFactor = 40
const counterFactor = 40

func GenerateRanks(banned []string, friendly []string, enemy []string, gameMap string, matchUps map[string]float32, mapPct map[string]map[string]float32) []brawlers.Brawler {
	var wg sync.WaitGroup
	allBrawlers := brawlers.GetAllBrawlers()
	pool := CreatePool(allBrawlers, banned, friendly, enemy)
	indexes := map[string]int{}

	var poolStats []brawlers.Brawler
	var mapStats []brawlers.Brawler
	var counterStats []brawlers.Brawler
	wg.Add(1)
	go FindPercentAgainstAll(pool, matchUps, &poolStats, &wg)
	wg.Add(1)
	go FindPercentMap(pool, mapPct, gameMap, &mapStats, &wg)
	wg.Add(1)
	go FindPercentCounter(pool, enemy, matchUps, &counterStats, &wg)
	wg.Wait()
	for i, x := range poolStats {
		indexes[x.Name] = i
		x.WinPct *= poolFactor
	}
	for _, x := range mapStats {
		poolStats[indexes[x.Name]].WinPct += mapFactor * x.WinPct
	}
	for _, x := range counterStats {
		poolStats[indexes[x.Name]].WinPct += counterFactor * x.WinPct
	}
	sort.Slice(poolStats, func(i, j int) bool {
		return poolStats[i].WinPct > poolStats[j].WinPct // sort in decreasing order
	})
	return poolStats
}

func CreatePool(brawlersList []string, banned []string, friendly []string, enemy []string) []string {
	ignore := map[string]bool{}
	var pool []string
	for _, x := range banned {
		ignore[x] = true
	}
	for _, x := range friendly {
		ignore[x] = true
	}
	for _, x := range enemy {
		ignore[x] = true
	}
	for _, x := range brawlersList {
		_, ok := ignore[x]
		if !ok {
			pool = append(pool, x)
		}
	}
	return pool
}

func FindPercentAgainstAll(brawlerList []string, matchUps map[string]float32, result *[]brawlers.Brawler, wg *sync.WaitGroup) {
	defer wg.Done()
	stats := map[string]float32{}
	N := len(brawlerList)
	for i := 0; i < N-1; i++ {
		b1 := brawlerList[i]
		for j := i + 1; j < N; j++ {
			b2 := brawlerList[j]
			b1WinPct := FindWinPct(b1, b2, matchUps)
			b2WinPct := 1 - b1WinPct
			_, ok := stats[b1]
			if ok {
				stats[b1] += b1WinPct
			} else {
				stats[b1] = b1WinPct
			}
			_, ok = stats[b2]
			if ok {
				stats[b2] += b2WinPct
			} else {
				stats[b2] = b2WinPct
			}
		}
	}
	for k := range stats {
		stats[k] /= float32(N - 1)
	}
	var keyValuePairs []brawlers.Brawler
	for k, v := range stats {
		keyValuePairs = append(keyValuePairs, brawlers.Brawler{Name: k, WinPct: v})
	}
	*result = keyValuePairs
}

func FindPercentMap(brawlerList []string, mapPct map[string]map[string]float32, gameMap string, result *[]brawlers.Brawler, wg *sync.WaitGroup) {
	defer wg.Done()
	var keyValuePairs []brawlers.Brawler
	currentMap := mapPct[gameMap]
	for _, x := range brawlerList {
		winPct := currentMap[x]
		keyValuePairs = append(keyValuePairs, brawlers.Brawler{Name: x, WinPct: winPct})
	}
	*result = keyValuePairs
}

func FindPercentCounter(brawlerList []string, enemy []string, matchUps map[string]float32, result *[]brawlers.Brawler, wg *sync.WaitGroup) {
	defer wg.Done()
	var keyValuePairs []brawlers.Brawler
	if len(enemy) == 0 {
		for _, x := range brawlerList {
			keyValuePairs = append(keyValuePairs, brawlers.Brawler{Name: x, WinPct: 100})
		}
		result = &keyValuePairs
		return
	}
	for _, x := range brawlerList {
		winPct := float32(0)
		for _, e := range enemy {
			winPct += FindWinPct(x, e, matchUps)
		}
		keyValuePairs = append(keyValuePairs, brawlers.Brawler{Name: x, WinPct: winPct / float32(len(enemy))})
	}
	*result = keyValuePairs
}

func FindWinPct(target string, enemy string, matchUps map[string]float32) float32 {
	var pair string
	var first bool
	if target < enemy {
		pair = target + "_" + enemy
		first = true
	} else {
		pair = enemy + "_" + target
		first = false
	}
	matchUp := matchUps[pair]
	if first {
		return matchUp
	} else {
		return 1 - matchUp
	}
}
