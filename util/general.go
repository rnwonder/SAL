package util

import (
	"cmp"
	"fmt"
	"net/url"
	"strconv"
)

func CalculatePageInfo(page string, limit string, total int) (int, int, int, int, int) {
	pageString := cmp.Or(page, "1")
	limitString := cmp.Or(limit, "10")
	pageInt, _ := strconv.Atoi(pageString)
	limitInt, _ := strconv.Atoi(limitString)
	startIndex := (pageInt - 1) * limitInt
	endIndex := pageInt * limitInt
	totalPages := total / limitInt

	// Account for remainder
	if total%limitInt > 0 {
		totalPages++
	}

	if totalPages < 1 {
		totalPages = 1
	}

	if endIndex > total {
		endIndex = total
	}
	return startIndex, endIndex, totalPages, limitInt, pageInt
}

func NextPage(page int, totalPages int) string {
	if page >= totalPages {
		return strconv.Itoa(totalPages)
	}
	return strconv.Itoa(page + 1)
}

func PrevPage(page int) string {
	if page <= 1 {
		return "1"
	}
	return strconv.Itoa(page - 1)
}

func EncodeMapToString(data map[string]interface{}) string {
	values := url.Values{}
	for key, value := range data {
		values.Add(key, fmt.Sprintf("%v", value))
	}
	return values.Encode()
}
