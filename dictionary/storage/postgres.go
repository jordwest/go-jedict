package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorageProvider struct {
	connectionString   string
	uncommittedEntries []Entry
	db                 *sqlx.DB
}

func NewPostgresStorageProvider(connectionString string) *PostgresStorageProvider {
	s := PostgresStorageProvider{
		connectionString:   connectionString,
		uncommittedEntries: make([]Entry, 0),
	}

	s.db = sqlx.MustConnect("postgres", connectionString)

	return &s
}

func (s *PostgresStorageProvider) fetchEntry(entry_id int) (Entry, error) {
	entry := Entry{}
	err := s.db.Select(&entry.Kanji, "SELECT kanji FROM jedict_kanji WHERE entry_id=$1", entry_id)
	if err != nil {
		return entry, err
	}

	err = s.db.Select(&entry.Readings, "SELECT reading FROM jedict_reading WHERE entry_id=$1", entry_id)
	if err != nil {
		return entry, err
	}

	err = s.db.Select(&entry.Meanings, "SELECT meaning FROM jedict_meaning WHERE entry_id=$1", entry_id)

	return entry, err
}

func (s *PostgresStorageProvider) LookupKanji(w string) (Entry, error) {
	entry := Entry{}
	var entry_id int

	err := s.db.Get(&entry_id, "SELECT entry_id FROM jedict_kanji"+
		" WHERE kanji=$1;", w)

	if err != nil {
		return entry, err
	}

	return s.fetchEntry(entry_id)
}

func (s *PostgresStorageProvider) LookupReading(w string) (Entry, error) {
	entry := Entry{}
	var entry_id int

	err := s.db.Get(&entry_id, "SELECT entry_id FROM jedict_reading"+
		" WHERE reading=$1;", w)

	if err != nil {
		return entry, err
	}

	return s.fetchEntry(entry_id)
}

func (s *PostgresStorageProvider) LookupMeaning(w string) (Entry, error) {
	entry := Entry{}
	var entry_id int

	err := s.db.Get(&entry_id, "SELECT entry_id FROM jedict_meaning"+
		" WHERE meaning=$1;", w)

	if err != nil {
		return entry, err
	}

	return s.fetchEntry(entry_id)
}

func (s *PostgresStorageProvider) LookupWord(w string) (Entry, error) {
	entry, err := s.LookupKanji(w)
	if err == nil {
		return entry, nil
	}

	entry, err = s.LookupReading(w)
	if err == nil {
		return entry, nil
	}

	entry, err = s.LookupMeaning(w)
	if err == nil {
		return entry, nil
	}

	return Entry{}, err
}

func (s *PostgresStorageProvider) StoreEntry(entry Entry) error {
	s.uncommittedEntries = append(s.uncommittedEntries, entry)
	return nil
}

func (s *PostgresStorageProvider) Commit() error {
	// Always succeed if there are no entries to commit
	if len(s.uncommittedEntries) == 0 {
		return nil
	}

	for _, entry := range s.uncommittedEntries {
		err := s.saveSingleEntry(entry)
		if err != nil {
			return err
		}
	}
	s.uncommittedEntries = make([]Entry, 0)

	return nil
}

func (s *PostgresStorageProvider) saveSingleEntry(entry Entry) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	var entry_id int
	err = tx.Get(&entry_id, `INSERT INTO jedict_entries (sequence_id)`+
		`VALUES ($1)`+
		`RETURNING id;`, entry.Sequence)
	if err != nil {
		return err
	}

	for _, kanji := range entry.Kanji {
		_, err := tx.Exec("INSERT INTO jedict_kanji (entry_id, kanji)"+
			"VALUES ($1, $2);", entry_id, kanji)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, reading := range entry.Readings {
		_, err := tx.Exec("INSERT INTO jedict_reading (entry_id, reading)"+
			"VALUES ($1, $2);", entry_id, reading)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, meaning := range entry.Meanings {
		_, err := tx.Exec("INSERT INTO jedict_meaning (entry_id, meaning)"+
			"VALUES ($1, $2);", entry_id, meaning)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *PostgresStorageProvider) Purge() error {
	_, err := s.db.Exec("DELETE FROM jedict_entries;")
	if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM jedict_kanji;")
	if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM jedict_reading;")
	if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM jedict_meaning;")
	return err
}

func (s *PostgresStorageProvider) UncommittedEntries() int {
	return len(s.uncommittedEntries)
}
