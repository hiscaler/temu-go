package helpers

import (
	"slices"

	"github.com/hiscaler/temu-go/entity"
)

// GetRegionByRegionId 根据区域 ID 获取所在区域
// https://seller.kuajingmaihuo.com/sop/view/231998342274104483#6mTvhA
func GetRegionByRegionId(regionId int) string {
	switch regionId {
	case 211, 219:
		return entity.AmericanRegion
	case 210, 76, 98, 13, 20, 32, 50, 52, 53, 54, 64, 68, 69, 72, 79, 90, 96, 108, 113, 114, 122, 141, 162, 163, 167, 180, 181, 186, 191, 91, 112, 151:
		return entity.EuropeanUnionRegion
	default:
		return entity.ChinaRegion
	}
}

// GetRegionBySiteId 根据站点 ID 获取所在区域
// https://seller.kuajingmaihuo.com/sop/view/231998342274104483#d78RUG
func GetRegionBySiteId(siteId int) string {
	if slices.Contains(entity.AmericanSiteIds, siteId) {
		return entity.AmericanRegion
	} else if slices.Contains(entity.EuropeanUnionSiteIds, siteId) {
		return entity.EuropeanUnionRegion
	} else {
		return entity.ChinaRegion
	}
}
