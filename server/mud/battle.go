package mud

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	//BattleTurnTime defines how long each fight turn last
	BattleTurnTime = 5 * time.Second
)

// Battle represents any battle between mobs and players
type Battle struct {
	PlayerSide  []*Player
	MobSide     []*Mob
	lock        sync.Mutex
	battleEnded chan struct{}
}

func (b *Battle) finishBattle() bool {
	select {
	case <-worldShutDown:
		return true
	case <-b.battleEnded:
		return true
	default:
		return false
	}
}

// Stop stops the processing of the battle
func (b *Battle) Stop() {
	if !b.finishBattle() {
		close(b.battleEnded)
	}
}

// Start starts the processing of the battle
func (b *Battle) Start() {
	go func() {
		for {
			if b.finishBattle() {
				return
			}
			hitNotifications := []string{}
			killNotifications := []string{}
			for _, p := range b.PlayerSide {
				att := p.GetLeftAttack() + p.GetRightAttack()
				mob := b.GetNextMob()
				def := mob.GetCurrentDefense()
				damage := max(att-def, 1)
				mob.CurrentHP -= damage
				hitNotifications = append(hitNotifications, fmt.Sprintf("%s inflicted %d damage to %s.", p.Name, damage, mob.ID))
				if mob.CurrentHP <= 0 {
					killNotifications = append(killNotifications, fmt.Sprintf("%s killed the %s!", p.Name, mob.ID))
				}
			}
			for _, m := range b.MobSide {
				att := m.GetAttack()
				player := b.GetNextPlayer()
				def := player.GetCurrentDefense()
				damage := max(att-def, 1)
				player.CurrentHP -= damage
				hitNotifications = append(hitNotifications, fmt.Sprintf("The %s inflicted %d damange to %s.", m.ID, damage, player.Name))
				if player.CurrentHP <= 0 {
					killNotifications = append(killNotifications, fmt.Sprintf("The %s killed %s!", m.ID, player.Name))
				}
			}
			b.NotifyAll(strings.Join(append(hitNotifications, killNotifications...), "\n"))
			playersToRemove := []*Player{}
			for _, p := range b.PlayerSide {
				if p.CurrentHP <= 0 {
					playersToRemove = append(playersToRemove, p)
				}
			}
			for _, p := range playersToRemove {
				b.RemovePlayer(p)
				p.Dead()
			}

			mobsToRemove := []*Mob{}
			for _, m := range b.MobSide {
				if m.CurrentHP <= 0 {
					mobsToRemove = append(mobsToRemove, m)
				}
			}
			for _, m := range mobsToRemove {
				b.RemoveMob(m)
				m.Dead()
			}

			if len(b.PlayerSide) == 0 {
				b.Stop()
			}

			if len(b.MobSide) == 0 {
				for _, p := range b.PlayerSide {
					p.IsFighting = false
				}
				b.NotifyAll("You won!")
				b.Stop()
			}
			time.Sleep(BattleTurnTime)
		}
	}()
}

// GetNextMob returns the next alive mob in the list
func (b *Battle) GetNextMob() *Mob {
	for _, m := range b.MobSide {
		if m.CurrentHP > 0 {
			return m
		}
	}
	return nil
}

// GetNextPlayer returns the next alive player in the list
func (b *Battle) GetNextPlayer() *Player {
	for _, p := range b.PlayerSide {
		if p.CurrentHP > 0 {
			return p
		}
	}
	return nil
}

// AddPlayer adds one player to the battle
func (b *Battle) AddPlayer(player *Player) {
	b.lock.Lock()
	defer b.lock.Unlock()

	for _, v := range b.PlayerSide {
		if v == player {
			return
		}
	}
	b.PlayerSide = append(b.PlayerSide, player)
	player.IsFighting = true
}

// RemovePlayer removes one player from the battle
func (b *Battle) RemovePlayer(player *Player) {
	for i, v := range b.PlayerSide {
		if v == player {
			b.PlayerSide = append(b.PlayerSide[:i], b.PlayerSide[i+1:]...)
			v.IsFighting = false
			return
		}
	}
	return
}

// AddMob adds one mob to the battle
func (b *Battle) AddMob(mob *Mob) {
	b.lock.Lock()
	defer b.lock.Unlock()

	for _, v := range b.MobSide {
		if v == mob {
			return
		}
	}
	b.MobSide = append(b.MobSide, mob)
}

// RemoveMob removes one mob from the battle
func (b *Battle) RemoveMob(mob *Mob) {
	for i, v := range b.MobSide {
		if v == mob {
			b.MobSide = append(b.MobSide[:i], b.MobSide[i+1:]...)
			return
		}
	}
	return
}

// mergeBattles merge two battles
func mergeBattles(a, b *Battle) *Battle {
	newBattle := &Battle{
		PlayerSide: []*Player{},
		MobSide:    []*Mob{},
	}

	if a != nil {
		a.lock.Lock()
		defer a.lock.Unlock()
		newBattle.PlayerSide = a.PlayerSide
		newBattle.MobSide = a.MobSide
		a.Stop()
	}

	if b != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	OUTER:
		for _, in := range b.PlayerSide {
			for _, present := range newBattle.PlayerSide {
				if in == present {
					continue OUTER
				}
			}
			newBattle.PlayerSide = append(newBattle.PlayerSide, in)
		}
		b.Stop()
	}

	newBattle.battleEnded = make(chan struct{})

	return newBattle
}

// IsPlayerFighting returns whether the player is fighting on this battle
func (b *Battle) IsPlayerFighting(player *Player) bool {
	for _, p := range b.PlayerSide {
		if p == player {
			return true
		}
	}
	return false
}

// IsMobFighting returns whether the mob is fighting on this battle
func (b *Battle) IsMobFighting(mob *Mob) bool {
	for _, m := range b.MobSide {
		if m == mob {
			return true
		}
	}
	return false
}

// NotifyAll sends a message to all the players on the battle
func (b *Battle) NotifyAll(message string) {
	for _, p := range b.PlayerSide {
		p.Notify(message)
	}
}
