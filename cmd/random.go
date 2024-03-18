package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/torbenconto/bible"
	"github.com/torbenconto/bible-cli/util"
	"log"
	"math/rand"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get a random verse from the Bible",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the Bible from the context
		ctxBible := util.GetFromContext(cmd.Context())

		bookName := cmd.Flag("book").Value.String()
		count, err := cmd.Flags().GetInt("count")
		if err != nil {
			log.Fatalf("Error getting count: %s", err)
		}

		var book bible.Book
		if bookName != "" {
			for _, b := range ctxBible.Books {
				if b.Name == bookName {
					book = b
					break
				}
			}

			if book.Name == "" {
				log.Fatalf("Custom book not found, run bible books to retrieve a list of available books")
			}
		} else {
			book = ctxBible.Books[rand.Intn(len(ctxBible.Books))]
		}

		for i := 0; i < count; i++ {
			// Select a random chapter from the book
			chapter := book.Chapters[rand.Intn(len(book.Chapters))]

			// Select a random verse from the chapter
			verse := chapter.Verses[rand.Intn(len(chapter.Verses))]
			fmt.Println(verse.Name, verse.Text)
		}
	},
}

func init() {
	randomCmd.Flags().IntP("count", "c", 1, "Number of random verses to get")
	randomCmd.Flags().StringP("book", "b", "", "Specify a book to get random verses from")
	rootCmd.AddCommand(randomCmd)
}
