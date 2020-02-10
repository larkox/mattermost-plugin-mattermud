package mud

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
)
