package marshaldb

import (
	"encoding/json"
	"io"
	"os"

	ald "github.com/leeola/aldente"
	"github.com/leeola/errors"
)

// data is the core data structure written to json on the fs.
type data struct {
	Machines []ald.MachineRecord
}

type MarshalDb struct {
	path string
}

func New(path string) (*MarshalDb, error) {
	return &MarshalDb{
		path: path,
	}, nil
}

func (db *MarshalDb) load() (data, error) {
	f, err := os.OpenFile(db.path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return data{}, errors.Wrap(err, "failed to open db")
	}
	defer f.Close()

	var d data
	err = json.NewDecoder(f).Decode(&d)
	if err != nil && err != io.EOF {
		return data{}, errors.Wrap(err, "failed to decode db")
	}

	return d, nil
}

func (db *MarshalDb) save(d data) error {
	f, err := os.OpenFile(db.path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open db")
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(&d); err != nil {
		return errors.Wrap(err, "failed to encode db")
	}

	return nil
}

func (db *MarshalDb) Add(m ald.MachineRecord) error {
	d, err := db.load()
	if err != nil {
		return err
	}

	d.Machines = append(d.Machines, m)

	return db.save(d)
}

func (db *MarshalDb) List() ([]ald.MachineRecord, error) {
	d, err := db.load()
	if err != nil {
		return nil, err
	}

	return d.Machines, nil
}
