package dedupe

import (
	"sort"
	"strconv"
	"strings"

	"proofpoint-flm/internal/domain"
	"proofpoint-flm/internal/normalize"
)

type Result struct {
	Episodes            []domain.Episode
	DuplicatesDetected  int
	StrategyDescription string
}

const strategyDescription = "Records are compared using normalized keys. A primary key uses (series, season, episode). Secondary keys are used only when one numeric field is missing: (series, 0, episode, title) when season is 0, and (series, season, 0, title) when episode is 0. For each duplicate cluster, the kept record is selected by priority: known AirDate, known EpisodeTitle, both season and episode known, then first appearance in input order."

func Apply(episodes []domain.Episode) Result {
	if len(episodes) == 0 {
		return Result{StrategyDescription: strategyDescription}
	}

	indexByKey := make(map[string]int, len(episodes)*2)
	uf := newUnionFind(len(episodes))

	for i, ep := range episodes {
		for _, key := range candidateKeys(ep) {
			if prev, ok := indexByKey[key]; ok {
				uf.union(i, prev)
			}
			indexByKey[key] = i
		}
	}

	clusters := make(map[int][]int)
	for i := range episodes {
		root := uf.find(i)
		clusters[root] = append(clusters[root], i)
	}

	selected := make([]domain.Episode, 0, len(clusters))
	duplicates := 0
	for _, indexes := range clusters {
		duplicates += len(indexes) - 1
		best := indexes[0]
		for _, idx := range indexes[1:] {
			if better(episodes[idx], episodes[best]) {
				best = idx
			}
		}
		selected = append(selected, episodes[best])
	}

	sortEpisodes(selected)
	return Result{
		Episodes:            selected,
		DuplicatesDetected:  duplicates,
		StrategyDescription: strategyDescription,
	}
}

func candidateKeys(ep domain.Episode) []string {
	series := normalize.NormalizeComparisonText(ep.SeriesName)
	title := normalize.NormalizeComparisonText(ep.EpisodeTitle)

	keys := []string{
		"A|" + series + "|" + strconv.Itoa(ep.SeasonNumber) + "|" + strconv.Itoa(ep.EpisodeNumber),
	}

	if ep.SeasonNumber == 0 {
		keys = append(keys, "B|"+series+"|0|"+strconv.Itoa(ep.EpisodeNumber)+"|"+title)
	}
	if ep.EpisodeNumber == 0 {
		keys = append(keys, "C|"+series+"|"+strconv.Itoa(ep.SeasonNumber)+"|0|"+title)
	}

	return keys
}

func better(candidate, current domain.Episode) bool {
	candidateAirKnown := normalize.IsAirDateKnown(candidate.AirDate)
	currentAirKnown := normalize.IsAirDateKnown(current.AirDate)
	if candidateAirKnown != currentAirKnown {
		return candidateAirKnown
	}

	candidateTitleKnown := normalize.IsTitleKnown(candidate.EpisodeTitle)
	currentTitleKnown := normalize.IsTitleKnown(current.EpisodeTitle)
	if candidateTitleKnown != currentTitleKnown {
		return candidateTitleKnown
	}

	candidateNumbersKnown := candidate.SeasonNumber > 0 && candidate.EpisodeNumber > 0
	currentNumbersKnown := current.SeasonNumber > 0 && current.EpisodeNumber > 0
	if candidateNumbersKnown != currentNumbersKnown {
		return candidateNumbersKnown
	}

	return candidate.InputOrder < current.InputOrder
}

func sortEpisodes(episodes []domain.Episode) {
	sort.Slice(episodes, func(i, j int) bool {
		left := episodes[i]
		right := episodes[j]

		leftSeries := strings.ToLower(left.SeriesName)
		rightSeries := strings.ToLower(right.SeriesName)
		if leftSeries != rightSeries {
			return leftSeries < rightSeries
		}
		if left.SeasonNumber != right.SeasonNumber {
			return left.SeasonNumber < right.SeasonNumber
		}
		if left.EpisodeNumber != right.EpisodeNumber {
			return left.EpisodeNumber < right.EpisodeNumber
		}

		leftTitle := strings.ToLower(left.EpisodeTitle)
		rightTitle := strings.ToLower(right.EpisodeTitle)
		if leftTitle != rightTitle {
			return leftTitle < rightTitle
		}

		return left.InputOrder < right.InputOrder
	})
}

type unionFind struct {
	parent []int
	rank   []int
}

func newUnionFind(size int) *unionFind {
	parent := make([]int, size)
	rank := make([]int, size)
	for i := range parent {
		parent[i] = i
	}
	return &unionFind{parent: parent, rank: rank}
}

func (u *unionFind) find(x int) int {
	if u.parent[x] != x {
		u.parent[x] = u.find(u.parent[x])
	}
	return u.parent[x]
}

func (u *unionFind) union(a, b int) {
	rootA := u.find(a)
	rootB := u.find(b)
	if rootA == rootB {
		return
	}

	if u.rank[rootA] < u.rank[rootB] {
		u.parent[rootA] = rootB
		return
	}
	if u.rank[rootA] > u.rank[rootB] {
		u.parent[rootB] = rootA
		return
	}

	u.parent[rootB] = rootA
	u.rank[rootA]++
}
