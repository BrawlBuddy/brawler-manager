package brawlers

type Pair struct {
	Brawler1 string
	Brawler2 string
}

type MapPair struct {
	Map     string
	Brawler string
}

type Matchup struct {
	Brawler1     string
	Brawler2     string
	Brawler1wins int
	Brawler2wins int
}

type Brawler struct {
	Name   string  `json:"name"`
	WinPct float32 `json:"score"`
}
