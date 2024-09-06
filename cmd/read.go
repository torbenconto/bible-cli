package cmd

import (
	"bufio"
	"fmt"
	"github.com/torbenconto/bible-cli/util"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read [starting] [flags]",
	Short: "Read the Bible",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the Bible from the context
		ctxBible := util.GetFromContext(cmd.Context())

		reader := bufio.NewReader(os.Stdin)

		if len(args) == 0 {
			args = append(args, "Genesis 1:1")
		}

		startVerse := args[0]

		startReading := false

		for _, book := range ctxBible.Books {
			for _, chapter := range book.Chapters {
				for _, verse := range chapter.Verses {
					// Remove book name from verse name
					if strings.ToLower(verse.Name) == strings.ToLower(startVerse) {
						fmt.Println("Press enter to show the next verse or q and then enter to quit")
						startReading = true
					}

					if startReading {
						name := strings.Replace(verse.Name, book.Name+" ", "", 1)
						fmt.Printf("%s: %s\n", name, verse.Text)
						input, _ := reader.ReadString('\n')
						if input == "q\n" {
							return
						}
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
