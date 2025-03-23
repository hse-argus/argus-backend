package schedule

import "time"

var (
	mapping = map[string]time.Duration{
		"1m":  time.Minute,
		"5m":  time.Minute * 5,
		"15m": time.Minute * 15,
		"1h":  time.Hour,
		"4h":  time.Hour * 4,
		"8h":  time.Hour * 8,
		"12h": time.Hour * 12,
		"24h": time.Hour * 24,
	}
)

func (s *Schedule) ParseScheduleDuration() (time.Duration, bool) {
	parsed, ok := mapping[s.Time]
	return parsed, ok
}
