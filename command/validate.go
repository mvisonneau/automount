package command

import (
	"os"

	"github.com/fatih/color"
	"github.com/mvisonneau/automount/pkg/fs"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

// Validate checks if dependencies are available on the host
func Validate(ctx *cli.Context) error {
	if err := configure(ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)

	table.SetHeader([]string{"Command", "Mandatory", "Available"})
	table.Append([]string{"blkid", "YES", boolToString(fs.IsBlkidAvailable())})
	table.Append([]string{"lsblk", "YES", boolToString(fs.IsLsblkAvailable())})
	table.Append([]string{"lvm", "NO", boolToString(fs.IsLvmAvailable())})
	table.Append([]string{"mdadm", "NO", boolToString(fs.IsMdadmAvailable())})

	table.Render()
	return nil
}

func boolToString(b bool) string {
	if b {
		return green("YES")
	}
	return red("NO")
}
