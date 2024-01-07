package players

import (
	"errors"
)

var (
	NotFoundErr = errors.New("not found")
)

type MemStore struct {
	list map[string]Player
}

func NewMemStore() *MemStore {
	list := make(map[string]Player)
	return &MemStore{
		list,
	}
}

func (m MemStore) Add(kitNumber string, player Player) error {
	m.list[kitNumber] = player
	return nil
}

func (m MemStore) Get(kitNumber string) (Player, error) {
	if val, ok := m.list[kitNumber]; ok {
		return val, nil
	}
	return Player{}, NotFoundErr
}

func (m MemStore) List() (map[string]Player, error) {
	return m.list, nil
}

func (m MemStore) Update(kitNumber string, player Player) error {
	if _, ok := m.list[kitNumber]; ok {
		m.list[kitNumber] = player
		return nil
	}
	return NotFoundErr
}

func (m MemStore) Remove(kitNumber string) error {
	delete(m.list, kitNumber)
	return nil
}
