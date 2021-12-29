package collection

func DistinctInt64(slc []int64) []int64 {
	var result []int64
	tempMap := map[int64]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}
