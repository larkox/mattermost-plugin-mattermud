package mud

import (
	"fmt"
	"time"
)

const (
	//MobRegenTime marks how long the routine sleep between regens
	MobRegenTime = 1 * time.Minute
	//MobSpawnTime marks how long does it take for a mob to reswpawn
	MobSpawnTime = 5 * time.Minute
)

func (m *Mob) finishMobRoutine() bool {
	select {
	case <-worldShutDown:
		return true
	default:
		return false
	}
}

// MobList represents a list of enemies
type MobList []*Mob

// Mob represents one single enemy
type Mob struct {
	// ID represents the type of monster. It is also the name shown to the player
	ID string
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
	newMob.start()
	return &newMob
}

// Show returns the string of how the user is seen
func (m *Mob) Show(canSeeHidden, canSeeInvisible bool) string {
	if (!canSeeHidden && m.IsHidden()) ||
		(!canSeeInvisible && m.IsInvisible()) {
		return ""
	}

	return fmt.Sprintf("A %s is here.", m.ID)
}

// IsHidden returns whether the character is hidden
func (m *Mob) IsHidden() bool {
	return m.Effects.GrantHidden()
}

// IsInvisible returns whether the character is invisible
func (m *Mob) IsInvisible() bool {
	return m.Effects.GrantInvisible()
}

// GetAttack returns the attack of the mob
func (m *Mob) GetAttack() int {
	str := m.GetCurrentStat(Strength)
	attEffectModifiers := m.Effects.GetAttackModifiers()
	return min(0, str+attEffectModifiers)
}

// GetCurrentStat returns the current stat of the mob
func (m *Mob) GetCurrentStat(s Stat) int {
	base := m.Stats[s]
	effectModifiers := m.Effects.GetStatModifiers(s)
	return min(0, base+effectModifiers)
}

// GetCurrentDefense returns the current defense of the mob
func (m *Mob) GetCurrentDefense() int {
	return m.GetCurrentStat(Constitution)
}

// Dead kills the mob
func (m *Mob) Dead() {
	m.DeadAt = time.Now()
}

// start runs the mob routine
func (m *Mob) start() {
	go func() {
		for {
			if m.finishMobRoutine() {
				return
			}
			if m.CurrentHP <= 0 && time.Now().Unix() > m.DeadAt.Add(MobSpawnTime).Unix() {
				m.CurrentHP = m.MaxHP
			}

			if m.CurrentHP > 0 {
				toRegen := max(1, int(float64(m.MaxHP)*0.1))
				m.CurrentHP = min(m.MaxHP, m.CurrentHP+toRegen)
			}
			time.Sleep(MobRegenTime)
		}
	}()
}
