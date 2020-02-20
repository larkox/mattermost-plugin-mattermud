package mud

import "time"

const (
	// ShoutLifespan defines how long a shout id is stored in a room before the garbage collector removes it
	ShoutLifespan = 10 * time.Second
	// GCSleepTime defines how long should the GC sleep between one room and another
	GCSleepTime = 1 * time.Second
)

var garbageDone = make(chan struct{})

func finishGarbageCollection() bool {
	select {
	case <-garbageDone:
		return true
	default:
		return false
	}
}

func (w *World) garbageCollector() {
	for {
		for _, room := range w.rooms {
			time.Sleep(GCSleepTime)
			if finishGarbageCollection() {
				return
			}
			t := time.Now()
			for k, v := range room.shouts {
				if t.Unix() < v.Add(ShoutLifespan).Unix() {
					delete(room.shouts, k)
					w.api.LogDebug("Shout deleted.")
				}
			}

			// for _, p := range room.Players {
			// 	p.Notify("An aura of cleanliness just passed through this room.")
			// }
		}
	}
}
