package mud

import (
	"fmt"
	"time"
)

// MobList represents a list of enemies
type MobList []*Mob

// Mob represents one single enemy
type Mob struct {
	// ID represents the type of monster
	ID string
	// Name represents the name shown to the player
	Name string
	// Stats are the stats of the mob
	Stats Stats
	// MaxHP denotes the Maximum Health points
	MaxHP int
	// CurrentHP denotes the current Health points
	CurrentHP int
	// Experience how many experience points the mob provides
	Experience int
	// Effects show all the magical effects that the mob is currently under
	Effects EffectList
	// Drops contains all the items dropped by the mob
	Drops []*Drop
	// DeadAt tells when the monster was defeated
	DeadAt time.Time
}

// Drop represents a drop from a monster with the probability to drop
type Drop struct {
	// Item is the item to drop
	Item Item
	// Probability is the chance to get the item as x out of 10000
	Probability int
}

// Spawn creates a new mob using another mob as template
func (m *Mob) Spawn() *Mob {
	newMob := *m
	newMob.Effects = EffectList{}
	newMob.CurrentHP = newMob.MaxHP
	return &newMob
}

// Show returns the string of how the user is seen
func (m *Mob) Show(canSeeHidden, canSeeInvisible bool) string {
	if (!canSeeHidden && m.IsHidden()) ||
		(!canSeeInvisible && m.IsInvisible()) {
		return ""
	}

	return fmt.Sprintf("A %s is here.", m.Name)
}

// IsHidden returns whether the character is hidden
func (m *Mob) IsHidden() bool {
	return m.Effects.GrantHidden()
}

// IsInvisible returns whether the character is invisible
func (m *Mob) IsInvisible() bool {
	return m.Effects.GrantInvisible()
}
