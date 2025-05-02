package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/cmdenv"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/secretkey"
)

func vaultAccessCommand() *cobra.Command {
	if vaultAccessCmd == nil {
		vaultAccessCmd = &cobra.Command{
			Use:     "access",
			Aliases: []string{"rights", "privilege", "permission", "permissions"},
			Short:   "Managing access to a vault",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		vaultAccessCmd.PersistentFlags().StringSliceP(cmdenv.EnvPublicKeysFlag.Name, cmdenv.EnvPublicKeysFlag.Shorthand, []string{}, cmdenv.EnvPublicKeysFlag.Usage)
		vaultAccessCmd.PersistentFlags().StringSliceP(cmdenv.EnvSearchFlag.Name, cmdenv.EnvSearchFlag.Shorthand, []string{}, cmdenv.EnvSearchFlag.Usage)
		vaultAccessCmd.PersistentFlags().BoolP(cmdenv.EnvSelfFlag.Name, cmdenv.EnvSelfFlag.Shorthand, false, cmdenv.EnvSelfFlag.Usage)
		vaultAccessCmd.PersistentFlags().BoolP(cmdenv.EnvK8sFlag.Name, cmdenv.EnvK8sFlag.Shorthand, false, cmdenv.EnvK8sFlag.Usage)
		vaultAccessCmd.PersistentFlags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage+" (used with k8s environment)")
		vaultAccessCmd.AddCommand(vaultAccessAddCommand())
		vaultAccessCmd.AddCommand(vaultAccessRemoveCommand())
	}
	return vaultAccessCmd
}

func vaultAccessAddCommand() *cobra.Command {
	if vaultAccessAddCmd == nil {
		vaultAccessAddCmd = &cobra.Command{
			Use:     "add",
			Aliases: []string{"allow", "grant", "share"},
			Short:   "Adds access to a vault for the given environments/public keys",
			Run: func(cmd *cobra.Command, args []string) {
				envSecretKey, err := secretkey.Get()
				if err != nil {
					utils.ExitOnError(err)
				}
				vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
				k8sPQ, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
				publicKeys, err := cmdenv.GetPublicKeys(cmd, true, k8sPQ)
				if err != nil {
					utils.ExitOnError(err)
				}
				vault, err := getVault(vaultFile)
				if err == nil {
					err = vault.Unlock(envSecretKey)
					if err == nil {
						for _, publicKey := range publicKeys {
							if _, err = vault.Share(publicKey); err != nil {
								break
							}
						}
						if err == nil {
							fmt.Println("Shared vault:", color.GreenString(vaultFile))
							utils.SafeExit()
						}
					}
				}
				utils.ExitOnError(err)
			},
		}
	}
	return vaultAccessAddCmd
}

func vaultAccessRemoveCommand() *cobra.Command {
	if vaultAccessRemoveCmd == nil {
		vaultAccessRemoveCmd = &cobra.Command{
			Use:     "remove",
			Aliases: []string{"rm", "deny", "revoke", "restrict", "delete", "del"},
			Short:   "Remove access to a vault for the given environments/public keys",
			Run: func(cmd *cobra.Command, args []string) {
				vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
				k8sPQ, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
				publicKeys, err := cmdenv.GetPublicKeys(cmd, false, k8sPQ)
				if err != nil {
					utils.ExitOnError(err)
				}
				vault, err := getVault(vaultFile)
				if err == nil {
					var envSecretKey *crypto.SecretKey
					if envSecretKey, err = secretkey.Get(); err == nil {
						err = vault.Unlock(envSecretKey)
					}
					if err == nil {
						pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
						if err = vault.Revoke(publicKeys, pq); err == nil {
							fmt.Println("Shared vault:", color.GreenString(vaultFile))
							utils.SafeExit()
						}
					}
				}
				utils.ExitOnError(err)
			},
		}
	}
	return vaultAccessRemoveCmd
}
