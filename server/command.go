package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func getHelp() string {
	return `Mattermud is the Multi-user dungeon integrated in Mattermost. The commands available are:
	start: Creates a player for you and starts the game
	help: Shows this help text

` + getIngameHelp()
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "mattermud",
		DisplayName:      "Mattermud",
		Description:      "Create a new mattermud player.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: help, start",
		AutoCompleteHint: "[command]",
	}
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     "mattermud",
		//IconURL:      fmt.Sprintf("/plugins/%s/profile.png", manifest.ID),
	}
}

// ExecuteCommand executes any command to mattermud
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	stringArgs := strings.Split(strings.TrimSpace(args.Command), " ")
	lengthOfArgs := len(stringArgs)

	if lengthOfArgs == 1 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), nil
	}

	command := stringArgs[1]

	switch command {
	case "help":
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), nil
	case "start":
		err := p.world.NewPlayer(args.UserId)
		if err != nil {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "There has been an error creating your player: "+err.Error()), nil
		}
		p.welcome(args.UserId)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Welcome to mattermud. The GM just messaged you to start the game."), nil
	default:
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), nil
	}
}
