package display

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// ==============================================================================================================================
//                                      SORTED DATA
// ==============================================================================================================================

type dataInterface interface {
	display() string
	update(m message) error
	key() string
	getLastUpdate() time.Time
	setStatus(status string) // Testing only!
}

type sortedDataInterface interface {
	display() []string
}

// -----------------------------------------------------------------------------------------------------------------------------

type sortedData struct {
	mType   string
	maxSize int
	sortAsc bool

	sync.RWMutex
	data   map[string]dataInterface
	ind    []string
	sorted bool
}

func newSortedData(mType string, sortAsc bool) *sortedData {
	return &sortedData{
		mType:  mType,
		data:   make(map[string]dataInterface),
		ind:    make([]string, 0),
		sorted: false,
	}
}

func (r *sortedData) update(m message) error {
	if _, ok := r.data[m.key]; !ok {
		return r.add(m)
	}

	r.Lock()
	defer r.Unlock()
	r.data[m.key].update(m)
	return nil
}

func (r *sortedData) add(m message) (err error) {
	var d dataInterface
	switch m.mType {
	case msgTypeES:
		d, err = newEngStatusType(m)
		if err != nil {
			return err
		}
	case msgTypeER:
		d, err = newEngRequestType(m)
		if err != nil {
			return err
		}
	case msgTypeEA:
		d, err = newEngAdpRequestType(m)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid message type: %q", m.mType)
	}

	r.Lock()
	defer r.Unlock()
	r.data[d.key()] = d
	r.sorted = false

	return nil
}

func (r *sortedData) display() []string {
	r.sort()

	var list []string
	r.RLock()
	defer r.RUnlock()
	for _, name := range r.ind {
		list = append(list, r.data[name].display())
	}
	return list
}

func (r *sortedData) sort() {
	if r.sorted {
		return
	}
	r.Lock()
	defer r.Unlock()
	r.ind = make([]string, 0)
	if len(r.data) > 0 {
		for k := range r.data {
			r.ind = append(r.ind, k)
		}
		if r.sortAsc {
			sort.Strings(r.ind)
		} else {
			sort.Sort(sort.Reverse(sort.StringSlice(r.ind)))
		}

	}
	r.sorted = true
}

func (r *sortedData) get() []dataInterface {
	r.sort()

	var result []dataInterface
	r.RLock()
	defer r.RUnlock()
	for _, name := range r.ind {
		result = append(result, r.data[name])
	}
	return result
}