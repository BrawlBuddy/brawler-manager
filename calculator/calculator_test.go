package calculator

import (
	"brawler-manager/brawlers"
	"math/rand"
	"reflect"
	"sort"
	"sync"
	"testing"
)

var brawlerList = []string{"Colt", "Surge", "Hank", "Barley", "Shelly", "Jessie", "Spike"}
var maps = []string{"Split", "Super Beach", "Sneaky Fields"}
var matchUps = createTestData(brawlerList)
var mapPct = createMapData(brawlerList, maps)

func createTestData(brawlerNames []string) map[brawlers.Pair]brawlers.Matchup {
	b := make([]string, len(brawlerNames))
	copy(b, brawlerNames)
	sort.Slice(b, func(i, j int) bool {
		return b[i] < b[j] // sort in decreasing order
	})
	matchUps := map[brawlers.Pair]brawlers.Matchup{}
	for i := 0; i < len(b)-1; i++ {
		for j := i + 1; j < len(b); j++ {
			wins := rand.Intn(301) + 100
			matchUps[brawlers.Pair{b[i], b[j]}] = brawlers.Matchup{
				b[i],
				b[j],
				wins,
				500 - wins}

		}
	}
	return matchUps
}

func createMapData(brawlerNames []string, maps []string) map[brawlers.MapPair]float32 {
	mapPct := map[brawlers.MapPair]float32{}
	for _, x := range brawlerNames {
		for _, m := range maps {
			mapPct[brawlers.MapPair{m, x}] = 0.15 + rand.Float32()*(0.85-0.15)
		}
	}
	return mapPct
}

func TestCreatePool(t *testing.T) {
	banned := []string{"Hank", "Barley"}
	friendly := []string{"Spike"}
	enemy := []string{"Shelly", "Colt"}
	case1 := CreatePool(brawlerList, banned, friendly, enemy)
	sort.Strings(case1)
	case2 := CreatePool(brawlerList, banned, []string{}, enemy)
	sort.Strings(case2)
	case3 := CreatePool(brawlerList, banned, []string{}, []string{})
	sort.Strings(case3)
	expected1 := []string{"Jessie", "Surge"}
	sort.Strings(expected1)
	expected2 := []string{"Jessie", "Surge", "Spike"}
	sort.Strings(expected2)
	expected3 := []string{"Jessie", "Surge", "Spike", "Shelly", "Colt"}
	sort.Strings(expected3)

	if !reflect.DeepEqual(case1, expected1) {
		t.Error("Not equal. Result:", case1, "Expected:", expected1)
	}
	if !reflect.DeepEqual(case2, expected2) {
		t.Error("Not equal. Result:", case2, "Expected:", expected2)
	}
	if !reflect.DeepEqual(case3, expected3) {
		t.Error("Not equal. Result:", case3, "Expected:", expected3)
	}
}

func TestFindPercentAgainstAll(t *testing.T) {
	var wg sync.WaitGroup
	var poolStats []brawlers.Brawler
	wg.Add(1)
	go FindPercentAgainstAll([]string{"Hank", "Surge", "Barley"}, matchUps, &poolStats, &wg)
	wg.Wait()
	t.Log("PoolStats:", poolStats)
	hankSurge := matchUps[brawlers.Pair{"Hank", "Surge"}]
	barleyHank := matchUps[brawlers.Pair{"Barley", "Hank"}]
	barleySurge := matchUps[brawlers.Pair{"Barley", "Surge"}]
	hankSurgePct := float32(hankSurge.Brawler1wins) / float32(hankSurge.Brawler1wins+hankSurge.Brawler2wins)
	barleyHankPct := float32(barleyHank.Brawler1wins) / float32(barleyHank.Brawler1wins+barleyHank.Brawler2wins)
	barleySurgePct := float32(barleySurge.Brawler1wins) / float32(barleySurge.Brawler1wins+barleySurge.Brawler2wins)
	t.Log("RealStats:")
	t.Log("HankSurge:", hankSurgePct, "--", 1-hankSurgePct)
	t.Log("BarleySurge:", barleySurgePct, "--", 1-barleySurgePct)
	t.Log("BarleyHank:", barleyHankPct, "--", 1-barleyHankPct)
	t.Log("hank wins", hankSurge.Brawler1wins)
}

func TestFindPercentMap(t *testing.T) {
	var wg sync.WaitGroup
	var mapStats []brawlers.Brawler
	wg.Add(1)
	go FindPercentMap(brawlerList, mapPct, maps[0], &mapStats, &wg)
	wg.Wait()
	t.Log(mapStats)
	for k, v := range mapPct {
		if k.Map != maps[0] {
			continue
		}
		t.Log(k.Brawler, v)
	}
}

func TestFindPercentCounter(t *testing.T) {
	var wg sync.WaitGroup
	var counterStats []brawlers.Brawler
	wg.Add(1)
	go FindPercentCounter([]string{"Surge", "Colt"}, []string{}, matchUps, &counterStats, &wg)
	wg.Wait()
	for _, x := range counterStats {
		if x.WinPct != 100 {
			t.Error("Not equal to 100, result:", x.WinPct)
		}
	}
}
