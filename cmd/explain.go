package cmd

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/torbenconto/bible-cli/util"
	"strings"
)

var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Explain a Bible verse using AI features",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the Bible from the context
		ctxBible := util.GetFromContext(cmd.Context())

		// Get the api key from the context
		apiKey := cmd.Context().Value("api_key").(string)

		verse := strings.ToLower(args[0])

		// Initialize OpenAI client
		client := openai.NewClient(apiKey)

		explainVerse := func(verseText string) (string, error) {
			response, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: openai.GPT3Dot5Turbo,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: fmt.Sprintf("Explain the meaning of the following Bible verse: %s", verseText),
						},
					},
				},
			)
			if err != nil {
				return "", err
			}

			return response.Choices[0].Message.Content, nil
		}

		// Call the explanation function and print the result
		for _, book := range ctxBible.Books {
			for _, chapter := range book.Chapters {
				for _, v := range chapter.Verses {
					if strings.ToLower(v.Name) == verse {
						explanation, err := explainVerse(v.Text)
						if err != nil {
							fmt.Println("Error explaining verse:", err)
							return
						}

						fmt.Println("Explanation:", explanation)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
