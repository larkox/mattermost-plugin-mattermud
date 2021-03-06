package main

import (
	"sync"

	"github.com/mattermost/mattermost-plugin-mattermud/server/mud"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	botUsername    = "mattermudgm"
	botDisplayName = "Mattermud GM"
	botDescription = "The game master of Mattermud."
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	// botUserID of the created bot account.
	botUserID string

	world mud.World
}

// OnActivate handles all initialization
func (p *Plugin) OnActivate() error {
	bot := &model.Bot{
		Username:    botUsername,
		DisplayName: botDisplayName,
		Description: botDescription,
	}
	botUserID, appErr := p.Helpers.EnsureBot(bot)
	if appErr != nil {
		return errors.Wrap(appErr, "failed to ensure bot user")
	}
	p.botUserID = botUserID

	p.world = mud.NewWorld(p.API, botUserID)
	err := p.world.Init()
	if err != nil {
		return errors.Wrap(err, "failed to init the world")
	}

	return p.API.RegisterCommand(getCommand())
}

// OnDeactivate handles all finalization
func (p *Plugin) OnDeactivate() error {
	p.world.Finalize()
	return nil
}
