package mock

import "github.com/reverse-market/backend/pkg/database/models"

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func subslice(s1 []int, s2 []int) bool {
	if len(s1) > len(s2) {
		return false
	}
	for _, e := range s1 {
		if !contains(s2, e) {
			return false
		}
	}
	return true
}

func toIDs(tags []*models.TagInRequest) []int {
	ids := make([]int, len(tags))
	for i, tag := range tags {
		ids[i] = tag.ID
	}

	return ids
}
