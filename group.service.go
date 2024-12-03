package temu

// productService 商品服务
type productService struct {
	Data          goodsService              // 商品基础数据
	Brand         goodsBrandService         // 商品品牌数据
	LifeCycle     goodsLifeCycleService     // 商品生命周期数据
	TopSelling    goodsTopSellingService    // 畅销商品数据
	Sales         goodsSalesService         // 销售数据
	Certification goodsCertificationService // 资质服务
}
