package cmd

import (
	"github.com/Bo0km4n/claude/pkg/proxy/config"
	"github.com/Bo0km4n/claude/pkg/proxy/repository"
	"github.com/Bo0km4n/claude/pkg/proxy/service"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	TabletIP       string
	TabletPort     string
	IsUDPMulticast bool
}

var ro = &RunOptions{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Proxy server",
	PreRun: func(cmd *cobra.Command, args []string) {
		config.InitConfig()

		// config.Config's value should not be change from outer.
		// But this application need using sudo permission,
		// and then oprating is complication keeping envrionment variables.
		// So temporary, oprating config's value from outer.
		// I'm going to resolve this problem some day.
		// config.Config.Tablet.IP = ro.TabletIP
		// config.Config.Tablet.Port = ro.TabletPort
		repository.InitDB()
	},
	Run: func(cmd *cobra.Command, args []string) {
		run(ro)
	},
}

func run(ro *RunOptions) {
	service.LaunchService(ro.IsUDPMulticast)
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&ro.TabletIP, "tablet_ip", "", ro.TabletIP, "tablet ip")
	runCmd.Flags().StringVarP(&ro.TabletPort, "tablet_port", "", ro.TabletPort, "tablet port")
	runCmd.Flags().BoolVarP(&ro.IsUDPMulticast, "is_udp_multicast", "", false, "use udp multicast")
}
