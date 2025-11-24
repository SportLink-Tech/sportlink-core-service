package common

// Stats keeps track of the number of matches won, lost, and drawn by a team
type Stats struct {
	Wins   int
	Losses int
	Draws  int
}

func NewStats(wins int, losses int, draws int) *Stats {
	return &Stats{Wins: wins, Losses: losses, Draws: draws}
}

func (ts *Stats) TotalMatches() int {
	return ts.Wins + ts.Losses + ts.Draws
}
