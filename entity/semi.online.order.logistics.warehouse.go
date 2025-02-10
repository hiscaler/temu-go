package entity

type SemiOnlineOrderLogisticsWarehouse struct {
	WarehouseId      string `json:"warehouseId"`      // 仓库 Id
	WarehouseName    string `json:"warehouseName"`    // 仓库名称
	RegionId1        int    `json:"regionId1"`        // 所属经营站点
	DefaultWarehouse bool   `json:"defaultWarehouse"` // 是否默认仓库
}
