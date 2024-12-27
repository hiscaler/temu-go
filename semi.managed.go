package temu

// 半托管专属服务
type semiManagedService struct {
	Order             semiOrderService             // 订单
	Logistics         semiLogisticsService         // 物流
	PlatformLogistics semiPlatformLogisticsService // 平台物流
}

// 物流服务
type semiLogisticsService struct {
	ServiceProvider semiLogisticsServiceProviderService // 服务商渠道
	Shipment        semiLogisticsShipmentService        // 发货
	Warehouse       semiLogisticsWarehouseService       // 仓库
}
