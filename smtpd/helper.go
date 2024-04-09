package smtpd

import (
	"fmt"

	"github.com/helloallan/terraform-provider-smtpd/sdk"
)

const mxRecord = "MX"

func flattenDNSRecords(dnsRecords *[]sdk.DNSRecord) map[string]map[string]string {
	records := make(map[string]map[string]string)

	for _, dnsRecord := range *dnsRecords {
		key := fmt.Sprintf("%s_%s", dnsRecord.DNSName, dnsRecord.DNSType)

		// Determine the appropriate record value format.
		value := dnsRecord.DNSValue
		if dnsRecord.DNSType == mxRecord {
			value = fmt.Sprintf("10 %s", value)
		}

		records[key] = map[string]string{
			"name":             dnsRecord.DNSName,
			"type":             dnsRecord.DNSType,
			"value":            value,
			"validation_state": dnsRecord.ValidationState,
		}
	}

	return records
}
