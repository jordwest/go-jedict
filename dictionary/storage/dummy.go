package storage

import "fmt"

type DummyStorageProvider struct {
	dummyEntry         Entry
	uncommittedEntries int
}

func NewDummyStorageProvider() *DummyStorageProvider {
	s := DummyStorageProvider{

		dummyEntry: Entry{
			Sequence: 100,
		},
		uncommittedEntries: 0,
	}

	return &s
}

func (s *DummyStorageProvider) LookupKanji(w string) (Entry, error) {
	return s.dummyEntry, nil
}
func (s *DummyStorageProvider) LookupWord(w string) (Entry, error) {
	return s.dummyEntry, nil
}
func (s *DummyStorageProvider) LookupMeaning(w string) (Entry, error) {
	return s.dummyEntry, nil
}
func (s *DummyStorageProvider) LookupReading(w string) (Entry, error) {
	return s.dummyEntry, nil
}

func (s *DummyStorageProvider) StoreEntry(e Entry) error {
	s.uncommittedEntries += 1
	return nil
}

func (s *DummyStorageProvider) UncommittedEntries() int {
	return s.uncommittedEntries
}

func (s *DummyStorageProvider) Commit() error {
	fmt.Printf("Committed %d entries to dummy storage\n", s.uncommittedEntries)
	s.uncommittedEntries = 0

	return nil
}

func (s *DummyStorageProvider) Purge() error {
	fmt.Printf("Purge requested on dummy storage\n")
	return nil
}
