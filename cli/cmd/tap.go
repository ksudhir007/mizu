package cmd

import (
	"errors"
	"fmt"
	"github.com/up9inc/mizu/cli/mizu"
	"github.com/up9inc/mizu/cli/uiUtils"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type MizuTapOptions struct {
	GuiPort                uint16
	Namespace              string
	AllNamespaces          bool
	Analysis               bool
	AnalysisDestination    string
	KubeConfigPath         string
	MizuImage              string
	PlainTextFilterRegexes []string
	TapOutgoing            bool
	SleepIntervalSec       uint16
}

var mizuTapOptions = &MizuTapOptions{}
var direction string
var regex *regexp.Regexp

const analysisMessageToConfirm = `NOTE: running mizu with --analysis flag will upload recorded traffic
to UP9 cloud for further analysis and enriched presentation options.
`

var tapCmd = &cobra.Command{
	Use:   "tap [POD REGEX]",
	Short: "Record ingoing traffic of a kubernetes pod",
	Long: `Record the ingoing traffic of a kubernetes pod.
Supported protocols are HTTP and gRPC.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		RunMizuTap(regex, mizuTapOptions)
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("POD REGEX argument is required")
		} else if len(args) > 1 {
			return errors.New("unexpected number of arguments")
		}

		var compileErr error
		regex, compileErr = regexp.Compile(args[0])
		if compileErr != nil {
			return errors.New(fmt.Sprintf("%s is not a valid regex %s", args[0], compileErr))
		}

		directionLowerCase := strings.ToLower(direction)
		if directionLowerCase == "any" {
			mizuTapOptions.TapOutgoing = true
		} else if directionLowerCase == "in" {
			mizuTapOptions.TapOutgoing = false
		} else {
			return errors.New(fmt.Sprintf("%s is not a valid value for flag --direction. Acceptable values are in/any.", direction))
		}

		if mizuTapOptions.Analysis {
			fmt.Printf(analysisMessageToConfirm)
			if !uiUtils.AskForConfirmation("Would you like to proceed [y/n]: ") {
				fmt.Println("You can always run mizu without analysis, aborting")
				os.Exit(0)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tapCmd)

	tapCmd.Flags().Uint16VarP(&mizuTapOptions.GuiPort, "gui-port", "p", 8899, "Provide a custom port for the web interface webserver")
	tapCmd.Flags().StringVarP(&mizuTapOptions.Namespace, "namespace", "n", "", "Namespace selector")
	tapCmd.Flags().BoolVar(&mizuTapOptions.Analysis, "analysis", false, "Uploads traffic to UP9 for further analysis (Beta)")
	tapCmd.Flags().StringVar(&mizuTapOptions.AnalysisDestination, "dest", "up9.app", "Destination environment")
	tapCmd.Flags().Uint16VarP(&mizuTapOptions.SleepIntervalSec, "upload-interval", "", 10, "Interval in seconds for uploading data to UP9")
	tapCmd.Flags().BoolVarP(&mizuTapOptions.AllNamespaces, "all-namespaces", "A", false, "Tap all namespaces")
	tapCmd.Flags().StringVarP(&mizuTapOptions.KubeConfigPath, "kube-config", "k", "", "Path to kube-config file")
	tapCmd.Flags().StringVarP(&mizuTapOptions.MizuImage, "mizu-image", "", fmt.Sprintf("gcr.io/up9-docker-hub/mizu/%s:latest", mizu.Branch), "Custom image for mizu collector")
	tapCmd.Flags().StringArrayVarP(&mizuTapOptions.PlainTextFilterRegexes, "regex-masking", "r", nil, "List of regex expressions that are used to filter matching values from text/plain http bodies")
	tapCmd.Flags().StringVarP(&direction, "direction", "", "in", "Record traffic that goes in this direction (relative to the tapped pod): in/any")
}
