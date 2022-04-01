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
var projectHierarchySchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"projectid": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

var departmentHierarchySchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"departmentname": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"departmentid": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"projects": &schema.Schema{
			Type:     schema.TypeList,
			Elem:     projectHierarchySchema,
			Optional: true,
			
		},
	},
}

var groupHierarchySchema = map[string]*schema.Schema{
	"groups": &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"groupname": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"groupid": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"departments": &schema.Schema {
					Type: schema.TypeList,
					Elem: departmentHierarchySchema,
					Optional: true,
				},
			},
		},
	},
}

func resourceHierarchy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHierarchyUpdateOrCreate,
		ReadContext:   resourceHierarchyRead,
		UpdateContext: resourceHierarchyUpdateOrCreate,
		DeleteContext: resourceHierarchyDelete,
		Schema:        groupHierarchySchema,
	}
}

func resourceHierarchyUpdateOrCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*Client)

	groupItems := d.Get("groups").([]interface{})
	
	extractedGroupItems := []Group{}
	
	// Unpacking JSON data into structs
	for _, group := range groupItems {
		groupObject := group.(map[string]interface{})
		departmentStructs := []Department{}
		for _, department := range (groupObject["departments"].([]interface{})) {
			projectStructs := []Project{}
			departmentObject := department.(map[string]interface{})
			for _, project := range departmentObject["projects"].([]interface{}) {
				projectObject := project.(map[string]interface{})
				projectStruct := Project {
					ProjectID: projectObject["projectid"].(string),
				}
				projectStructs = append(projectStructs, projectStruct)
			}

			departmentStruct := Department {
				Projects: projectStructs,
				DepartmentID: departmentObject["departmentid"].(string),
				DepartmentName: departmentObject["departmentname"].(string),
			}

			departmentStructs = append(departmentStructs, departmentStruct)
			
		}
		groupStruct := Group{
			Departments: departmentStructs,
			GroupName: groupObject["groupname"].(string),
			GroupID: groupObject["groupid"].(string),
		}
		extractedGroupItems = append(extractedGroupItems, groupStruct)
	}

	_, err := client.postGroups("api/group_hierarchy", extractedGroupItems)
	if err != nil {
		return diag.FromErr(err)
	}

	// always run
	// Set to unix time to force the resource to reapply every time
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceHierarchyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Use this slice for detailed terraform diagnostics
	var diags diag.Diagnostics

	c := m.(*Client)
	groupHierarchy, err := c.getEndpointList("api/group_hierarchy")

	if err != nil || groupHierarchy == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error Retrieving groups",
			Detail:   fmt.Sprintf("groups %v", groupHierarchy),
		})

		return diag.FromErr(err)
	}

	if err := d.Set("groups", groupHierarchy); err != nil {
		return diag.FromErr(err)
	}

	// always run
	// Set to unix time to force the resource to reapply every time
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	
	return diags
}

func resourceHierarchyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceHierarchyRead(ctx, d, m)
}

func resourceHierarchyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
