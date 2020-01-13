package memorykg

import (
	"errors"
	"sync"

	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
)

// KeygroupStorage saves a set of all available keygroup.
type KeygroupStorage struct {
	keygroups map[commons.KeygroupName]struct{}
	sync.RWMutex
}

// New creates a new KeygroupStorage.
func New() (kS *KeygroupStorage) {
	kS = &KeygroupStorage{
		keygroups: make(map[commons.KeygroupName]struct{}),
	}

	return
}

// Create adds a keygroup to the KeygroupStorage.
func (kS *KeygroupStorage) Create(k keygroup.Keygroup) error {
	log.Debug().Msgf("CreateKeygroup from memorykg: in %#v", k)
	kS.RLock()
	_, ok := kS.keygroups[k.Name]
	kS.RUnlock()

	if ok {
		return nil
	}

	kS.Lock()
	kS.keygroups[k.Name] = struct{}{}
	kS.Unlock()

	return nil
}

// Delete removes a keygroup from the KeygroupStorage.
func (kS *KeygroupStorage) Delete(k keygroup.Keygroup) error {
	log.Debug().Msgf("DeleteKeygroup from memorykg: in %#v", k)
	kS.RLock()
	_, ok := kS.keygroups[k.Name]
	kS.RUnlock()

	if !ok {
		return errors.New("memorykg: no such keygroup")
	}

	kS.Lock()
	delete(kS.keygroups, k.Name)
	kS.Unlock()

	return nil
}

// Exists checks if a keygroup exists in the KeygroupStorage.
func (kS *KeygroupStorage) Exists(k keygroup.Keygroup) bool {
	kS.RLock()
	_, ok := kS.keygroups[k.Name]
	kS.RUnlock()

	log.Debug().Msgf("Exists from memorykg: in %#v, out %#v", k, ok)

	return ok
}
