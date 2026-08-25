package hysteria2

import "sync"

// userIDAllocator maps user names to stable numeric IDs across hot updates.
// The underlying service API only carries an opaque int handle through the
// connection context, so reusing list indices misattributes traffic after
// users are added/removed and can panic when the new list is shorter.
type userIDAllocator struct {
	mu       sync.RWMutex
	next     int
	byName   map[string]int
	nameByID map[int]string
}

func newUserIDAllocator() *userIDAllocator {
	return &userIDAllocator{
		byName:   make(map[string]int),
		nameByID: make(map[int]string),
	}
}

// reset rebuilds the mapping for the given names, preserving IDs of names
// that survive the update. Returns the parallel id slice for the names.
func (a *userIDAllocator) reset(names []string) []int {
	a.mu.Lock()
	defer a.mu.Unlock()

	alive := make(map[int]struct{}, len(names))
	ids := make([]int, len(names))
	for i, name := range names {
		id, ok := a.byName[name]
		if !ok {
			id = a.next
			a.next++
			a.byName[name] = id
		}
		a.nameByID[id] = name
		alive[id] = struct{}{}
		ids[i] = id
	}
	for id := range a.nameByID {
		if _, ok := alive[id]; !ok {
			delete(a.nameByID, id)
		}
	}
	return ids
}

func (a *userIDAllocator) name(id int) (string, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	name, ok := a.nameByID[id]
	return name, ok
}
