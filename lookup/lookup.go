/*
`lookup` is a simple command line tool for testing lookups on the
dictionary.

Usage:
	lookup -db CONNECTION_STRING [-kanji EXPRESSION]
	                             [-meaning EXPRESSION]
	                             [-reading EXPRESSION]

CONNECTION_STRING is the postgresql connection string for the database containing
the dictionary. The database must be prepopulated.

EXPRESSION is a kanji expression, english expression, or hiragana/katakana
expression respectively.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/jordwest/go-jedict/storage"
	"strings"
)

func printEntry(entry storage.Entry) {
	fmt.Printf("%s\n", strings.Join(entry.Kanji, " "))
	fmt.Printf("%s\n----\n", strings.Join(entry.Readings, " "))
	fmt.Printf("%s\n", strings.Join(entry.Meanings, " "))
}

func main() {
	// Connect to the dictionary database using the provided connetion string
	var kanji string
	var reading string
	var meaning string
	var connectionString string
	flag.StringVar(&connectionString, "db", "", "PostgreSQL connection string")
	flag.StringVar(&kanji, "kanji", "", "Lookup word by kanji expression")
	flag.StringVar(&reading, "reading", "", "Lookup word by reading")
	flag.StringVar(&meaning, "meaning", "", "Lookup word by english meaning")
	flag.Parse()

	provider := storage.NewPostgresStorageProvider(connectionString)

	// -kanji option for looking up an expression by kanji
	if kanji != "" {
		entry, err := provider.LookupKanji(kanji)
		if err != nil {
			fmt.Errorf("Error looking up kanji: %s\n", err)
			return
		}
		printEntry(entry)
	}

	if reading != "" {
		entry, err := provider.LookupReading(reading)
		if err != nil {
			fmt.Errorf("Error looking up reading: %s\n", err)
			return
		}
		printEntry(entry)
	}

	if meaning != "" {
		entry, err := provider.LookupMeaning(meaning)
		if err != nil {
			fmt.Errorf("Error looking up meaning: %s\n", err)
			return
		}
		printEntry(entry)
	}
}
