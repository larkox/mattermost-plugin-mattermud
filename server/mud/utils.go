package mud

import (
	"fmt"
	"strings"
)

// Min returns the minimum between two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// jsonRoomsToRooms convert imported json rooms to usable Rooms in the game
func jsonRoomsToRooms(in map[string]*JSONRoom) (map[string]*Room, error) {
	out := make(map[string]*Room)

	for k, v := range in {
		out[k] = &Room{
			ID:               v.ID,
			Name:             v.Name,
			AreaID:           v.AreaID,
			ShortDescription: v.ShortDescription,
			LongDescription:  v.LongDescription,
			Mobs:             MobList{},
			Players:          []*Player{},
			Neighbours:       make(map[Direction]*RoomDoor),
		}
	}

	for id, room := range in {
		for direction, door := range room.Neighbours {
			var directionKey Direction
			switch direction {
			case "north":
				directionKey = North
			case "south":
				directionKey = South
			case "east":
				directionKey = East
			case "west":
				directionKey = West
			}

			roomID := room.AreaID + "_" + door.Room
			if strings.HasPrefix(door.Room, "__EXT__") {
				roomID = door.Room[7:]
			}

			neighbourRoom, ok := out[roomID]
			if !ok {
				return nil, fmt.Errorf("Cannot find neighbour with id %s for room %s", roomID, id)
			}

			out[id].Neighbours[directionKey] = &RoomDoor{
				isHidden:    door.IsHidden,
				isInvisible: door.IsInvisible,
				isLocked:    door.IsLocked,
				room:        neighbourRoom,
			}
		}
	}

	return out, nil
}
