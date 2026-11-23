package temu

import "testing"

func TestGoodsImageUploadRequestValidateAcceptsImageDataURL(t *testing.T) {
	request := GoodsImageUploadRequest{
		Image: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQIHWP4z8DwHwAFgAI/ScLw5wAAAABJRU5ErkJggg==",
	}

	if err := request.validate(); err != nil {
		t.Fatalf("validate image data URL: %v", err)
	}
}

func TestGoodsImageUploadRequestValidateRejectsInvalidImageDataURL(t *testing.T) {
	request := GoodsImageUploadRequest{Image: "data:image/png;base64,not-base64"}

	if err := request.validate(); err == nil {
		t.Fatal("expected invalid image data URL to fail validation")
	}
}
