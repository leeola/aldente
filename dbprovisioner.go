package aldente

import "sync"

// DbProvisioner wraps a provisioner and updates the record to the database.
type DbProvisioner struct {
	Provisioner
	updated bool
	db      Database
	mr      MachineRecord
	l       *sync.Mutex
}

func NewDbProvisioner(db Database, mr MachineRecord, p Provisioner) *DbProvisioner {
	return &DbProvisioner{
		Provisioner: p,
		db:          db,
		mr:          mr,
		l:           &sync.Mutex{},
	}
}

func (p *DbProvisioner) Record() (ProviderRecord, error) {
	pr, err := p.Provisioner.Record()
	if err != nil {
		return nil, err
	}

	p.l.Lock()
	defer p.l.Unlock()
	if !p.updated {
		p.mr.ProviderRecord = pr
		if err := p.db.UpdateMachine(p.mr); err != nil {
			return nil, err
		}
		p.updated = true
	}

	return pr, nil
}

func (p *DbProvisioner) Wait() error {
	// ensure Record gets called, so we always write to the db at least once.
	if _, err := p.Record(); err != nil {
		return err
	}

	return p.Wait()
}
