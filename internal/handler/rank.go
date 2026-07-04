package handler

type rankTier struct {
	name	string
	low		float64
	high	float64
}

var rankTiers = []rankTier{
	{"Iron I", 0, 0.1},
	{"Iron II", 0.1, 0.2},
	{"Iron III", 0.2, 0.3},
	{"Bronze I", 0.3, 0.4},
	{"Bronze II", 0.4, 0.5},
	{"Bronze III", 0.5, 0.6},
	{"Silver I", 0.6, 0.7},
	{"Silver II", 0.7, 0.8},
	{"Silver III", 0.8, 0.9},
	{"Gold I", 0.9, 1.0},
	{"Gold II", 1.0, 1.1},
	{"Gold III", 1.1, 1.2},
	{"Platinum I", 1.2, 1.3},
	{"Platinum II", 1.3, 1.4},
	{"Platinum III", 1.4, 1.5},
	{"Diamond I", 1.5, 1.6},
	{"Diamond II", 1.6, 1.7},
	{"Diamond III", 1.7, 1.8},
	{"Spiker", 1.8, 1.9},
	{"Ace", 1.9, 2.0},
}

func calculateRank(played, points int32, efficiencyRate float64) string {
	if played < 10 {
		return "Unranked"
	}

	avgPoints := 0.0
	if played > 0 {
		avgPoints = float64(points) / float64(played)
	}

	mmr := avgPoints * efficiencyRate

	for _, tier := range rankTiers {
		if mmr >= tier.low && mmr < tier.high {
			return tier.name
		}
	}
	if mmr >= 2.0 {
		return "Sensei"
	}
	return "Iron I"
}
