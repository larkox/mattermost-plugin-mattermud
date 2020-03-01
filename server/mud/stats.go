package mud

import (
	"encoding/json"
)

// Stat denotes one single stat (Strenght, Constitution...)
type Stat int

// Stats denotes all the stats for a single element
type Stats map[Stat]int

const (
	// Strength directly affect the physical attack and maximum weight to carry
	Strength Stat = iota
	// Constitution directly affect defense and HP
	Constitution
	// Dexterity directly affect evasion and hit rate
	Dexterity
	// Intelligence directly affect magical attack and maximum mana
	Intelligence
	// Wisdom directly affect magical defense and restoration power
	Wisdom
	// Luck affect critical rate, loot and others
	Luck
	// StatsLength is just used to check the number of stats. Any new stat should be added before this.
	StatsLength
)

// StatsJSON represents the stats in JSON format
type StatsJSON struct {
	Strength     int
	Constitution int
	Dexterity    int
	Intelligence int
	Wisdom       int
	Luck         int
}

// MarshalJSON marshals the stats into JSON
func (s Stats) MarshalJSON() ([]byte, error) {
	var sJSON StatsJSON
	sJSON.Strength = s[Strength]
	sJSON.Constitution = s[Constitution]
	sJSON.Dexterity = s[Dexterity]
	sJSON.Intelligence = s[Intelligence]
	sJSON.Wisdom = s[Wisdom]
	sJSON.Luck = s[Luck]
	return json.Marshal(sJSON)
}

// UnmarshalJSON unmarshal the stats from JSON
func (s Stats) UnmarshalJSON(b []byte) error {
	s = make(map[Stat]int)
	var sJSON StatsJSON
	err := json.Unmarshal(b, &sJSON)
	if err != nil {
		return err
	}
	s[Strength] = sJSON.Strength
	s[Constitution] = sJSON.Constitution
	s[Dexterity] = sJSON.Dexterity
	s[Intelligence] = sJSON.Intelligence
	s[Wisdom] = sJSON.Wisdom
	s[Luck] = sJSON.Luck

	return nil
}
