package mud

// EffectList list the different effects that any player, mob or item may have
type EffectList []*Effect

// Effect denotes any kind of effect that any player, mob or item may have
type Effect struct {
	Attack         int
	StatsModifiers Stats
	SeeHidden      bool
	SeeInvisible   bool
	GrantInvisible bool
	GrantHidden    bool
}

// GetAttackModifiers returns the modifiers to attack provided by all the effects on the list
func (el EffectList) GetAttackModifiers() int {
	modifier := 0
	for _, v := range el {
		modifier = modifier + v.Attack
	}
	return modifier
}

// GetStatModifiers returns the modifiers to any stat provided by all the effects on the list
func (el EffectList) GetStatModifiers(s Stat) int {
	modifier := 0
	for _, v := range el {
		modifier = modifier + v.StatsModifiers[s]
	}
	return modifier
}

// CanSeeInvisible returns whether any effect lets you see invisible things
func (el EffectList) CanSeeInvisible() bool {
	for _, v := range el {
		if v.SeeInvisible {
			return true
		}
	}
	return false
}

// CanSeeHidden returns whether any effect lets you see hidden things
func (el EffectList) CanSeeHidden() bool {
	for _, v := range el {
		if v.SeeHidden {
			return true
		}
	}
	return false
}

// GrantInvisible returns whether any effect renders you invisible
func (el EffectList) GrantInvisible() bool {
	for _, v := range el {
		if v.GrantInvisible {
			return true
		}
	}
	return false
}

// GrantHidden returns whether any effect renders you hidden
func (el EffectList) GrantHidden() bool {
	for _, v := range el {
		if v.GrantHidden {
			return true
		}
	}
	return false
}
