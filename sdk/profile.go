package sdk

import "encoding/json"

// DNSRecord struct
type DNSRecord struct {
	DNSName         string "json:name"
	DNSType         string "json:type"
	DNSValue        string "json:value"
	ValidationState string "json:validation_state"
}

// Profile struct
type Profile struct {
	ProfileId                 string      `json:profile_id`
	ProfileName               string      `json:profile_name`
	LinkDomainDefaultRedirect string      `json:link_domain_default_redirect`
	BounceDomain              string      `json:bounce_domain`
	LinkDomain                string      `json:link_domain`
	SendingDomain             string      `json:sending_domain`
	CreatedAtUtc              int64       `json:created_at_utc`
	ModifiedAtUtc             int64       `json:modified_at_utc`
	DNSRecords                []DNSRecord `json:dns_records`
}

// LoadFromJSON update object from json
func (p *Profile) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &p)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (p *Profile) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
