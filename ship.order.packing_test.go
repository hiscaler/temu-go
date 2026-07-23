package temu

import (
	"strings"
	"testing"
	"time"

	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

// TestShipOrderPackingService_Match 测试装箱发货校验：按发货单号调用 Match，验证接口可正常返回且无错误
func TestShipOrderPackingService_Match(t *testing.T) {
	req := ShipOrderPackingMatchRequest{
		DeliveryOrderSnList: []string{"FH2408231977953"},
	}
	_, err := temuClient.Services.ShipOrder.Packing.Match(ctx, req)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.Packing.Match(ctx, %s)", jsonx.ToJson(req, "{}"))
}

// assertPredictVolumeCross 校验 PredictVolume 返回值的四种交叉情况：
//  1. 有体积 + 无错误 → 成功（wantErr=false 时要求；wantErr=nil 时允许）
//  2. 无体积 + 有错误 → 失败（wantErr=true 时要求；wantErr=nil 时允许）
//  3. 有体积 + 有错误 → 交叉异常，一律失败
//  4. 无体积 + 无错误 → 交叉异常，一律失败
func assertPredictVolumeCross(t *testing.T, volume string, err error, wantErr *bool) {
	t.Helper()
	hasVolume := strings.TrimSpace(volume) != ""
	hasErr := err != nil

	switch {
	case hasVolume && !hasErr:
		if wantErr != nil && *wantErr {
			assert.Failf(t, "期望失败却成功", "得到体积且无错误: volume=%q", volume)
		}
	case !hasVolume && hasErr:
		if wantErr != nil && !*wantErr {
			assert.Failf(t, "期望成功却失败", "得到错误且无体积: err=%v", err)
		}
	case hasVolume && hasErr:
		assert.Failf(t, "交叉异常：同时返回体积和错误", "volume=%q, err=%v", volume, err)
	default:
		assert.Fail(t, "交叉异常：既无体积也无错误")
	}
}

func boolPtr(v bool) *bool { return &v }

// TestShipOrderPackingService_PredictVolume 测试获取预估体积：覆盖参数边界与体积/错误四种交叉返回情况
func TestShipOrderPackingService_PredictVolume(t *testing.T) {
	fifty := make([]string, 50)
	for i := range fifty {
		fifty[i] = "FH2408231977953"
	}
	fiftyOne := make([]string, 51)
	for i := range fiftyOne {
		fiftyOne[i] = "FH2408231977953"
	}

	tests := []struct {
		name                 string
		deliveryOrderNumbers []string
		wantErr              *bool  // nil 表示仅校验交叉互斥，不限定成败
		wantErrContains      string // 非空时校验错误文案
		forbidErrContains    string // 非空时禁止错误文案（用于确认未命中数量校验）
	}{
		{
			name:                 "成功_单个发货单号_应有体积无错误",
			deliveryOrderNumbers: []string{"FH2607223990152"},
			wantErr:              boolPtr(false),
		},
		{
			name:                 "成功_多个发货单号_应有体积无错误",
			deliveryOrderNumbers: []string{"FH2408231977953", "FH2408231977953"},
			wantErr:              boolPtr(false),
		},
		{
			name:                 "边界_发货单号数量为50_仅校验交叉互斥且不得因数量校验失败",
			deliveryOrderNumbers: fifty,
			wantErr:              nil,
			forbidErrContains:    "发货单号列表数量须在 1-50 之间",
		},
		{
			name:                 "失败_发货单号为nil_应无体积有错误",
			deliveryOrderNumbers: nil,
			wantErr:              boolPtr(true),
			wantErrContains:      "发货单号列表数量须在 1-50 之间",
		},
		{
			name:                 "失败_发货单号为空切片_应无体积有错误",
			deliveryOrderNumbers: []string{},
			wantErr:              boolPtr(true),
			wantErrContains:      "发货单号列表数量须在 1-50 之间",
		},
		{
			name:                 "失败_发货单号数量为51_应无体积有错误",
			deliveryOrderNumbers: fiftyOne,
			wantErr:              boolPtr(true),
			wantErrContains:      "发货单号列表数量须在 1-50 之间",
		},
		{
			name:                 "失败_无效发货单号_应无体积有错误",
			deliveryOrderNumbers: []string{"INVALID-ORDER-SN"},
			wantErr:              boolPtr(true),
		},
		{
			name:                 "失败_空白发货单号_应无体积有错误",
			deliveryOrderNumbers: []string{"   "},
			wantErr:              boolPtr(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volume, err := temuClient.Services.ShipOrder.Packing.PredictVolume(ctx, tt.deliveryOrderNumbers...)
			assertPredictVolumeCross(t, volume, err, tt.wantErr)
			if tt.wantErrContains != "" {
				assert.ErrorContainsf(t, err, tt.wantErrContains, "错误文案应包含 %q", tt.wantErrContains)
			}
			if tt.forbidErrContains != "" && err != nil {
				assert.NotContainsf(t, err.Error(), tt.forbidErrContains, "不应命中数量校验错误: %v", err)
			}
			if tt.wantErr != nil && !*tt.wantErr {
				t.Logf("predictVolume: %s", volume)
			}
		})
	}
}

// TestShipOrderPackingService_SendForSelf 测试自送发货：对待装箱发货单组装自送信息后调用 Send，验证可成功发货
func TestShipOrderPackingService_SendForSelf(t *testing.T) {
	// 发货地址
	addresses, err := temuClient.Services.Mall.DeliveryAddress.Query(ctx)
	assert.Nilf(t, err, "temuClient.Services.Mall.DeliveryAddress.Query(ctx): error")
	assert.Equal(t, true, len(addresses) > 0, "temuClient.Services.Mall.DeliveryAddress.Query(ctx): results")
	address := addresses[0]

	params := ShipOrderQueryParams{
		Status: null.IntFrom(entity.ShipOrderStatusWaitingPacking),
	}
	params.PageSize = 1
	items, _, _, _, err := temuClient.Services.ShipOrder.Query(ctx, params)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.Query(ctx, %s)", jsonx.ToJson(params, "{}"))
	if len(items) != 0 {
		shipOrder := items[0]
		// 必须打印箱唛
		if !shipOrder.IsPrintBoxMark {
			_, err = temuClient.Services.Goods.Barcode.BoxMark(ctx, shipOrder.DeliveryOrderSn)
			assert.Nilf(t, err, "temuClient.Services.Goods.Barcode.BoxMark(ctx, %s)", shipOrder.DeliveryOrderSn)
		}

		driverName := shipOrder.DriverName
		if driverName == "" {
			driverName = "Zhang San"
		}
		req := ShipOrderPackingSendRequest{
			DeliveryAddressId:   address.ID,
			DeliveryOrderSnList: []string{shipOrder.DeliveryOrderSn},
			DeliverMethod:       null.IntFrom(entity.DeliveryMethodSelf),
			SelfDeliveryInfo: &ShipOrderPackingSendSelfDeliveryInformation{
				// DriverUid:             0,
				DriverName: driverName,
				// PlateNumber:           "",
				// DeliveryContactNumber: "",
				// DeliveryContactAreaNo: "",
				ExpressPackageNum: len(shipOrder.PackageList),
			},
		}
		_, err = temuClient.Services.ShipOrder.Packing.Send(ctx, req)
		assert.Nilf(t, err, "temuClient.Services.ShipOrder.Packing.Send(ctx, %s)", jsonx.ToJson(req, "{}"))
	} else {
		t.Logf("not found waitingPackage status purchase order")
	}
}

// TestShipOrderPackingService_SendForPlatformRecommendation 测试平台推荐物流发货：对待装箱发货单组装平台推荐配送信息后调用 Send，验证可成功发货
func TestShipOrderPackingService_SendForPlatformRecommendation(t *testing.T) {
	// 发货地址
	addresses, err := temuClient.Services.Mall.DeliveryAddress.Query(ctx)
	assert.Nil(t, err, "temuClient.Services.Mall.DeliveryAddress.Query(ctx): error")
	assert.Equal(t, true, len(addresses) > 0, "temuClient.Services.Mall.DeliveryAddress.Query(ctx): results")
	address := addresses[0]

	// 快递公司
	companies, err := temuClient.Services.Logistics.Companies(ctx)
	assert.Nilf(t, err, "temuClient.Services.Logistics.Companies(ctx): error")
	assert.Equal(t, true, len(companies) > 0, "temuClient.Services.Logistics.Companies(ctx): results")
	company := companies[0]

	params := ShipOrderQueryParams{
		Status: null.IntFrom(entity.ShipOrderStatusWaitingPacking),
	}
	params.PageSize = 1
	items, _, _, _, err := temuClient.Services.ShipOrder.Query(ctx, params)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.Query(ctx, %s)", jsonx.ToJson(params, "{}"))
	if len(items) != 0 {
		shipOrder := items[0]
		// 必须打印箱唛
		if !shipOrder.IsPrintBoxMark {
			_, err = temuClient.Services.Goods.Barcode.BoxMark(ctx, shipOrder.DeliveryOrderSn)
			assert.Nilf(t, err, "temuClient.Services.Goods.Barcode.BoxMark(ctx, %s)", shipOrder.DeliveryOrderSn)
		}

		d, _ := time.ParseInLocation(time.DateTime, time.Now().Format(time.DateOnly)+" 18:00:00", temuClient.TimeLocation)
		req := ShipOrderPackingSendRequest{
			DeliveryAddressId:   address.ID,
			DeliveryOrderSnList: []string{shipOrder.DeliveryOrderSn},
			DeliverMethod:       null.IntFrom(entity.DeliveryMethodPlatformRecommendation),
			ThirdPartyDeliveryInfo: &ShipOrderPackingSendPlatformRecommendationDeliveryInformation{
				ExpressCompanyId:          company.ShipId,
				TmsChannelId:              0,
				ExpressCompanyName:        company.ShipName,
				StandbyExpress:            false,
				ExpressDeliverySn:         shipOrder.ExpressDeliverySn,
				PredictTotalPackageWeight: shipOrder.PredictTotalPackageWeight,
				ExpectPickUpGoodsTime:     d.UnixMilli(),
				ExpressPackageNum:         len(shipOrder.PackageList),
				MinChargeAmount:           0.01,
				MaxChargeAmount:           0.02,
				PredictId:                 123, // ?
			},
		}
		if req.ThirdPartyDeliveryInfo.PredictTotalPackageWeight < 1000 {
			req.ThirdPartyDeliveryInfo.PredictTotalPackageWeight = 1000
		}
		_, err = temuClient.Services.ShipOrder.Packing.Send(ctx, req)
		assert.Nilf(t, err, "temuClient.Services.ShipOrder.Packing.Send(ctx, %s)", jsonx.ToJson(req, "{}"))
	} else {
		t.Logf("not found waitingPackage status purchase order")
	}
}
