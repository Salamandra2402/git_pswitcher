package profile

import (
	"os"
	"testing"
)

func TestWriteOne(t *testing.T) {
	db := JsonFileDb{path: "/tmp/test_jsondb"}
	var err error
	profile := Profile{Name: "1", Email: "jsahd@as.com"}
	os.Create(db.path)

	if err = db.WriteProfiles([]Profile{profile}); err != nil {
		t.Errorf("Can't write one profile: %s", err)
	}

	// check file
	var profiles []Profile
	if profiles, err = db.GetProfiles(); err != nil {
		t.Errorf("Can't get profiles: %s", err)
	}

	if len(profiles) != 1 {
		t.Errorf("There must be the one element")
	}

}

func TestValidation(t *testing.T) {
	var err error
	profile := Profile{Name: "1"}
	if err = profile.SetEmail("a.user13@mail.ru"); err != nil {
		t.Errorf("email validation failed")
	}

}
