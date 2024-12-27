package temu

// 半托管专属服务
type semiManagedService struct {
	Order     semiOrderService     // 订单
	Logistics semiLogisticsService // 物流
}

// 物流服务
type semiLogisticsService struct {
	Warehouse semiLogisticsWarehouseService
}
