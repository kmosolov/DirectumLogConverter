package DirectumLogConverter

import (
	jsoniter "github.com/json-iterator/go"
	"sort"
)

type MapItem struct {
	Key   string
	Value interface{}
	index uint64
}

type MapSlice []MapItem

func (ms MapSlice) Len() int           { return len(ms) }
func (ms MapSlice) Less(i, j int) bool { return ms[i].index < ms[j].index }
func (ms MapSlice) Swap(i, j int)      { ms[i], ms[j] = ms[j], ms[i] }

var indexCounter uint64

func nextIndex() uint64 {
	indexCounter++
	return indexCounter
}

func (ms *MapSlice) UnmarshalJSON(b []byte) error {
	m := map[string]MapItem{}
	var json = jsoniter.ConfigDefault
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	for k, v := range m {
		*ms = append(*ms, MapItem{Key: k, Value: v.Value, index: v.index})
	}
	sort.Sort(*ms)
	return nil
}

func (mi *MapItem) UnmarshalJSON(b []byte) error {
	var v interface{}
	var json = jsoniter.ConfigDefault
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	mi.Value = v
	mi.index = nextIndex()
	return nil
}
