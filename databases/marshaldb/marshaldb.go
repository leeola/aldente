package marshaldb

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	ald "github.com/leeola/aldente"
	"github.com/leeola/errors"
)

// data is the core data structure written to json on the fs.
type data struct {
	Groups map[string]map[string]ald.MachineRecord
}

type Config struct {
	Path string `toml:"path"`
}

type MarshalDb struct {
	config Config
}

func New(c Config) (*MarshalDb, error) {
	if c.Path == "" {
		return nil, errors.New("missing required config: Path")
	}

	if err := os.MkdirAll(filepath.Dir(db.config.Path), 0755); err != nil {
		return nil, err
	}

	return &MarshalDb{
		config: c,
	}, nil
}

func (db *MarshalDb) load() (data, error) {
	f, err := os.OpenFile(db.config.Path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return data{}, errors.Wrap(err, "failed to open db")
	}
	defer f.Close()

	var d data
	err = json.NewDecoder(f).Decode(&d)

	// If it's eof, create a new map and return that.
	if err == io.EOF {
		return data{
			Groups: map[string]map[string]ald.MachineRecord{},
		}, nil
	}

	if err != nil {
		return data{}, errors.Wrap(err, "failed to decode db")
	}

	return d, nil
}

func (db *MarshalDb) save(d data) error {
	f, err := os.OpenFile(db.config.Path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open db")
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(&d); err != nil {
		return errors.Wrap(err, "failed to encode db")
	}

	return nil
}

func (db *MarshalDb) Groups() ([]string, error) {
	d, err := db.load()
	if err != nil {
		return nil, err
	}

	// build a slice from the map
	s := make([]string, len(d.Groups))
	var i int
	for n, _ := range d.Groups {
		s[i] = n
		i++
	}

	return s, nil
}

func (db *MarshalDb) GroupMachines(group string) ([]ald.MachineRecord, error) {
	d, err := db.load()
	if err != nil {
		return nil, err
	}

	ms, ok := d.Groups[group]
	if !ok {
		return nil, errors.Errorf("group not found: %s", group)
	}

	// build a slice from the map
	s := make([]ald.MachineRecord, len(ms))
	var i int
	for _, m := range ms {
		s[i] = m
		i++
	}

	return s, nil
}

func (db *MarshalDb) CreateGroup(group string, machines []ald.MachineConfig) error {
	d, err := db.load()
	if err != nil {
		return err
	}

	if _, ok := d.Groups[group]; ok {
		return errors.Errorf("group already exists: %s", group)
	}

	ms := map[string]ald.MachineRecord{}
	for _, m := range machines {
		if m.Name == "" {
			return errors.New("missing required field: Name")
		}
		if m.Provider == "" {
			return errors.New("missing required field: Provider")
		}

		// ensure no duplicate names
		if _, ok := ms[m.Name]; ok {
			return errors.Errorf("duplicate machine name: %s", m.Name)
		}

		ms[m.Name] = ald.MachineRecord{
			Name:     m.Name,
			Group:    group,
			Provider: m.Provider,
		}
	}

	d.Groups[group] = ms
	return db.save(d)
}

func (db *MarshalDb) UpdateMachine(m ald.MachineRecord) error {
	if m.Name == "" {
		return errors.New("missing required field: Name")
	}
	if m.Group == "" {
		return errors.New("missing required field: Group")
	}
	if m.Provider == "" {
		return errors.New("missing required field: Provider")
	}

	d, err := db.load()
	if err != nil {
		return err
	}

	ms, ok := d.Groups[m.Group]
	if !ok {
		return errors.Errorf("group not found: %s", m.Group)
	}

	// only update a machine if it exists in the group.
	// Ie, it may not have existed in the config when the group was created.
	if _, ok := ms[m.Name]; !ok {
		return errors.Errorf("group does not contain record for machine: %s", m.Name)
	}

	// update the record, and then save the struct
	ms[m.Name] = m
	return db.save(d)
}
