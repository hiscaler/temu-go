package temu

// 半托管专属服务
type semiManagedService struct {
	Order             semiOrderService             // 订单
	Logistics         semiLogisticsService         // 物流
	PlatformLogistics semiPlatformLogisticsService // 平台物流
	OnlineOrder       semiOnlineOrderService       // 在线下单
}

// 物流服务
type semiLogisticsService struct {
	ServiceProvider semiOnlineOrderLogisticsServiceProviderService // 服务商渠道
	Shipment        semiOnlineOrderLogisticsShipmentService        // 发货
	Warehouse       semiOnlineOrderLogisticsWarehouseService       // 仓库
}
