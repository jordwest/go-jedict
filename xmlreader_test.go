package jedict

import (
	"os"
	"testing"

	"github.com/jordwest/go-jedict/storage"
)

/*
func Test_ReadXMLIntoStorage(t *testing.T) {
	p := storage.NewDummyStorageProvider()
	err := ReadXMLIntoStorage("./JMdict", p)
	if err != nil {
		t.Error(err)
	}
}
*/

func Test_ReadXMLIntoPostgresStorage(t *testing.T) {
	p := storage.NewPostgresStorageProvider(os.Getenv("YOMUYOMU_DBURL"))
	err := ReadXMLIntoStorage("./JMdict", p)
	if err != nil {
		t.Error(err)
	}
}
