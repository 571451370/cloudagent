package qga

var (
	store = make(map[string]map[interface{}]interface{})
)

// StoreSet set variable k with value v to namespace ns
func StoreSet(ns string, k interface{}, v interface{}) {
	store[ns][k] = v
}

// StoreGet get variable k from namespace ns
func StoreGet(ns string, k interface{}) (interface{}, bool) {
	m, ok := store[ns]
	if !ok {
		return nil, false
	}
	v, ok := m[k]
	return v, ok
}

// StoreDel del variable k from namespace ns
func StoreDel(ns string, k interface{}) {
	m, ok := store[ns]
	if !ok {
		return
	}
	delete(m, k)
}
