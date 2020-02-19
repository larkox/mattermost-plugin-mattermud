package mud

import "fmt"

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
	// Notify sends a message to a player
	Notify func(message string)
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
func (p *Player) Move(d Direction) {
	if p.IsSleeping {
		p.Notify("You cannot sleepwalk.")
		return
	}

	if !p.CurrentRoom.CanMove(d, p.CanSeeHidden(), p.CanSeeInvisible()) {
		if p.CanSeeDoor(d) {
			p.Notify("The door is locked.")
			return
		}
		p.Notify("You cannot go in that direction.")
		return
	}

	p.CurrentRoom.Exit(p, d)
	p.CurrentRoom = p.CurrentRoom.GetNeighbourRoom(d)
	p.CurrentRoom.Enter(p, d)
	p.ShowRoom()
}

// CanSeeDoor checks whether a locked door can be seen in certain direction
func (p *Player) CanSeeDoor(d Direction) bool {
	return p.CurrentRoom.CanSeeDoor(d, p.CanSeeHidden(), p.CanSeeInvisible())
}

// LookRoom returns the current room long description
func (p *Player) LookRoom() {
	if p.IsSleeping {
		p.Notify("No matter how hard you look, you see nothing while asleep.")
		return
	}
	p.Notify(p.CurrentRoom.Show(p.UserID, p.CanSeeHidden(), p.CanSeeInvisible(), true))
}

// ShowRoom returns the string for the current room
func (p *Player) ShowRoom() {
	if p.IsSleeping {
		p.Notify("You cannot see much while sleeping.")
		return
	}
	p.Notify(p.CurrentRoom.Show(p.UserID, p.CanSeeHidden(), p.CanSeeInvisible(), false))
}

// Show returns the string of how the user is seen
func (p *Player) Show(canSeeHidden, canSeeInvisible bool) string {
	if (!canSeeHidden && p.IsHidden()) ||
		(!canSeeInvisible && p.IsInvisible()) {
		return ""
	}

	if p.IsSleeping {
		return fmt.Sprintf("%s is sleeping here.", p.Name)
	}

	return fmt.Sprintf("%s is here.", p.Name)
}

// NotifyExitingPlayer checks if the exitingPlayer can be seen, and sends a message to the player.
func (p *Player) NotifyExitingPlayer(exitingPlayer *Player, d Direction) {
	if p.IsSleeping {
		return
	}

	if (!exitingPlayer.IsHidden() || p.CanSeeHidden()) &&
		(!exitingPlayer.IsInvisible() || p.CanSeeInvisible()) {
		var dString string
		switch d {
		case North:
			dString = "North"
		case South:
			dString = "South"
		case West:
			dString = "West"
		case East:
			dString = "East"
		}
		message := exitingPlayer.Name + " left to the " + dString + "."
		p.Notify(message)
	}
}

// NotifyEnteringPlayer checks if the enteringPlayer can be seen, and sends a message to the player.
func (p *Player) NotifyEnteringPlayer(enteringPlayer *Player, d Direction) {
	if p.IsSleeping {
		return
	}

	if (!enteringPlayer.IsHidden() || p.CanSeeHidden()) &&
		(!enteringPlayer.IsInvisible() || p.CanSeeInvisible()) {
		var dString string
		switch d {
		case North:
			dString = "South"
		case South:
			dString = "North"
		case West:
			dString = "East"
		case East:
			dString = "West"
		}
		message := enteringPlayer.Name + " came from the " + dString + "."
		p.Notify(message)
	}
}

// Say prints a message for all players on the same room
func (p *Player) Say(message string) {
	if p.IsSleeping {
		p.Notify("Is hard to talk while sleeping.")
		return
	}
	p.CurrentRoom.Say(p.UserID, p.Name, message, p.IsHidden(), p.IsInvisible())
	p.Notify("You said: " + message)
}

// Shout prints a message for all players on the same area
func (p *Player) Shout(message string) {
	if p.IsSleeping {
		p.Notify("No matter how loud you shout. Nobody can hear you in your dreams.")
		return
	}
	p.CurrentRoom.Shout(p.UserID, p.Name, message, p.IsHidden(), p.IsInvisible())
	p.Notify("You shouted: " + message)
}

// Hear prints a message from another user that can be heard
func (p *Player) Hear(playerName, message string, isHidden, isInvisible bool) {
	if p.IsSleeping {
		return
	}

	showName := playerName
	if (isHidden && !p.CanSeeHidden()) ||
		(isInvisible && !p.CanSeeInvisible()) {
		showName = "Someone"
	}
	p.Notify(fmt.Sprintf("%s says: %s", showName, message))
}

// Sleep puts the player to sleep.
func (p *Player) Sleep() {
	if p.IsSleeping {
		p.Notify("You are already deep asleep.")
		return
	}

	p.IsSleeping = true
	p.Notify("You lay down and start to sleep.")
}

// Wake wakes up the player.
func (p *Player) Wake() {
	if !p.IsSleeping {
		p.Notify("You are already awake.")
		return
	}

	p.IsSleeping = false
	p.Notify("You wake up and stand up.")
}
