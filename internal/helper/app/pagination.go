package app

import (
	"TTPanel/pkg/convert"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) int {
	pageGet := convert.StrTo(c.Query("page")).MustInt()
	pagePost := convert.StrTo(c.PostForm("page")).MustInt()
	if pageGet <= 0 && pagePost <= 0 {
		return 1
	}
	if pageGet > 0 {
		return pageGet
	}
	if pagePost > 0 {
		return pagePost
	}
	return 1
}

func GetOffsetLimits(page, limit int) (offset, limits int) {
	return (page - 1) * limit, limit
}
