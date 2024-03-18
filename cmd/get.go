package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/torbenconto/bible-cli/util"
	"log"
)

var getCmd = &cobra.Command{
	Use:   "get [verse]",
	Short: "Get a verse from the Bible",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Get the Bible from the context
		ctxBible := util.GetFromContext(cmd.Context())

		verse := args[0]

		verses := ctxBible.GetVerse(verse)

		if len(verses) == 0 {
			log.Fatal("No results found")
		}

		for _, v := range verses {
			fmt.Println(v.Name, v.Text)
		}

	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
