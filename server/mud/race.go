package mud

// Race denotes which race a character is (human, elf, dwarf...)
type Race int

const (
	// Human are the most equilibrated class, excelling only in Luck
	Human Race = iota
	// Elf have more dexterity, intelligence and wisdom, but have little strength or constitution
	Elf
	// Dwarf have more strenght, constitution and wisdom, but have little dexterity or intelligence
	Dwarf
)
