/*
`jedict` is the command line tool for importing the JMDICT XML dictionary
and performing lookups on the database.

Usage:
	jedict -db CONNECTION_STRING [-import XML_PATH]
		                         [-kanji EXPRESSION]
	                             [-meaning EXPRESSION]
	                             [-reading EXPRESSION]

XML_PATH is the location of the JMDICT XML file. See README.md for more details.

CONNECTION_STRING is the postgresql connection string for the database containing
the dictionary. The database must be prepopulated.

EXPRESSION is a kanji expression, english expression, or hiragana/katakana
expression respectively.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/jordwest/go-jedict/dictionary"
	"github.com/jordwest/go-jedict/dictionary/storage"
	"strings"
)

func printEntry(entry storage.Entry) {
	fmt.Printf("%s\n", strings.Join(entry.Kanji, " "))
	fmt.Printf("%s\n----\n", strings.Join(entry.Readings, " "))
	fmt.Printf("%s\n", strings.Join(entry.Meanings, ", "))
}

func doImport(p storage.StorageWriter, xmlfile string) error {
	progress := make(chan float32)

	// Print out the progress while the import is running
	go func(progress chan float32) {
		incomplete := true
		var completion float32 = 0.0
		for incomplete == true {
			completion, incomplete = <-progress
			fmt.Printf("\rImport progress: %.2f%%", completion*100)
		}
		fmt.Printf("\n")
	}(progress)

	err := dictionary.ReadXMLIntoStorage("./JMdict", p, progress)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Connect to the dictionary database using the provided connetion string
	var xmlfile string
	var kanji string
	var reading string
	var meaning string
	var connectionString string
	flag.StringVar(&connectionString, "db", "", "PostgreSQL connection string")
	flag.StringVar(&xmlfile, "import", "", "JMdict file to import")
	flag.StringVar(&kanji, "kanji", "", "Lookup word by kanji expression")
	flag.StringVar(&reading, "reading", "", "Lookup word by reading")
	flag.StringVar(&meaning, "meaning", "", "Lookup word by english meaning")
	flag.Parse()

	provider := storage.NewPostgresStorageProvider(connectionString)

	// -import option for performing database import
	if xmlfile != "" {
		err := doImport(provider, xmlfile)
		if err != nil {
			fmt.Errorf("Error importing dictionary: %s\n", err)
			return
		}
	}

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
