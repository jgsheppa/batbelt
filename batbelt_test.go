package batbelt_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
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

	people, err := batbelt.ReadJSONFile(readList, filepath)
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

func TestReadJSONFile_String(t *testing.T) {
	t.Parallel()

	filepath := "./testdata/json_read_string_file.json"

	exampleText := "Example"

	var readList string

	belt := batbelt.NewBatbelt()
	belt.CreateJSONFile(exampleText, filepath)

	got, err := batbelt.ReadJSONFile(readList, filepath)
	if err != nil {
		t.Fatalf("could not create and read: %e", belt.Error())
	}

	want := "Example"

	if !cmp.Equal(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}

	belt.RemoveFile(filepath)
	if belt.Error() != nil {
		t.Fatalf("could not remove file: %e", belt.Error())
	}
}

func ExampleReadJSONFile() {
	filepath := "./testdata/json_read_file.json"

	list := []Person{{Name: "Todd"}, {Name: "Sally"}, {Name: "Gizem"}}

	var readList []Person

	belt := batbelt.NewBatbelt()
	belt.CreateJSONFile(list, filepath)

	people, err := batbelt.ReadJSONFile(readList, filepath)
	if err != nil {
		// handle error
	}

	fmt.Println(people)
	// Output: [{Todd} {Sally} {Gizem}]
}

func TestGeneratePassword_Length(t *testing.T) {
	t.Parallel()

	password := batbelt.GeneratePassword("abcdefg", 8)

	want := 8
	got := len(password)

	if !cmp.Equal(got, want) {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestGeneratePassword_Characters(t *testing.T) {
	t.Parallel()

	chars := "abcdefg"
	password := batbelt.GeneratePassword(chars, 8)

	if !strings.ContainsAny(password, chars) {
		t.Errorf("password contains unwanted characters: %s", password)
	}
}

func TestGeneratePassword_UnwantedCharacters(t *testing.T) {
	t.Parallel()

	chars := "abcdefg"
	password := batbelt.GeneratePassword(chars, 8)

	unwantedChars := "jklmnop"

	if strings.ContainsAny(password, unwantedChars) {
		t.Errorf("password contains unwanted characters: %s", password)
	}
}

func ExampleGeneratePassword() {
	chars := "abcdefg"
	password := batbelt.GeneratePassword(chars, 8)

	fmt.Println(len(password))
	// Output: 8
}
