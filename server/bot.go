package main

import (
	"strings"

	"github.com/mattermost/mattermost-plugin-mattermud/server/mud"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

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

func (p *Plugin) welcome(userID string) {
	player, err := p.world.GetPlayer(userID)
	if err != nil {
		return
	}
	p.world.Notify(userID, "Welcome to MatterMUD")
	player.ShowRoom()
}
