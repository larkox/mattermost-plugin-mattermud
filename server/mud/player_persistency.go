package mud

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

const (
	// ASSleepTime defines how long should the Auto Save sleep between one save and another
	ASSleepTime = 30 * time.Minute
)

func finishAutosave() bool {
	select {
	case <-worldShutDown:
		return true
	default:
		return false
	}
}

// JSONPlayer represent a player as stored in the persistant store
type JSONPlayer struct {
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
	// CurrentRoom shows the id of the room on which the player is currently on
	CurrentRoom string
	// MaxHP denotes the Maximum Health points
	MaxHP int
	// CurrentHP denotes the current Health points
	CurrentHP int
}

// autoSave stores the player information periodically into the persistant memory
func (w *World) autoSave() {
	for {
		time.Sleep(ASSleepTime)
		if finishAutosave() {
			return
		}

		w.SavePlayers()
	}
}

// SavePlayers store the list of players on the persistant memory
func (w *World) SavePlayers() error {
	jsonPlayers := make([]*JSONPlayer, 0, len(w.players))
	for _, v := range w.players {
		jsonPlayers = append(jsonPlayers, playerToJSONPlayer(v))
	}
	marshalledPlayers, jsonErr := json.Marshal(jsonPlayers)

	if jsonErr != nil {
		return jsonErr
	}

	appErr := w.api.KVSet(playerListKey(), marshalledPlayers)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

// GetPlayers get all the players from the persistant memory and loads them into the world
func (w *World) GetPlayers() error {
	marshalledPlayers, appErr := w.api.KVGet(playerListKey())
	if appErr != nil {
		return errors.New(appErr.Error())
	}
	w.api.LogDebug(string(marshalledPlayers))

	var jsonPlayers []*JSONPlayer
	jsonErr := json.Unmarshal(marshalledPlayers, &jsonPlayers)
	if jsonErr != nil {
		return jsonErr
	}

	for _, v := range jsonPlayers {
		w.players[v.UserID] = w.jsonPlayerToPlayer(v)
	}

	return nil
}

func playerListKey() string {
	return "players"
}
