package cmd

import (
	"TT-Watch/schedule/watch"

	"github.com/spf13/cobra"
)

var Watch = &cobra.Command{
	Use: "watch",
	Run: watch.Run,
}
