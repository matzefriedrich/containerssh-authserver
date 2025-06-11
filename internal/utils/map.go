package utils

type MapMergeStrategy int

const (
	AddMissing MapMergeStrategy = 1 << iota
	Override
	Remove
)

func MergeMaps[K comparable, V any](original, candidate map[K]V, strategy MapMergeStrategy) map[K]V {

	result := copyMap(original)

	for key, value := range candidate {
		found := hasKey(result, key)
		if found {
			if strategy&Override == Override {
				result[key] = value
				continue
			}
			if strategy&Remove == Remove {
				delete(result, key)
			}
		} else {
			if strategy&AddMissing == AddMissing {
				result[key] = value
			}
		}
	}

	return result
}

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

func hasKey[K comparable, V any](m map[K]V, key K) bool {
	_, exists := m[key]
	return exists
}
