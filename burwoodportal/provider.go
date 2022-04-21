package burwoodportal

import (
	"context"
	//"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//"log"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default: "https://api.bcs.burwood.com",
				Description: "Desired hostname. Only needed if interactions with non-production environments are desired.",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PORTAL_USERNAME", nil),
				Description: "Burwood portal username used for authentication with the Burwood portal REST API.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("PORTAL_PASSWORD", nil),
				Description: "Burwood portal password used for authentication with the Burwood portal REST API.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"burwoodportal_projects": resourceProject(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"burwoodportal_hierarchy":   dataSourceGroupHierarchy(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Set up and return an authenticated client for the Burwood portal api.

	// Credential config
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	//  Host config
	var host *string
	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// Initialize the cilent or handle errors
	if (username != "") && (password != "") {
		c, err := NewClient(host, &username, &password)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Burwood client",
				Detail:   "Unable to authenticate user for authenticated Burwood client",
			})

			return nil, diags
		}

		return c, diags
	}

	c, err := NewClient(host, nil, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Burwood client",
			Detail:   "Unable to create anonymous Burwood client",
		})
		return nil, diags
	}

	return c, diags

}
