package mud

// PlayerClass denotes the class of a player (Warrior, Mage, Rogue...)
type PlayerClass int

const (
	// Warrior excells in Strength and Constitution and related skills
	Warrior PlayerClass = iota
	// Mage excells in Intelligence and Wisdom and related skills
	Mage
	// Rogue excells in Dexterity and Luck and related skills
	Rogue
)
