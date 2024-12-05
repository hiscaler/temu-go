package is

import "regexp"

var millisecondPattern, packageNumberPattern, purchaseOrderNumberPattern, shipOrderNumberPattern *regexp.Regexp

func init() {
	millisecondPattern = regexp.MustCompile(`^[1-9][0-9]{12}$`)
	packageNumberPattern = regexp.MustCompile(`^(?i)pc[0-9]{13}$`)
	purchaseOrderNumberPattern = regexp.MustCompile(`^(?i)wb[0-9]{12,13}$`)
	shipOrderNumberPattern = regexp.MustCompile(`^(?i)fh[0-9]{13}$`)
}
