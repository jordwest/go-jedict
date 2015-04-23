package jedict

import "./storage"

type Dictionary struct {
	reader *storage.StorageReader
}

func NewDictionary(storageReader *StorageReader) (d Dictionary, err error) {
	d = Dictionary{}
	d.reader = storageReader
}

func (d Dictionary) FindOneByKanji(kanji string) (Entry, error) {
	return d.reader.LookupKanji(kanji)
}

func (d Dictionary) FindOneByReading(reading string) (Entry, error) {
	return d.reader.LookupReading(reading)
}

func (d Dictionary) FindOneByMeaning(meaning string) (Entry, error) {
	return d.reader.LookupMeaning(meaning)
}
