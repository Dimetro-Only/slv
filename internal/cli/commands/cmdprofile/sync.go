package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/profiles"
)

func profilePullCommand() *cobra.Command {
	if profilePullCmd == nil {
		profilePullCmd = &cobra.Command{
			Use:     "pull",
			Aliases: []string{"sync"},
			Short:   "Pulls the latest changes for the current profile from remote repository",
			Run: func(cmd *cobra.Command, args []string) {
				profile, err := profiles.GetCurrentProfile()
				if err != nil {
					utils.ExitOnError(err)
				}
				if err = profile.Pull(); err != nil {
					utils.ExitOnError(err)
				}
				fmt.Printf("Successfully pulled changes into profile: %s\n", color.GreenString(profile.Name()))
			},
		}
	}
	return profilePullCmd
}

func profilePushCommand() *cobra.Command {
	if profilePushCmd == nil {
		profilePushCmd = &cobra.Command{
			Use:   "push",
			Short: "Pushes the changes in the current profile to the pre-configured remote repository",
			Run: func(cmd *cobra.Command, args []string) {
				profile, err := profiles.GetCurrentProfile()
				if err != nil {
					utils.ExitOnError(err)
				}
				if err = profile.Push(); err != nil {
					utils.ExitOnError(err)
				}
				fmt.Printf("Successfully pushed changes from profile: %s\n", color.GreenString(profile.Name()))
			},
		}
	}
	return profilePushCmd
}
