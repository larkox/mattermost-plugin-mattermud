package mud

// EquipmentSlot represents each slot of equipment
type EquipmentSlot int

// PlayerEquipment stores all the equipped items from a player
type PlayerEquipment map[EquipmentSlot]*Equipment

const (
	// Head represents any item that can be wear on the head, like helmets, masks or caps
	Head EquipmentSlot = iota
	// Chest represents any item that can be wear on the chest, like jackets, chest guards...
	Chest
	// Legs represent any item that can be wear on the legs, like grieves or pants
	Legs
	// Feet represent any item that can be wear on the feet, like boots or sandals
	Feet
	// RightHand represent any item that can be wielded while wielded on the right hand. It is used also as default for items that can be wielded
	RightHand
	// LeftHand represent any item that can be wielded while wielded on the left hand
	LeftHand
	// Necklace represent any item that hang from the neck like necklaces
	Necklace
	// RightRing represent any item that can be wear as a ring while wear on the right hand. It is used also as default for items that can be wear as rings
	RightRing
	// LeftRing represent any item that can be wear as a ring while wear on the left hand.
	LeftRing
)

// GetRightAttack gets the attack of your right hand weapon
func (e PlayerEquipment) GetRightAttack() int {
	return e[RightHand].GetAttack()
}

// GetLeftAttack gets the attack of your left hand weapon
func (e PlayerEquipment) GetLeftAttack() int {
	return e[LeftHand].GetAttack()
}

// GetAttackModifiers gets all attack modifiers from the equipment
func (e PlayerEquipment) GetAttackModifiers() int {
	modifier := 0
	for k, v := range e {
		if k == RightHand || k == LeftHand {
			continue
		}
		modifier += v.GetAttack()
	}
	return modifier
}

// GetStatModifiers get all the modifiers for certain stat from the equipment
func (e PlayerEquipment) GetStatModifiers(s Stat) int {
	modifier := 0
	for _, v := range e {
		modifier += v.GetStat(s)
	}
	return modifier
}

// CanSeeInvisible returns whether any piece of equipment lets you see the invisible
func (e PlayerEquipment) CanSeeInvisible() bool {
	for _, v := range e {
		if v.CanSeeInvisible() {
			return true
		}
	}
	return false
}

// CanSeeHidden returns whether any piece of equipment lets you see hidden objects
func (e PlayerEquipment) CanSeeHidden() bool {
	for _, v := range e {
		if v.CanSeeHidden() {
			return true
		}
	}
	return false
}

// GrantInvisible returns whether any piece of equipment renders you invisible
func (e PlayerEquipment) GrantInvisible() bool {
	for _, v := range e {
		if v.GrantInvisible() {
			return true
		}
	}
	return false
}

// GrantHidden returns whether any piece of equipment renders you hidden
func (e PlayerEquipment) GrantHidden() bool {
	for _, v := range e {
		if v.GrantHidden() {
			return true
		}
	}
	return false
}

// Equipment represents a single piece of equipment
type Equipment struct {
	Item
	// Slot denotes where it is wear. For wielded items will always denote RightHand even if it can be wielded on both hands. Same for rings.
	Slot EquipmentSlot
	// StatsModifiers denotes how much modify each stat
	StatsModifiers Stats
	// Attack denotes how much attack it grants
	Attack int
	// MagicEffects denotes all the magical effects this item has
	MagicEffects EffectList
}

// GetAttack returns the attack of the item
func (e *Equipment) GetAttack() int {
	if e == nil {
		return 0
	}

	return e.Attack
}

// GetStat returns the stat modifier of the item
func (e *Equipment) GetStat(s Stat) int {
	if e == nil {
		return 0
	}

	return e.StatsModifiers[s]
}

// CanSeeInvisible returns whether this equip has any magical effect that let you see the invisible
func (e *Equipment) CanSeeInvisible() bool {
	if e == nil {
		return false
	}

	return e.MagicEffects.CanSeeInvisible()
}

// CanSeeHidden returns whether this equip has any magical effect that let you see hidden things
func (e *Equipment) CanSeeHidden() bool {
	if e == nil {
		return false
	}

	return e.MagicEffects.CanSeeHidden()
}

// GrantInvisible returns whether this equip renders you invisible
func (e *Equipment) GrantInvisible() bool {
	if e == nil {
		return false
	}

	return e.MagicEffects.GrantInvisible()
}

// GrantHidden returns whether this equip renders you hidden
func (e *Equipment) GrantHidden() bool {
	if e == nil {
		return false
	}

	return e.MagicEffects.GrantHidden()
}
