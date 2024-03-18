package util

import (
	"bufio"
	"fmt"
	"github.com/torbenconto/bible"
	"github.com/torbenconto/bible-cli/config"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func LoadSourceFile(b *bible.Bible) *bible.Bible {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(filepath.Join(home, fmt.Sprintf(".bible/versions/%s/%s.txt", b.Version.Name, b.Version.Name)))
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Version %s not found locally", b.Version.Name)
			log.Println("Downloading the version")
			config.InitVersion(b.Version)

			// Bad but only way to make it look clean
			os.Exit(1)
		}
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		badline := false

		for _, badLine := range b.Version.BadLines {
			if strings.Contains(line, badLine) {
				badline = true
			}
		}

		if badline {
			continue
		}

		// Split the line into words
		words := strings.Fields(line)

		// Identify the book name, chapter number and verse
		bookName := ""
		chapterNumber := 0
		verseStartIndex := 0
		for i, word := range words {
			if unicode.IsDigit(rune(word[0])) {
				if i == 0 { // The first word starting with a digit is the chapter number
					chapterNumber, _ = strconv.Atoi(word)
				} else { // The second word starting with a digit is the start of the verse
					verseStartIndex = i
					break
				}
			} else {
				bookName += word + " "
			}
		}
		bookName = strings.TrimSpace(bookName)

		verseName := bookName + " " + words[verseStartIndex]
		verseText := strings.Join(words[verseStartIndex+1:], " ")

		// Check if the book already exists, if not, create a new book
		var currentBook *bible.Book
		for i := range b.Books {
			if b.Books[i].Name == bookName {
				currentBook = &b.Books[i]
				break
			}
		}
		if currentBook == nil {
			newBook := bible.NewBook(bookName, []bible.Chapter{})
			b.Books = append(b.Books, *newBook)
			currentBook = &b.Books[len(b.Books)-1]
		}

		// Check if the chapter already exists, if not, create a new chapter
		var currentChapter *bible.Chapter
		for i := range currentBook.Chapters {
			if currentBook.Chapters[i].Number == chapterNumber {
				currentChapter = &currentBook.Chapters[i]
				break
			}
		}
		if currentChapter == nil {
			newChapter := bible.Chapter{Number: chapterNumber, Verses: []bible.Verse{}}
			currentBook.Chapters = append(currentBook.Chapters, newChapter)
			currentChapter = &currentBook.Chapters[len(currentBook.Chapters)-1]
		}

		// Add the verse to the current chapter
		currentChapter.Verses = append(currentChapter.Verses, *bible.NewVerse(verseName, verseText))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return b
}
