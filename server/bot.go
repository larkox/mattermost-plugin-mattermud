package main

import (
	"fmt"
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
	}
}

func (p *Plugin) handleMove(player *mud.Player, d mud.Direction) {
	p.postBotDM(player.UserID, player.Move(d))
}

func (p *Plugin) handleLook(player *mud.Player) {
	p.postBotDM(player.UserID, player.Look())
}

func (p *Plugin) postBotDM(userID string, message string) error {
	channel, appError := p.API.GetDirectChannel(userID, p.botUserID)
	if appError != nil {
		return appError
	}
	if channel == nil {
		return fmt.Errorf("could not get direct channel for bot and user_id=%s", userID)
	}

	_, appError = p.API.CreatePost(&model.Post{
		UserId:    p.botUserID,
		ChannelId: channel.Id,
		Message:   message,
	})

	return appError
}

func (p *Plugin) welcome(userID string) {
	p.postBotDM(userID, "Welcome to MatterMUD")
	p.postBotDM(userID, p.world.GetUser(userID).GetRoom())
}
