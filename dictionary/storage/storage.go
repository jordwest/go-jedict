package storage

import "errors"

// Error indicating that a matching word could not be found.
// Returned by Lookup(...) functions
var NotFoundError = errors.New("Word not found")

type Entry struct {
	// Unique identifier per record
	Sequence int

	Kanji []string

	// Readings in kana
	Readings []string

	// List of meanings in English
	Meanings []string

	PositionOfSpeech string

	Miscellaneous string

	Dialect string
}

// Implement StorageReader to provide lookup capability on a dictionary storage
type StorageReader interface {
	// Look up a word by kanji expression
	LookupKanji(string) (Entry, error)

	// Look up a word by hiragana/katakana expression
	LookupReading(string) (Entry, error)

	// Look up a word by English expression
	LookupMeaning(string) (Entry, error)

	// Look up a word by any of kanji, reading, or English expression
	// (returns the first result in that order)
	LookupWord(string) (Entry, error)
}

// Implement StorageWriter to provide dictionary persistence
type StorageWriter interface {
	// Store an entry in the database. The entry is not necessarily saved until
	// Commit() is called. The record is overwritten if one with the sequence
	// number already exists
	StoreEntry(Entry) error

	// Returns the number of entries that have not been committed to disk
	UncommittedEntries() int

	// Save all changes to the database
	Commit() error

	// Purge all entries from this storage (DESTRUCTIVE!)
	Purge() error
}

type StorageProvider interface {
	StorageReader
	StorageWriter
}
