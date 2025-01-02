package temu

// 在线下单服务
type semiOnlineOrderService struct {
	Logistics semiOnlineOrderLogisticsService
}

// 在线下单物流服务
type semiOnlineOrderLogisticsService struct {
	ServiceProvider semiOnlineOrderLogisticsServiceProviderService // 服务商
	Shipment        semiOnlineOrderLogisticsShipmentService        // 发货
	Warehouse       semiOnlineOrderLogisticsWarehouseService       // 仓库
}
