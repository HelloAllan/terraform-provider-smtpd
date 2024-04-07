package smtpd

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/helloallan/terraform-provider-smtpd/sdk"
)

func resourceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProfileCreate,
		ReadContext:   resourceProfileRead,
		UpdateContext: resourceProfileUpdate,
		DeleteContext: resourceProfileDelete,
		Schema: map[string]*schema.Schema{
			"profile_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Required: true,
				ForceNew: true,
			},
			"bounce_domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"link_domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sending_domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dns_records": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
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

func resourceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	profileName := d.Get("profile_name").(string)
	profile, err := c.GetProfileByName(profileName)
	if err != nil {
		profile, err = c.CreateProfile(&sdk.Profile{
			ProfileName:               d.Get("profile_name").(string),
			SendingDomain:             d.Get("sending_domain").(string),
			LinkDomainDefaultRedirect: d.Get("link_domain_default_redirect").(string),
			LinkDomain:                d.Get("link_domain").(string),
			BounceDomain:              d.Get("bounce_domain").(string),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(profile.ProfileId)
	d.Set("sending_domain", profile.SendingDomain)
	d.Set("link_domain_default_redirect", profile.LinkDomainDefaultRedirect)
	d.Set("link_domain", profile.LinkDomain)
	d.Set("bounce_domain", profile.BounceDomain)

	resourceProfileRead(ctx, d, m)

	return diags
}

func resourceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	profileId := d.Id()

	profileSetup, err := c.GetProfile(profileId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("profile_name", profileSetup.ProfileName)
	d.Set("profile_id", profileSetup.ProfileId)
	d.Set("state", profileSetup.State)

	dnsRecords := flattenDNSRecords(&profileSetup.DNSRecords)
	if err := d.Set("dns_records", dnsRecords); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sdk.Client)

	profileId := d.Id()

	if d.HasChange("state") {
		err := c.NoOp(profileId)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceProfileRead(ctx, d, m)
}

func resourceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	profileId := d.Id()

	err := c.NoOp(profileId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
