package temu

import (
	"encoding/json"
	"testing"
)

func TestGoodsCreateProductSemiManagedUsesDocumentedLanguageStrategyField(t *testing.T) {
	request := GoodsCreateRequest{
		ProductSemiManagedReq: &GoodsCreateProductSemiManaged{
			BindSiteIds:          []int{100},
			SemiLanguageStrategy: 1,
		},
		MaterialMultiLanguages: []string{"ja"},
	}

	payload, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal goods create request: %v", err)
	}

	var values map[string]any
	if err := json.Unmarshal(payload, &values); err != nil {
		t.Fatalf("unmarshal goods create request: %v", err)
	}
	semiManaged, ok := values["productSemiManagedReq"].(map[string]any)
	if !ok {
		t.Fatalf("missing semi managed request: %#v", values)
	}
	if semiManaged["semiLanguageStrategy"] != float64(1) {
		t.Fatalf("unexpected semi language strategy: %#v", semiManaged)
	}
	if _, exists := semiManaged["emiLanguageStrategy"]; exists {
		t.Fatalf("unexpected legacy language strategy field: %#v", semiManaged)
	}
	if languages, ok := values["materialMultiLanguages"].([]any); !ok || len(languages) != 1 || languages[0] != "ja" {
		t.Fatalf("unexpected material languages: %#v", values["materialMultiLanguages"])
	}
}
