package motley

type Motley interface {
	Status(groupName string) error
}
