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
			Description: "PO to use for this budget (custom terminology for this field may be present in the portal UI, e.g. ChartField)",
		},
		"grant": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Description: "Grant to use for this budget.",
		},
		"amount": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
			Description: "Dollar amount to use for the budget. Acts as a float data type (decimals allowed).",
		},
		"billingaccountid": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			Description: "GCP billing account ID to use for consumption on this budget.",
		},
		"expirationdate": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Description: "YYYY-MM-DD format. Date after which to mark the budget as consumed regardless of spend on it.",
		},
		"dateissued": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			Description: "YYYY-MM-DD format. Budget issue date. Used in budget alerting emails.",
		},
		"dateactivated": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			Description: "Budget activation date. Date on which the budget activate its billing account and tracking consumption.",
		},
		"datesuspended": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			Description: "Date on which the budget was deactivate and marked consumed.",
		},
		"state": &schema.Schema {
			Type:	schema.TypeString,
			Default: "Future",
			Optional: true,
			Description: "Valid values are 'Active' and 'Future'. WARNING! If set to 'Active', this budget will mark existing active budgets as consumed and set the GCP project's billing account to the specified billingaccountid!",
		},
		"recurring": &schema.Schema {
			Type:	schema.TypeBool,
			Default: false,
			Optional: true,
			Description: "Boolean; whether the budget should be a recurring monthly budget or a standard budget.",
		},
	},
}



var projectSchema = map[string]*schema.Schema{
	"projectid": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: "GCP Project ID",
	},
	"projectname": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Description: "Project name  as shown in the portal.",
	},
	"primarycontactemail": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Description: "The project primary contact email address.",
	},
	"billingcontactemail": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Description: "Primary billing contact email.",
	},
	"aftercredits": &schema.Schema{
		Type:     schema.TypeString,
		Default: "Suspend",
		Optional: true,
		Description: "Valid values: 'Bill' or 'Suspend'. Only set to bill if post-budget free spend is desired.",
	},
	"aftercreditsaccount": &schema.Schema {
		Type: schema.TypeString,
		Optional: true,
		Description: "GCP billing account to use for post-credit consumption. Only applies if aftercredits is set to 'Suspend' ",
	},
	"aftercreditspo": &schema.Schema {
		Type: schema.TypeString,
		Optional: true,
		Description: "Purchase Order for afterCredits consumption.",
	},
	"paidbillingaccount": &schema.Schema {
		Type: schema.TypeString,
		Optional: true,
		Description: "The project GCP billing account ID. WARNING! This will change the project's billing account in GCP!",
	},
	"totalbudget": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: "Total budget dollar amount on the project.",
	},
	"recurringbudget": &schema.Schema {
		Type: schema.TypeBool,
		Default: false,
		Optional: true,
		Description: "Boolean. Whether project budgets should recur on a monthly basis.",
	},
	"departmentid": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: "Department ID to for the project. If invalid or not given, the project will be placed into the 'Unaffiliated Projects' department. Department names and ID's can be seen in the group hierarchy data source, and a code example can be seen in the guides.",
	},
	"departmentname": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: "Department name that the project is under.",
	},
	"latestbudget": &schema.Schema {
		Type: schema.TypeList,
		Elem: budgetSchema,
		Optional: true,
		Description: "Most recently added budget. If given as a subblock, a new budget will be appended to the prjoect. See the budget schema for more details.",
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
