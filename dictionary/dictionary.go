package dictionary

import "github.com/jordwest/go-jedict/dictionary/storage"

type Dictionary struct {
	reader storage.StorageReader
}

func NewDictionary(storageReader storage.StorageReader) (d Dictionary, err error) {
	d = Dictionary{}
	d.reader = storageReader
	return d, err
}

func (d Dictionary) FindOneByKanji(kanji string) (storage.Entry, error) {
	return d.reader.LookupKanji(kanji)
}

func (d Dictionary) FindOneByReading(reading string) (storage.Entry, error) {
	return d.reader.LookupReading(reading)
}

func (d Dictionary) FindOneByMeaning(meaning string) (storage.Entry, error) {
	return d.reader.LookupMeaning(meaning)
}

/* TODO:

func (d Dictionary) FindAllByKanji(meaning string, limit int, offset int) ([]storage.Entry, error)
func (d Dictionary) FindAllByReading(meaning string, limit int, offset int) ([]storage.Entry, error)
func (d Dictionary) FindAllByMeaning(meaning string, limit int, offset int) ([]storage.Entry, error)

*/
