# MUD client and server for Mattermost

Mattermud is the Multi-user dungeon integrated in Mattermost. The commands available are:
	start: Creates a player for you and starts the game
	help: Shows this help text

Ingame commands:
	n, s, e, w, north, south, east, west: Movement commands
	look: Show again the description of the room, with extra information
	status: Shows your current HP
	kill [mob]: Starts attacking the mob with that name. Example: kill bunny
	sleep: Starts to sleep. This will silence almost all notifications from the game
	wake: You wake up
	say [something you want to say]: Says something so all players in the same room will see it. Example: say Hello everyone!
	shout [something you want to shout]: Shouts something so all players in the same area will see it. Example: shout Hello everyone!
	help: Shows the ingame help