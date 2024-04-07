package smtpd

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/helloallan/terraform-provider-smtpd/sdk"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SMTPD_HOST", nil),
			},
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SMTPD_API_KEY", nil),
			},
			"api_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SMTPD_API_SECRET", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"smtpd_profile": resourceProfile(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"smtpd_profile": dataSourceProfile(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	key := d.Get("api_key").(string)
	secret := d.Get("api_secret").(string)

	var host *string
	var version *string

	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (key != "") && (secret != "") {
		c, err := sdk.NewClient(host, &key, &secret, version)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create smtpd client",
				Detail:   "Unable to authenticate user for authenticated smtpd client",
			})

			return nil, diags
		}

		return c, diags
	}

	c, err := sdk.NewClient(host, nil, nil, version)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create smtpd client",
			Detail:   "Unable to create anonymous smtpd client",
		})
		return nil, diags
	}

	return c, diags
}
