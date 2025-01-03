package temu

// 半托管专属服务
type semiManagedService struct {
	Order       semiOrderService       // 订单
	OnlineOrder semiOnlineOrderService // 在线下单
}

// 在线下单服务
type semiOnlineOrderService struct {
	Logistics semiOnlineOrderLogisticsService // 物流服务
	Package   semiOnlineOrderPackageService   // 包裹服务
}

// 在线下单物流服务
type semiOnlineOrderLogisticsService struct {
	ServiceProvider semiOnlineOrderLogisticsServiceProviderService // 服务商
	Shipment        semiOnlineOrderLogisticsShipmentService        // 发货
	Warehouse       semiOnlineOrderLogisticsWarehouseService       // 仓库
}

// 在线下单包裹服务
type semiOnlineOrderPackageService struct {
	Unshipped semiOnlineOrderUnshippedPackageService // 未发货
	Shipped   semiOnlineOrderShippedPackageService   // 已发货
}
