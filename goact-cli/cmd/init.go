package cmd

import (
	"errors"
	"fmt"
	"github.com/sevenreup/goact/goact-cli/actions"
	"github.com/spf13/cobra"
	"os"
)

type BootStrapConfig struct {
}

var packageManager string
var tailwindFlag bool
var viewDir string

var rootCmd = &cobra.Command{Use: "goact"}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Goact Project Init",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if tailwindFlag && viewDir == "" {
			return errors.New("viewDir is required when using tailwind flag")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		action := actions.InitAction{
			PackageManager: packageManager,
			UseTailwind:    tailwindFlag,
			ViewDir:        viewDir,
		}
		action.HandleAction()
	},
}

func init() {
	initCmd.PersistentFlags().StringVarP(&packageManager, "packageManger", "p", "npm", "The package manager to use")
	initCmd.Flags().BoolVarP(&tailwindFlag, "tailwind", "t", false, "Include Tailwind")
	initCmd.Flags().StringVarP(&viewDir, "viewDir", "v", "", "Directory location for views (required with tailwind)")
}

func Execute() {
	rootCmd.AddCommand(initCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
