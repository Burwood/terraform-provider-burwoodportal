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
			Description: "GCP project id",
		},
	},
}

var departmentDataSourceSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"departmentname": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			Description: "Department name as it appears in the portal.",
		},
		"departmentid": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			Description: "Unique department ID used under the hood to relate the department to projects and groups.",
		},
		"projects": &schema.Schema{
			Type:     schema.TypeList,
			Elem:     projectDataSourceSchema,
			Computed: true,
			Description: "List of projects. See project schema.",
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
					Description: "Group name as it appears in the portal.",
				},
				"groupid": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
					Description: "Unique group ID used under the hood to relate groups to departments.",
				},
				"departments": &schema.Schema {
					Type: schema.TypeList,
					Elem: departmentDataSourceSchema,
					Computed: true,
					Description: "List of departments. Projects are nested underneath departments. See department schema.",
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


