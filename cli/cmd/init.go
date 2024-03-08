package cmd

import (
	"fmt"
	"github.com/sevenreup/goact/cli/actions"
	"github.com/spf13/cobra"
	"os"
)

type BootStrapConfig struct {
}

var packageManager string

var rootCmd = &cobra.Command{Use: "goact"}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Goact Project Init",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		action := actions.InitAction{
			PackageManager: packageManager,
		}
		action.HandleAction()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&packageManager, "packageManger", "npm", "config file (default is $HOME/.cobra.yaml)")
}

func Execute() {
	rootCmd.AddCommand(initCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
