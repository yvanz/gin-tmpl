package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Version struct {
	Major string
	Minor string
	Patch string
	Build string
}

var (
	AppVersion = Version{
		Major: "0", Minor: "0", Patch: "1",
		Build: "$Id$",
	}
)

func (v Version) String() string {
	ver := fmt.Sprintf("Version: %s.%s.%s", v.Major, v.Minor, v.Patch)

	return fmt.Sprintf("%s\nBuild: %s", ver, v.Build)
}

func (v *Version) Tiny() string {
	return fmt.Sprintf("%s.%s.%s-%s", v.Major, v.Minor, v.Patch, v.Build)
}

func NewVerCommand(serviceName string) *cobra.Command {
	if serviceName == "" {
		serviceName = "demo"
	}

	return &cobra.Command{
		Use:   "version",
		Short: "Prints version.",
		Run: func(*cobra.Command, []string) {
			fmt.Printf("%s\n%s\n", serviceName, AppVersion)
		},
	}
}
