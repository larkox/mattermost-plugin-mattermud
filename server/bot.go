package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-plugin-mattermud/server/mud"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func getIngameHelp() string {
	return `Ingame commands:
	n, s, e, w, north, south, east, west: Movement commands
	look: Show again the description of the room, with extra information
	status: Shows your current HP
	kill [mob]: Starts attacking the mob with that name. Example: kill bunny
	sleep: Starts to sleep. This will silence almost all notifications from the game
	wake: You wake up
	say [something you want to say]: Says something so all players in the same room will see it. Example: say Hello everyone!
	shout [something you want to shout]: Shouts something so all players in the same area will see it. Example: shout Hello everyone!
	help: Shows the ingame help`
}

// MessageHasBeenPosted checks if the message is a DM from an user, and process the message as a command in the game
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	if p.botUserID == post.UserId {
		return
	}

	ch, appErr := p.API.GetDirectChannel(p.botUserID, post.UserId)
	if appErr != nil {
		p.API.LogError("error getting direct channel: " + appErr.Error())
		return
	}

	if ch.Id != post.ChannelId {
		return
	}

	player, err := p.world.GetPlayer(post.UserId)
	if err != nil {
		p.API.LogError("user not initiated: " + err.Error())
		return
	}
	if player == nil {
		p.API.LogError("player not initiated: " + err.Error())
		return
	}

	args := strings.Split(post.Message, " ")

	switch strings.ToLower(args[0]) {
	case "n":
		p.handleMove(player, mud.North)
	case "north":
		p.handleMove(player, mud.North)
	case "s":
		p.handleMove(player, mud.South)
	case "south":
		p.handleMove(player, mud.South)
	case "e":
		p.handleMove(player, mud.East)
	case "east":
		p.handleMove(player, mud.East)
	case "w":
		p.handleMove(player, mud.West)
	case "west":
		p.handleMove(player, mud.West)
	case "look":
		p.handleLook(player)
	case "sleep":
		p.handleSleep(player)
	case "wake":
		p.handleWake(player)
	case "say":
		p.handleSay(player, args[1:])
	case "shout":
		p.handleShout(player, args[1:])
	case "kill":
		p.handleKill(player, args[1:])
	case "status":
		p.handleStatus(player)
	case "help":
		p.handleHelp(player)
	default:
		p.handleDefault(player)
	}
}

func (p *Plugin) handleMove(player *mud.Player, d mud.Direction) {
	player.Move(d)
}

func (p *Plugin) handleLook(player *mud.Player) {
	player.LookRoom()
}

func (p *Plugin) handleSleep(player *mud.Player) {
	player.Sleep()
}

func (p *Plugin) handleWake(player *mud.Player) {
	player.Wake()
}

func (p *Plugin) handleSay(player *mud.Player, args []string) {
	message := strings.Join(args, " ")
	player.Say(message)
}

func (p *Plugin) handleShout(player *mud.Player, args []string) {
	message := strings.Join(args, " ")
	player.Shout(message)
}

func (p *Plugin) handleKill(player *mud.Player, args []string) {
	objective := strings.Join(args, " ")
	player.Kill(objective)
}

func (p *Plugin) handleStatus(player *mud.Player) {
	player.Notify(fmt.Sprintf("%d/%d HP", player.CurrentHP, player.MaxHP))
}

func (p *Plugin) handleHelp(player *mud.Player) {
	player.Notify(getIngameHelp())
}

func (p *Plugin) handleDefault(player *mud.Player) {
	player.Notify("I do not understand what you say. Type `help` if you want to check all the available commands.")
}

func (p *Plugin) welcome(userID string) {
	player, err := p.world.GetPlayer(userID)
	if err != nil {
		return
	}
	p.world.Notify(userID, "Welcome to MatterMUD")
	player.ShowRoom()
}
