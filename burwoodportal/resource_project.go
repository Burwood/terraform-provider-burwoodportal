package burwoodportal

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"encoding/json"
	"net/http"
	"strings"
)

var budgetSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"ponumber": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"grant": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"amount": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"billingaccountid": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"expirationdate": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"dateissued": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"dateactivated": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"datesuspended": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": &schema.Schema {
			Type:	schema.TypeString,
			Default: "Future",
			Optional: true,
		},
		"recurring": &schema.Schema {
			Type:	schema.TypeBool,
			Default: true,
			Optional: true,
		},
	},
}



var projectSchema = map[string]*schema.Schema{
	"projectid": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"projectname": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"primarycontactemail": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"billingcontactemail": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"aftercredits": &schema.Schema{
		Type:     schema.TypeString,
		Default: "Suspend",
		Optional: true,
	},
	"aftercreditsaccount": &schema.Schema {
		Type: schema.TypeString,
		Optional: true,
	},
	"aftercreditspo": &schema.Schema {
		Type: schema.TypeString,
		Optional: true,
	},
	"paidbillingaccount": &schema.Schema {
		Type: schema.TypeString,
		Optional: true,
	},
	"totalbudget": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
	"recurringbudget": &schema.Schema {
		Type: schema.TypeBool,
		Default: false,
		Optional: true,
	},
	"departmentid": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"departmentname": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
	"latestbudget": &schema.Schema {
		Type: schema.TypeList,
		Elem: budgetSchema,
		Optional: true,
	},
}
 

func resourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceProjectRead,
		UpdateContext: resourceProjectCreateOrUpdate,
		DeleteContext: resourceProjectDelete,
		CreateContext: resourceProjectCreateOrUpdate, 
		Schema:      projectSchema,
	}
}

func (c *Client) postProject(projectID string, postBody Project) (*Project, error) {
	postBodyMarshaled, err := json.Marshal(postBody)
	if err != nil {
		return nil, err
	}

	processedBody := strings.NewReader(string(postBodyMarshaled))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/project/%s", c.HostURL, projectID), processedBody)
	if err != nil {
		return nil, err
	}
	

	responseBody, err := c.doRequest(req, nil)

	// Unmarshal response JSON into a map data structure
	responseBodyUnmarshal := &Project{}
	err = json.Unmarshal(responseBody, &responseBodyUnmarshal)
	if responseBodyUnmarshal == nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return responseBodyUnmarshal, nil
}


func resourceProjectCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics { 
	var diags diag.Diagnostics

	c := m.(*Client)
	projectID := d.Get("projectid").(string)
	projectStruct := Project {
		ProjectID: projectID,
		ProjectName: d.Get("projectname").(string),
		PrimaryContactEmail: d.Get("primarycontactemail").(string),
		BillingContactEmail: d.Get("billingcontactemail").(string),
		AfterCredits: d.Get("aftercredits").(string),
		AfterCreditsAccount: d.Get("aftercreditsaccount").(string),
		AfterCreditsPO: d.Get("aftercreditspo").(string),
		PaidBillingAccount: d.Get("paidbillingaccount").(string),
		TotalBudget: d.Get("totalbudget").(string),
		RecurringBudget: d.Get("recurringbudget").(bool),
		DepartmentID: d.Get("departmentid").(string),
	}

	response, err := c.postProject(projectID, projectStruct)

	if err != nil || response == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error Creating Project",
			Detail:   fmt.Sprintf("Project: %v", projectStruct),
		})

		return diag.FromErr(err)
	}


	// Latest allowance comes in as a list,
	// but that's only so we can use the TypeList schema
	// for the allowance.
	allowanceList := d.Get("latestbudget").([]interface{})

	if (len(allowanceList) > 1) {
		// Undefined behavior for multiple allowances.
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error With latestbudget",
			Detail:   fmt.Sprintf("Can't specify multiple budgets."),
		})

		return diag.FromErr(err)
	} else if (len(allowanceList) == 1) {
		allowanceObject := allowanceList[0].(map[string]interface{})
		allowanceStruct := Allowance {
			PONumber: allowanceObject["ponumber"].(string),
			Grant: allowanceObject["grant"].(string),
			Amount: allowanceObject["amount"].(int),
			BillingAccountID: allowanceObject["billingaccountid"].(string),
			ExpirationDate: allowanceObject["expirationdate"].(string),
			DateSuspended: allowanceObject["datesuspended"].(string),
			DateActivated: allowanceObject["dateactivated"].(string),
			State: allowanceObject["state"].(string),
			Recurring: allowanceObject["recurring"].(bool),
		}
		err = c.postBudget(projectID, "project", allowanceStruct)
	
		if err != nil  {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error Creating Budget",
				Detail:   fmt.Sprintf("Budget: %v", allowanceStruct),
			})
	
			return diag.FromErr(err)
		}
	} 

	d.SetId(projectID)

	return diags
}



func (c *Client) postBudget(entityID string, scope string, postBody Allowance) (error) {
	postBodyMarshaled, err := json.Marshal(postBody)
	if err != nil {
		return err
	}

	processedBody := strings.NewReader(string(postBodyMarshaled))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/%s/%s/add_budget", c.HostURL, scope, entityID), processedBody)
	if err != nil {
		return err
	}
	

	_, err = c.doRequest(req, nil)

	if err != nil {
		return err
	}

	return nil
}



func (c *Client) getProject(projectID string) (*Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/project/%s", c.HostURL, projectID), nil)
	if err != nil {
		return nil, err
	}
	
	responseBody, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}
	responseBodyUnmarshal := &Project{}
	err = json.Unmarshal(responseBody, &responseBodyUnmarshal)
	if responseBodyUnmarshal == nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return responseBodyUnmarshal, nil
}

func (c *Client) getLatestProjectBudget(projectID string) (*Allowance, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/project/%s/budgets", c.HostURL, projectID), nil)
	if err != nil {
		return nil, err
	}
	
	responseBody, err := c.doRequest(req, nil)

	if err != nil {
		return nil, err
	}

	responseBodyUnmarshal := []Allowance{}
	err = json.Unmarshal(responseBody, &responseBodyUnmarshal)
	if responseBodyUnmarshal == nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	// Get the most recently configured budget object.
	// Should be the last element in the JSON response.
	latestBudgetObject := &Allowance{}
	if len(responseBodyUnmarshal) > 0 {
		latestBudgetObject = &responseBodyUnmarshal[len(responseBodyUnmarshal) - 1]
	}
	return latestBudgetObject, nil
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Use this slice for detailed terraform diagnostics
	var diags diag.Diagnostics

	c := m.(*Client)

	projectID := d.Get("projectid")
	projectObject, err := c.getProject(projectID.(string))

	if err != nil || projectObject == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error Retrieving Project",
			Detail:   fmt.Sprintf("Project: %v", projectObject),
		})

		return diag.FromErr(err)
	}

	d.SetId(projectID.(string))
	d.Set("projectid", projectID.(string))
	d.Set("projectname", projectObject.ProjectName)
	d.Set("primarycontactemail", projectObject.PrimaryContactEmail)
	d.Set("billingcontactemail", projectObject.BillingContactEmail)
	d.Set("aftercredits", projectObject.AfterCredits)
	d.Set("aftercreditsaccount", projectObject.AfterCreditsAccount)
	d.Set("aftercreditspo", projectObject.AfterCreditsPO)
	d.Set("paidbillingaccount", projectObject.PaidBillingAccount)
	d.Set("totalbudget", projectObject.TotalBudget)
	d.Set("recurringbudget", projectObject.RecurringBudget)
	d.Set("departmentid", projectObject.DepartmentID)
	d.Set("departmentname", projectObject.DepartmentName)


	if projectObject.ProjectID != "" { // Handling the case where no such project exists
		budgetObject, err := c.getLatestProjectBudget(projectID.(string))
		d.Set("latestbudget", budgetObject)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error Retrieving Latest Budget",
				Detail:   fmt.Sprintf("Project: %v", projectObject),
			})
			return diag.FromErr(err)
		}
	}
	

	return diags
}


func (c *Client) deleteProject(projectID string) (*Project, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/project/%s", c.HostURL, projectID), nil)
	if err != nil {
		return nil, err
	}
	
	responseBody, err := c.doRequest(req, nil)
	responseBodyUnmarshal := &Project{}
	err = json.Unmarshal(responseBody, &responseBodyUnmarshal)
	if err != nil {
		return nil, err
	}

	return responseBodyUnmarshal, nil
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics { 
	var diags diag.Diagnostics
	c := m.(*Client)
	projectID := d.Get("projectid")
	deleteResponse, err := c.deleteProject(projectID.(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error Deleting Project",
			Detail:   fmt.Sprintf("Response: %v", deleteResponse),
		})

		return diag.FromErr(err)
	}
	
	return diags
}
