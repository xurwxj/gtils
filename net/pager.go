package net

import (
	"strconv"
	"strings"
)

// PaginationFromMap generate pagination from map
func PaginationFromMap(m map[string]string, pSize, defaultSize int64) (map[string]interface{}, int, int) {
	rs := make(map[string]interface{})
	for k, v := range m {
		rs[k] = v
	}
	var page = 1
	var pageSize = defaultSize
	if pSize > 0 {
		pageSize = pSize
	}
	qPSize, has := m["pageSize"]
	if has {
		qPSizeV, qPSerr := strconv.ParseInt(qPSize, 10, 64)
		if qPSerr == nil {
			pageSize = qPSizeV
		}
	}
	if pageSize < 1 {
		pageSize = 10
	}
	rs["pageSize"] = pageSize
	qPage, pHas := m["page"]
	if pHas {
		qPageV, qPErr := strconv.Atoi(qPage)
		if qPErr == nil {
			page = qPageV
		}
	}
	if page < 1 {
		page = 1
	}
	rs["page"] = page
	rs["offset"] = (page - 1) * int(pageSize)
	// return map[string]interface{}{"page": page, "pageSize": pageSize, "offset": (page - 1) * int(pageSize)}, int(page), int(pageSize)
	return rs, int(page), int(pageSize)
}

// QueryBytesToMap query bytes convert to map
// s := "A=B&C=D&E=F"
func QueryBytesToMap(query []byte) map[string]string {
	m := make(map[string]string)
	queries := strings.Split(string(query), "&")
	for _, q := range queries {
		q = strings.TrimSpace(q)
		if q != "" {
			zs := strings.Split(q, "=")
			if len(zs) == 2 {
				k := strings.TrimSpace(zs[0])
				v := strings.TrimSpace(zs[1])
				if k != "" && v != "" {
					m[k] = v
				}
			}
		}
	}
	return m
}

// PaginationInfo generate pagination result info object
func PaginationInfo(page, pageSize, defaultSize, total int) map[string]int {
	if pageSize == 0 {
		pageSize = defaultSize
	}
	if pageSize == 0 {
		pageSize = 10
	}
	if page == 0 {
		page = 1
	}
	totalPage := total / pageSize
	if total%pageSize > 0 {
		totalPage++
	}
	return map[string]int{"page": page, "pageSize": pageSize, "pages": totalPage, "total": total}
}
