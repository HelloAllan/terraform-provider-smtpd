package sdk

import "encoding/json"

// DNSRecord struct
type DNSRecord struct {
	DNSName         string `json:"name"`
	DNSType         string `json:"type"`
	DNSValue        string `json:"value"`
	ValidationState string `json:"validation_state"`
}

// CreateProfile struct
type CreateProfile struct {
	ProfileName               string `json:"profile_name"`
	BounceDomain              string `json:"bounce_domain",omitempty`
	LinkDomain                string `json:"link_domain",omitempty`
	SendingDomain             string `json:"sending_domain",omitempty`
	LinkDomainDefaultRedirect string `json:"link_domain_default_redirect",omitempty`
}

// Profile struct
type Profile struct {
	ProfileId                 string      `json:"profile_id"`
	ProfileName               string      `json:"profile_name"`
	State                     string      `json:"state"`
	LinkDomainDefaultRedirect string      `json:"link_domain_default_redirect",omitempty`
	BounceDomain              string      `json:"bounce_domain",omitempty`
	LinkDomain                string      `json:"link_domain",omitempty`
	SendingDomain             string      `json:"sending_domain",omitempty`
	CreatedAtUtc              int64       `json:"created_at_utc",omitempty`
	ModifiedAtUtc             int64       `json:"modified_at_utc",omitempty`
	DNSRecords                []DNSRecord `json:"dns_records",omitempty`
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
	data, err := json.Marshal(&CreateProfile{
		p.ProfileName,
		p.BounceDomain,
		p.LinkDomain,
		p.SendingDomain,
		p.LinkDomainDefaultRedirect,
	})
	if err != nil {
		return "", err
	}
	return string(data), nil
}
