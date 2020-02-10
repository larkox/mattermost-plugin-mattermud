package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "mattermud",
		DisplayName:      "Mattermud",
		Description:      "Create a new mattermud player.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: ",
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
	err := p.world.NewPlayer(args.UserId)
	if err != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "There has been an error creating your player: "+err.Error()), nil
	}
	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Welcome to mattermud. The GM just messaged you to start the game."), nil
}
