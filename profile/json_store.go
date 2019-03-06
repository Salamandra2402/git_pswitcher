package profile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	dbFile = "/home/%s/.git_pswitcher"
)

//JsonFileDb work with simple json file
type JsonFileDb struct {
	path string
}

//CreateDefaultJsonFileDb creates default db (inside the Linux home dir)
func CreateDefaultJsonFileDb() JsonFileDb {
	user := os.Getenv("USER")
	path := fmt.Sprintf(dbFile, user)
	return JsonFileDb{path: path}
}

//GetProfiles get all profiles as a list of Profile structures
func (db JsonFileDb) GetProfiles() ([]Profile, error) {
	var err error
	if _, err := os.Stat(db.path); err != nil {
		var tempFile *os.File
		if tempFile, err = os.Create(db.path); err != nil {
			return nil, fmt.Errorf("Can't create file <%s>: %s", db.path, err)
		}
		tempFile.Close()
		return make([]Profile, 0), nil
	}
	byteValue, _ := ioutil.ReadFile(db.path)
	var profiles []Profile
	if err = json.Unmarshal([]byte(byteValue), &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}

//WriteProfiles overwrite db
func (db JsonFileDb) WriteProfiles(profiles []Profile) error {
	var err error
	var writeFile *os.File
	if writeFile, err = os.OpenFile(db.path, os.O_WRONLY, os.ModeExclusive); err != nil {
		return err
	}
	defer writeFile.Close()
	var b []byte
	if b, err = json.Marshal(profiles); err != nil {
		return err
	}
	writeFile.Truncate(0)
	writeFile.Seek(0, 0)
	if _, err = writeFile.Write(b); err != nil {
		return err
	}

	return nil
}

//AddProfile add new profile to db
func (db JsonFileDb) AddProfile(p Profile, force bool) error {
	var err error
	var profiles []Profile
	if profiles, err = db.GetProfiles(); err != nil {
		return err
	}

	for num, elem := range profiles {
		if elem.Name == p.Name {
			if force {
				// update email, because forced
				profiles[num].Email = p.Email
				// write to file
				if err = db.WriteProfiles(profiles); err != nil {
					return err
				}
				return nil
			} else {
				return fmt.Errorf("user with this name is already exist")
			}
		}
	}

	profiles = append(profiles, p)

	// write to file
	if err = db.WriteProfiles(profiles); err != nil {
		return err
	}
	return nil
}

//RemoveProfile remove a profile from db
func (db JsonFileDb) RemoveProfile(name string) error {
	var err error
	var profiles []Profile

	if profiles, err = db.GetProfiles(); err != nil {
		return err
	}

	for num, elem := range profiles {
		if elem.Name == name {
			// found elem to delete
			slice := append(profiles[:num], profiles[num+1:]...)
			if err = db.WriteProfiles(slice); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("can't find a profile to remove")
}

//GetProfile get profile by name
func (db JsonFileDb) GetProfile(name string) (Profile, error) {
	var err error
	var profiles []Profile
	if profiles, err = db.GetProfiles(); err != nil {
		return Profile{}, err
	}
	for _, item := range profiles {
		if item.Name == name {
			return item, nil
		}
	}
	return Profile{}, fmt.Errorf("Can't find the profile with this name")
}
