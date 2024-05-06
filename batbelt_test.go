package batbelt_test

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func TestCreateJSONFile(t *testing.T) {
	t.Parallel()

	filepath := "./testdata/json_file.json"

	list := []struct {
		Name string `json:"name"`
	}{{Name: "Todd"}, {Name: "Sally"}, {Name: "Gizem"}}

	belt := batbelt.NewBatbelt()
	belt.CreateJSONFile(list, filepath)

	if belt.Error() != nil {
		t.Fatalf("could not remove file: %e", belt.Error())
	}

	b, err := os.ReadFile("testdata/json_file.json")
	if err != nil {
		t.Errorf("could not read json file: %e", err)
	}

	var readList []struct {
		Name string `json:"name"`
	}

	err = json.Unmarshal(b, &readList)
	if err != nil {
		t.Errorf("could not unmarshal json file: %e", err)
	}

	want := "Gizem"
	got := readList[2].Name

	if !cmp.Equal(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}

	belt.RemoveFile(filepath)
	if belt.Error() != nil {
		t.Fatalf("could not remove file: %e", belt.Error())
	}
}

func ExampleBatbelt_CreateJSONFile() {
	filepath := "./example_file.json"

	list := []struct {
		Name string `json:"name"`
	}{{Name: "Philipp"}, {Name: "Sally"}, {Name: "Gizem"}}

	belt := batbelt.NewBatbelt()
	belt.CreateJSONFile(list, filepath)

	if belt.Error() != nil {
		// handle error
	}
}

type Person struct {
	Name string `json:"name"`
}

func TestReadJSONFile(t *testing.T) {
	t.Parallel()

	filepath := "./testdata/json_read_file.json"

	list := []Person{{Name: "Todd"}, {Name: "Sally"}, {Name: "Gizem"}}

	var readList []Person

	belt := batbelt.NewBatbelt()
	belt.CreateJSONFile(list, filepath)

	people, err := batbelt.ReadJSONFile[[]Person](readList, filepath)
	if err != nil {
		t.Fatalf("could not create and read: %e", belt.Error())
	}

	want := "Gizem"
	got := people[2].Name

	if !cmp.Equal(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}

	belt.RemoveFile(filepath)
	if belt.Error() != nil {
		t.Fatalf("could not remove file: %e", belt.Error())
	}

}

func ExampleBatbelt_ReadJSONFile() {
	filepath := "./testdata/json_read_file.json"

	list := []Person{{Name: "Todd"}, {Name: "Sally"}, {Name: "Gizem"}}

	var readList []Person

	belt := batbelt.NewBatbelt()
	belt.CreateJSONFile(list, filepath)

	_, err := batbelt.ReadJSONFile[[]Person](readList, filepath)
	if err != nil {
		// handle error
	}

	// Output: [{Name: "Todd"}, {Name: "Sally"}, {Name: "Gizem"}]
}
