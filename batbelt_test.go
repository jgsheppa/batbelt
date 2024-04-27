package batbelt_test

import (
	"errors"
	"os"
	"testing"

	"github.com/jgsheppa/batbelt"
)

func TestRemoveFile(t *testing.T) {
	t.Parallel()

	filepath := "./testdata/file_to_remove.txt"

	belt := batbelt.NewBatbelt()
	belt.RemoveFile(filepath)

	if belt.Error() != nil {
		t.Fatalf("could not remove file: %e", belt.Error())
	}

	if _, err := os.Stat(filepath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("file should not exist: %e", err)
	}

	if err := os.WriteFile(filepath, []byte("remove me for testing purposes!"), 0644); err != nil {
		t.Fatalf("could not create testfile: %e", err)
	}
}

func ExampleBatbelt_RemoveFile() {
	filepath := "./testdata/file_to_remove.txt"

	belt := batbelt.NewBatbelt()
	belt.RemoveFile(filepath)
	if belt.Error() != nil {
		// return or log error...
	}

}
