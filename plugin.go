package main

import (
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const pluginCommand = "antimatter"

var postGrabSize = 25

// AntimatterPlugin implements the interface expected by the Mattermost server to communicate
// between the server and plugin processes.
type AntimatterPlugin struct {
	plugin.MattermostPlugin
}

// Creates a command to test plugin slash commands
func createPluginCommand() *model.Command {
	return &model.Command{
		Trigger: pluginCommand,
	}
}

// OnActivate contains code to register the plugin command with the API
func (p *AntimatterPlugin) OnActivate() error {
	if err := p.API.RegisterCommand(createPluginCommand()); err != nil {
		return errors.Wrapf(err, "failed to register %s command", pluginCommand)
	}
	return nil
}

// ExecuteCommand executes the plugin's slash command
func (p *AntimatterPlugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	input := strings.Split(strings.TrimPrefix(args.Command, "/"+pluginCommand+" "), " ")

	response := "Unknown arguments, type /antimatter help for help."
	var errOut *model.AppError

	switch input[0] {
	case "help":
		response = "/antimatter [mode]\nModes:\n\tclear - deletes all posts in channel"
	case "clear":
		if !p.API.HasPermissionToChannel(args.UserId, args.ChannelId, model.PERMISSION_DELETE_OTHERS_POSTS) {
			response = "You do not have the permissions necessary to run this command."
		}
		numDeleted := 0
		for postsLeft := true; postsLeft; {
			postlist, err := p.API.GetPostsForChannel(args.ChannelId, 0, postGrabSize)
			if err != nil {
				errOut = err
				break
			}
			posts := postlist.ToSlice()
			if len(posts) > 0 {
				for i := 0; i < len(posts); i++ {
					err := p.API.DeletePost(posts[i].Id)
					if err != nil {
						errOut = err
						break
					}
					numDeleted++
				}
			} else {
				postsLeft = false
			}
		}
		if errOut != nil {
			break
		}
		response = "Successfully deleted " + strconv.Itoa(numDeleted) + " posts."
	}

	if errOut != nil {
		response = "Error: see System Log for details."
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         response,
	}, errOut
}

// This example demonstrates a plugin that handles HTTP requests which respond by greeting the
// world.
func main() {
	plugin.ClientMain(&AntimatterPlugin{})
}
