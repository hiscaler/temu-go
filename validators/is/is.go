package is

import "regexp"

var millisecondPattern,
	packageNumberPattern,
	originalPurchaseOrderNumberPattern,
	purchaseOrderNumberPattern,
	shipOrderNumberPattern,
	mobilePhoneNumberPattern,
	telNumberPattern,
	telNumberAreaCodePattern,
	imageUrlPattern *regexp.Regexp

func init() {
	millisecondPattern = regexp.MustCompile(`^[1-9][0-9]{12}$`)
	packageNumberPattern = regexp.MustCompile(`^(?i)pc[0-9]{13,15}$`)
	originalPurchaseOrderNumberPattern = regexp.MustCompile(`^(?i)wp[0-9]{12,15}$`)
	purchaseOrderNumberPattern = regexp.MustCompile(`^(?i)wb[0-9]{12,15}$`)
	shipOrderNumberPattern = regexp.MustCompile(`^(?i)fh[0-9]{13,15}$`)
	mobilePhoneNumberPattern = regexp.MustCompile(`^(?:\+?86)?1(?:3\d{3}|5[^4\D]\d{2}|8\d{3}|7(?:[0-35-9]\d{2}|4(?:0\d|1[0-2]|9\d))|9[0-35-9]\d{2}|6[2567]\d{2}|4(?:(?:10|4[01])\d{3}|[68]\d{4}|[579]\d{2}))\d{6}$`)
	telNumberPattern = regexp.MustCompile(`^(0\d{2,3}-?)?\d{8}$`)
	telNumberAreaCodePattern = regexp.MustCompile(`^0(\d{2,3})$`)
	imageUrlPattern = regexp.MustCompile(`^https?://.*\.(?:png|jpg|jpeg)(?:\?.*)?$`)
}
