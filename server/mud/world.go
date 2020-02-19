package mud

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	// startingRoom is the room where all players start to play
	startingRoom = "midgaard_temple"
)

// World stores all the information from the game
type World struct {
	api       plugin.API
	botUserID string
	rooms     map[string]*Room
	players   map[string]*Player
}

// JSONArea is the struct of area files of mattermud
type JSONArea struct {
	ID    string      `json:"area_id"`
	Rooms []*JSONRoom `json:"rooms"`
}

// JSONRoom is the struct of rooms on area files of mattermud
type JSONRoom struct {
	// AreaID is the unique ID for the area. Not present on each room on the json file
	AreaID string
	// ID is the unique ID of the room inside the area. Final ID will be AreaID + _ + ID
	ID string `json:"id"`
	// Name is the name shown to the player
	Name string `json:"name"`
	// ShortDescription is the name shown to the player
	ShortDescription string `json:"short_description"`
	// LongDescription is the name shown to the player when using the command look
	LongDescription string `json:"long_description"`
	// Mobs is the list of IDs of Mobs in the area
	Mobs []string `json:"mobs"`
	// Neighbours is map of rooms neighbour to this one
	Neighbours map[string]JSONNeighbour `json:"neighbours"`
}

// JSONNeighbour is the struct for room transitions on area files of mattermud
type JSONNeighbour struct {
	// IsHidden shows whether the transition is hidden
	IsHidden bool `json:"is_hidden"`
	// IsInvisible shows whether the transition is invisible
	IsInvisible bool `json:"is_invisible"`
	// IsLocked shows whether the transition is locked behind a door
	IsLocked bool `json:"is_locked"`
	// Room is the ID of the room this transition connects to. External transitions will have the __EXT__ prefix
	Room string `json:"id"`
	// KeyID is the key needed to open the door
	KeyID string `json:"key_id"`
}

// NewWorld creates a new world
func NewWorld(api plugin.API, botUserID string) World {
	return World{
		api:       api,
		botUserID: botUserID,
	}
}

// Init initializes the world
func (w *World) Init() error {
	bundlePath, err := w.api.GetBundlePath()
	if err != nil {
		return errors.Wrap(err, "couldn't get bundle path")
	}

	areasPath := filepath.Join(bundlePath, "assets", "areas")
	jsonRooms := make(map[string]*JSONRoom)
	err = filepath.Walk(areasPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		var area JSONArea
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&area); err != nil {
			return nil
		}

		for _, v := range area.Rooms {
			v.ID = area.ID + "_" + v.ID
			v.AreaID = area.ID
			jsonRooms[v.ID] = v
		}

		return nil
	})

	if err != nil {
		return errors.WithMessage(err, "OnActivate/loadWorld failed")
	}
	w.rooms, err = jsonRoomsToRooms(jsonRooms)
	if err != nil {
		return err
	}

	w.players = make(map[string]*Player)
	go w.garbageCollector()
	return nil
}

// NewPlayer creates a new player for userID and place it on the starting room
func (w *World) NewPlayer(userID string) error {
	if _, ok := w.players[userID]; ok {
		return errors.New("You already have a character in mattermud")
	}
	w.players[userID] = &Player{
		UserID:      userID,
		Name:        "Placeholder",
		CurrentRoom: w.rooms[startingRoom],
		Notify: func(message string) {
			w.Notify(userID, message)
		},
	}

	return nil
}

// GetPlayer returns a player from the player list
func (w *World) GetPlayer(userID string) (*Player, error) {
	return w.players[userID], nil
}

func (w *World) String() string {
	out := ""
	for id, room := range w.rooms {
		out += fmt.Sprintf("ID: %s\nArea: %s\nName: %s\nShort Description: %s\nLong Description: %s\n",
			id,
			room.AreaID,
			room.Name,
			room.ShortDescription,
			room.LongDescription)
		for d, door := range room.Neighbours {
			out += fmt.Sprintf("Door in direction %d towards %s\n", d, door.room.ID)
		}
	}
	return out
}

// Notify sends a message to the user
func (w *World) Notify(userID, message string) {
	channel, appError := w.api.GetDirectChannel(userID, w.botUserID)
	if appError != nil {
		w.api.LogError("failed to notify user, err=" + appError.Error())
		return
	}
	if channel == nil {
		w.api.LogError("failed to get direct channel")
		return
	}

	_, appError = w.api.CreatePost(&model.Post{
		UserId:    w.botUserID,
		ChannelId: channel.Id,
		Message:   message,
	})

	if appError != nil {
		w.api.LogError("failed to notify user, err=" + appError.Error())
		return
	}
}

// Finalize handles all the important task when plugin gets disabled.
func (w *World) Finalize() {
	close(garbageDone)
}
