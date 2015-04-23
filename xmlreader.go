package jedict

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/jordwest/go-jedict/storage"
)

type jmDictXmlDoc struct {
	JmDict jmDict `xml:"JMdict"`
}

type jmDict struct {
	XMLName xml.Name `xml:"JMdict"`
	Entries []entry  `xml:"entry"`
}

type entry struct {
	XMLName xml.Name `xml:"entry"`
	Ent_seq int      `xml:"ent_seq"`
	K_ele   []k_ele  `xml:"k_ele"`
	R_ele   []r_ele  `xml:"r_ele"`
	//Sense   sense   `xml:"sense"`
	Gloss        []gloss `xml:"sense>gloss"`
	PartOfSpeech string  `xml:"sense>pos"`
}

type entity struct {
	XMLName xml.Name `xml:"!ENTITY"`
	content string   `xml:,innerxml`
}

type r_ele struct {
	Reb string `xml:"reb"`
}

type k_ele struct {
	Keb string `xml:"keb"`
}

type gloss struct {
	Content string `xml:",chardata"`
	Lang    string `xml:"lang,attr"`
}

type sense struct {
	gloss []string `xml:"gloss"`
	Pos   string   `xml:"pos"`
}

// Converts an XML read entity into an entry ready for storage
func (e *entry) convertToStorageEntry() storage.Entry {
	sEntry := storage.Entry{
		Sequence: e.Ent_seq,
		Kanji:    make([]string, 0),
		Readings: make([]string, 0),
		Meanings: make([]string, 0),
	}
	for _, reading := range e.R_ele {
		sEntry.Readings = append(sEntry.Readings, reading.Reb)
	}
	for _, kanji := range e.K_ele {
		sEntry.Kanji = append(sEntry.Kanji, kanji.Keb)
	}
	for _, meaning := range e.Gloss {
		// For now, skip languages other than the default (english)
		if meaning.Lang != "" {
			continue
		}
		sEntry.Meanings = append(sEntry.Meanings, meaning.Content)
	}

	return sEntry
}

func ReadXMLIntoStorage(filename string, provider storage.StorageWriter) error {
	reader, err := os.Open(filename)
	defer reader.Close()

	if err != nil {
		fmt.Printf("Error opening dictionary xml file: %s\n", err)
		return err
	}

	decoder := xml.NewDecoder(reader)
	decoder.Strict = false

	// Load directives
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.Directive: // Reading an XML directive
			directive := token.(xml.Directive)
			results, err := FindEntities(&directive)
			if err != nil {
				fmt.Printf("Entities could not be found: %s\n", err)
				continue
			}
			decoder.Entity = results
		case xml.StartElement: // Reading an XML element
			var entry entry
			if startElement.Name.Local == "entry" {
				decoder.DecodeElement(&entry, &startElement)
			}
			err = provider.StoreEntry(entry.convertToStorageEntry())
			if err != nil {
				fmt.Errorf("Error storing entry: %s\n %+v", err, entry)
				return err
			}
			if provider.UncommittedEntries() >= 5000 {
				err = provider.Commit()
				if err != nil {
					fmt.Errorf("Error committing entries: %s", err)
					return err
				}
			}
		}
	}

	provider.Commit()
	if err != nil {
		fmt.Errorf("Error committing entries: %s", err)
		return err
	}

	return nil

}

func FindEntities(d *xml.Directive) (map[string]string, error) {
	directiveStr := string(*d)

	// Go's xml package doesn't parse the DOCTYPE directive, we
	// need to parse this manually. Since we're only looking for
	// ENTITY definitions, a regex will do.
	re, err := regexp.Compile("<!ENTITY\\s([0-9A-z]+)\\s\"(.+)\">")
	if err != nil {
		fmt.Printf("Error compiling regex: %s", err)
		return nil, err
	}

	entities := make(map[string]string)

	results := re.FindAllStringSubmatch(directiveStr, -1)
	for _, entity := range results {
		entities[entity[1]] = entity[2]
	}

	return entities, nil
}

func (dic *jmDict) ReadInto(reader io.Reader) error {
	decoder := xml.NewDecoder(reader)
	decoder.Strict = false

	// Load directives
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.Directive: // Reading an XML directive
			directive := token.(xml.Directive)
			results, err := FindEntities(&directive)
			if err != nil {
				fmt.Printf("Entities could not be found: %s\n", err)
				continue
			}
			//fmt.Printf("Entities: %+v\n", results)
			decoder.Entity = results
		case xml.StartElement: // Reading an XML element
			var entry entry
			if startElement.Name.Local == "entry" {
				decoder.DecodeElement(&entry, &startElement)
			}
			dic.Entries = append(dic.Entries, entry)
		}
	}

	//return decoder.Decode(&dict)
	return nil

}
