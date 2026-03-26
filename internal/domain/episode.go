package domain

// Episode represents a cleaned, comparable episode record.
type Episode struct {
	SeriesName    string
	SeasonNumber  int
	EpisodeNumber int
	EpisodeTitle  string
	AirDate       string
	InputOrder    int
}

// RecordOutcome captures sanitize results for one raw row.
type RecordOutcome struct {
	Episode   Episode
	Discarded bool
	Corrected bool
}
