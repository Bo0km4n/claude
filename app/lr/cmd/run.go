package cmd

import (
	"github.com/Bo0km4n/claude/app/lr/config"
	"github.com/Bo0km4n/claude/app/lr/repository"
	"github.com/Bo0km4n/claude/app/lr/service"
	"github.com/spf13/cobra"
)

type runOptions struct {
	TabletIP   string
	TabletPort string
}

var ro = &runOptions{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run LR server",
	PreRun: func(cmd *cobra.Command, args []string) {
		config.InitConfig()

		// config.Config's value should not be change from outer.
		// But this application need using sudo permission,
		// and then oprating is complication keeping envrionment variables.
		// So temporary, oprating config's value from outer.
		// I'm going to resolve this problem some day.
		config.Config.Tablet.IP = ro.TabletIP
		config.Config.Tablet.Port = ro.TabletPort
		repository.InitDB()
	},
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	service.LaunchService()
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&ro.TabletIP, "tablet_ip", "", ro.TabletIP, "tablet ip")
	runCmd.Flags().StringVarP(&ro.TabletPort, "tablet_port", "", ro.TabletPort, "tablet port")
}
