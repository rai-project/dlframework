package dlframework

import (
	strings "strings"

	"golang.org/x/sync/syncmap"
)

func syncMapKeys(m *syncmap.Map) []string {
	keys := []string{}
	m.Range(func(key0 interface{}, value interface{}) bool {
		key, ok := key0.(string)
		if !ok {
			return true
		}
		keys = append(keys, key)
		return true
	})
	return keys
}

func CleanString(str string) string {
	r := strings.NewReplacer(" ", "_", "-", "_")
	res := r.Replace(str)
	return strings.ToLower(res)
}

func Partition(in []interface{}, partitionSize int) (out [][]interface{}) {
	cnt := (len(in)-1)/partitionSize + 1
	for i := 0; i < cnt; i++ {
		start := i * partitionSize
		end := (i + 1) * partitionSize
		if end > len(in) {
			end = len(in)
		}
		part := in[start:end]
		out = append(out, part)
	}

	return out
}

func PartitionStringList(in []string, partitionSize int) (out [][]string) {
	cnt := (len(in)-1)/partitionSize + 1
	for ii := 0; ii < cnt; ii++ {
		start := ii * partitionSize
		end := (ii + 1) * partitionSize
		if end > len(in) {
			end = len(in)
		}
		part := in[start:end]
		out = append(out, part)
	}

	return out
}
