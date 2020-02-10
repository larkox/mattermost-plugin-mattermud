package mud

// Player represents one single player
type Player struct {
	// UserID is Mattermos UserID
	UserID string
	// Name is the character name shown in the game
	Name string
	// Stats are the base stats of the character
	Stats Stats
	// Class is the character class
	Class PlayerClass
	// Race is the character race
	Race Race
	// Level is the current experience leve
	Level int
	// Experience how many experience points the player has. It is used for levelling up
	Experience int
	// IsSleeping shows whether the player is sleeping and should receive messages from the bot or not
	IsSleeping bool
	// Inventory contains all the items carried by the character
	Inventory []*Item
	// Equip contains the currently equipped items
	Equip PlayerEquipment
	// Effects show all the magical effects that the character is currently under
	Effects EffectList
	// CurrentRoom shows on which room the player is currently on
	CurrentRoom *Room
}

// GetLeftAttack returns the attack with the weapon on the left hand
func (p *Player) GetLeftAttack() int {
	baseAtt := p.Equip.GetLeftAttack()
	if baseAtt == 0 {
		return 0
	}
	str := p.GetCurrentStat(Strength)
	attEquipModifiers := p.Equip.GetAttackModifiers()
	attEffectModifiers := p.Effects.GetAttackModifiers()
	return min(0, str+baseAtt+attEquipModifiers+attEffectModifiers)
}

// GetRightAttack returns the attack with the weapon on the right hand
func (p *Player) GetRightAttack() int {
	baseAtt := p.Equip.GetRightAttack()
	if baseAtt == 0 {
		return 0
	}
	str := p.GetCurrentStat(Strength)
	attEquipModifiers := p.Equip.GetAttackModifiers()
	attEffectModifiers := p.Effects.GetAttackModifiers()
	return min(0, str+baseAtt+attEquipModifiers+attEffectModifiers)
}

// GetCurrentStat returns the current stat of the character
func (p *Player) GetCurrentStat(s Stat) int {
	base := p.Stats[s]
	equipModifiers := p.Equip.GetStatModifiers(s)
	effectModifiers := p.Effects.GetStatModifiers(s)
	return min(0, base+equipModifiers+effectModifiers)
}

// CanSeeHidden returns whether the character can see hidden objects
func (p *Player) CanSeeHidden() bool {
	return p.Equip.CanSeeHidden() || p.Effects.CanSeeHidden()
}

// CanSeeInvisible returns whether the character can see invisible objects
func (p *Player) CanSeeInvisible() bool {
	return p.Equip.CanSeeInvisible() || p.Effects.CanSeeInvisible()
}

// IsHidden returns whether the character is hidden
func (p *Player) IsHidden() bool {
	return p.Equip.GrantHidden() || p.Effects.GrantHidden()
}

// IsInvisible returns whether the character is invisible
func (p *Player) IsInvisible() bool {
	return p.Equip.GrantInvisible() || p.Effects.GrantInvisible()
}

// Move moves a character in certain direction, and returns the message to show to the player
func (p *Player) Move(d Direction) string {
	if !p.CurrentRoom.CanMove(d, p.CanSeeHidden(), p.CanSeeInvisible()) {
		if p.CanSeeDoor(d) {
			return "The door is locked."
		}
		return "You cannot go in that direction."
	}

	p.CurrentRoom = p.CurrentRoom.GetNeighbourRoom(d)
	return p.CurrentRoom.String()
}

// CanSeeDoor checks whether a locked door can be seen in certain direction
func (p *Player) CanSeeDoor(d Direction) bool {
	return p.CurrentRoom.CanSeeDoor(d, p.CanSeeHidden(), p.CanSeeInvisible())
}

// Look returns the current room long description
func (p *Player) Look() string {
	return p.CurrentRoom.LongDescription
}

// GetRoom returns the string for the current room
func (p *Player) GetRoom() string {
	return p.CurrentRoom.String()
}
