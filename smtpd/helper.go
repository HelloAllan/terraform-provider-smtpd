package smtpd

import (
	"fmt"

	"github.com/helloallan/terraform-provider-smtpd/sdk"
)

const mxRecord = "MX"

func flattenDNSRecords(dnsRecords *[]sdk.DNSRecord) []interface{} {
	if dnsRecords != nil {
		records := make([]interface{}, len(*dnsRecords), len(*dnsRecords))

		for i, dnsRecord := range *dnsRecords {
			record := make(map[string]interface{})

			record["name"] = dnsRecord.DNSName
			record["type"] = dnsRecord.DNSType
			record["validation_state"] = dnsRecord.ValidationState

			if record["type"] == mxRecord {
				record["value"] = fmt.Sprintf("10 %s", dnsRecord.DNSValue)
			} else {
				record["value"] = dnsRecord.DNSValue
			}

			records[i] = record
		}

		return records
	}

	return make([]interface{}, 0)
}
