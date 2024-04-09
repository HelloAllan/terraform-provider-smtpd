package smtpd

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/helloallan/terraform-provider-smtpd/sdk"
)

func dataSourceProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProfileRead,
		Schema: map[string]*schema.Schema{
			"profile_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"profile_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"link_domain_default_redirect": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bounce_domain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"link_domain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"sending_domain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_records": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"validation_state": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	profileName := d.Get("profile_name").(string)

	profile, err := c.GetProfileByName(profileName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(profile.ProfileId)
	d.Set("profile_id", profile.ProfileId)
	d.Set("profile_name", profile.ProfileName)
	d.Set("state", profile.State)
	d.Set("link_domain_default_redirect", profile.LinkDomainDefaultRedirect)
	d.Set("bounce_domain", profile.BounceDomain)
	d.Set("link_domain", profile.LinkDomain)
	d.Set("sending_domain", profile.SendingDomain)

	profileSetup, err := c.GetProfile(profile.ProfileId)
	if err != nil {
		return diag.FromErr(err)
	}

	dnsRecords := flattenDNSRecords(&profileSetup.DNSRecords)
	if err := d.Set("dns_records", dnsRecords); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
