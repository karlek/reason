// Package save implements functionality for saving the current game session.
package save

import (
	"bytes"
	"encoding/gob"
	"os"

	"github.com/karlek/reason/beastiary"
	"github.com/karlek/reason/fauna"
	"github.com/karlek/reason/item"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/object"
	"github.com/mewkiz/pkg/errorsutil"
)

// Save is a save file.
type Save struct {
	Path   string
	exists bool
}

func init() {
	// Register types to be saved / loaded.
	gob.Register(new(beastiary.Creature))
	gob.Register(new(fauna.Doodad))
	gob.Register(new(item.Item))
	gob.Register(new(object.Object))
	gob.Register(new(area.Area))
}

// New makes a new save with a reference to a path, this is used for both new
// saves and loading old ones.
func New(path string) (sav *Save, err error) {
	sav = &Save{Path: path}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		sav.exists = false
	} else {
		sav.exists = true
	}

	return sav, nil
}

// Skeleton is a wrapper for which objects to be saved.
type Skeleton struct {
	Area area.Area
	Hero beastiary.Creature
}

// Save writes the current game session to disk.
func (sav *Save) Save(a area.Area, hero beastiary.Creature) (err error) {
	sBlob := Skeleton{
		a,
		hero,
	}
	var buffer bytes.Buffer
	err = gob.NewEncoder(&buffer).Encode(&sBlob)
	if err != nil {
		return errorsutil.ErrorfColor("encode error: %#v", err)
	}

	f, err := os.Create(sav.Path)
	if err != nil {
		return errorsutil.ErrorfColor("%s", err)
	}
	_, err = f.Write(buffer.Bytes())
	if err != nil {
		return errorsutil.ErrorfColor("%s", err)
	}

	return nil
}

// Load returns a skeleton with all loaded information.
func (sav *Save) Load() (s *Skeleton, err error) {
	f, err := os.Open(sav.Path)
	if err != nil {
		return nil, errorsutil.ErrorfColor("%s", err)
	}

	sto := new(Skeleton)
	err = gob.NewDecoder(f).Decode(sto)
	if err != nil {
		return nil, errorsutil.ErrorfColor("decode error: %#v", err)
	}

	return sto, nil
}

// Exists returns true if a save file exists.
func (sav *Save) Exists() bool {
	return sav.exists
}
