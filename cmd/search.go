package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/torbenconto/bible-cli/util"
	"log"
	"strings"

	"github.com/torbenconto/bible"
)

var searchCmd = &cobra.Command{
	Use:   "search [term]",
	Short: "Search for terms in the Bible",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the Bible from the context
		ctxBible := util.GetFromContext(cmd.Context())

		caseSensitive, err := cmd.Flags().GetBool("case-sensitive")
		if err != nil {
			log.Fatalf("Error getting case-sensitive: %s", err)
		}
		bookName, err := cmd.Flags().GetString("book")
		if err != nil {
			log.Fatalf("Error getting book name: %s", err)
		}
		maxResults, err := cmd.Flags().GetInt("max")
		if err != nil {
			log.Fatalf("Error getting max results: %s", err)
		}
		wordsOnly, err := cmd.Flags().GetBool("words-only")
		if err != nil {
			log.Fatalf("Error getting words-only: %s", err)
		}

		verse := args[0]

		var foundVerses []bible.Verse
		query := verse
		for _, book := range ctxBible.Books {
			// If bookName is provided, only search in that book
			if bookName != "" && book.Name != bookName {
				continue
			}
			for _, chapter := range book.Chapters {
				for _, v := range chapter.Verses {
					text := v.Text
					if !caseSensitive {
						query = strings.ToLower(query)
						text = strings.ToLower(text)
					}
					if wordsOnly {
						// Split text into words using whitespace and punctuation as delimiters
						words := strings.FieldsFunc(text, func(r rune) bool {
							return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'))
						})
						if containsPhrase(words, strings.Fields(query)) {
							foundVerses = append(foundVerses, v)
						}
					} else {
						if strings.Contains(text, query) {
							foundVerses = append(foundVerses, v)
						}
					}
				}
			}
		}
		if len(foundVerses) == 0 {
			log.Fatal("No results found")
		}

		if maxResults > 0 && len(foundVerses) > maxResults {
			foundVerses = foundVerses[:maxResults]
		}

		for _, v := range foundVerses {
			fmt.Println(v.Name, v.Text)
		}
	},
}

func init() {
	searchCmd.Flags().BoolP("case-sensitive", "c", false, "Make the search case sensitive")
	searchCmd.Flags().StringP("book", "b", "", "Specify a book to search in")
	searchCmd.Flags().BoolP("words-only", "w", true, "Only return results that match the entire query")
	searchCmd.Flags().IntP("max", "m", -1, "Maximum number of results to return")
	rootCmd.AddCommand(searchCmd)
}

// containsPhrase checks if the phrase is contained in the text, preserving word order
func containsPhrase(text, phrase []string) bool {
	for i := 0; i <= len(text)-len(phrase); i++ {
		found := true
		for j := range phrase {
			if text[i+j] != phrase[j] {
				found = false
				break
			}
		}
		if found {
			return true
		}
	}
	return false
}
