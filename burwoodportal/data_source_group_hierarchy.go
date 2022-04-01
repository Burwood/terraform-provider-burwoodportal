package burwoodportal

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)


// Because this endpoint is a nesting doll,
// break up the schema into list types.
var projectDataSourceSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"projectid": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}

var departmentDataSourceSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"departmentname": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"departmentid": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"projects": &schema.Schema{
			Type:     schema.TypeList,
			Elem:     projectDataSourceSchema,
			Computed: true,
		},
	},
}

var groupDataSourceSchema = map[string]*schema.Schema{
	"groups": &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"groupname": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"groupid": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"departments": &schema.Schema {
					Type: schema.TypeList,
					Elem: departmentDataSourceSchema,
					Computed: true,
				},
			},
		},
	},
}

// Define group schema.
func dataSourceGroupHierarchy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcehierarchyRead,
		Schema:      groupDataSourceSchema,
	}
}

// Read groups from the portal API endpoint.
func dataSourcehierarchyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Use this slice for detailed terraform diagnostics
	var diags diag.Diagnostics

	c := m.(*Client)
	groups, err := c.getEndpointList("api/group_hierarchy")

	if err != nil || groups == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error Retrieving Groups",
			Detail:   fmt.Sprintf("Groups %v", groups),
		})

		return diag.FromErr(err)
	}

	if err := d.Set("groups", groups); err != nil {
		return diag.FromErr(err)
	}

	// always run
	// Set to unix time to force the resource to reapply every time
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}


