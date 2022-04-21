# Burwood portal terraform provider
This is a custom terraform provider written to interface with the cloudbilling portal REST API and to automate the configuration of projects in the portal.

## Examples
Anotated code examples can be seen [here](https://github.com/Burwood/terraform-provider-burwoodportal/blob/8074a7796936dfb16eeda8f2bbb2c1e41d724c23/examples/provider-overview/projects-and-budgets.tf).

## REST API
This provider consumes the Burwood portal's public REST API, as documented here:

https://app.swaggerhub.com/apis-docs/Burwood-Group/burwood_cloud_services/

## Authentication
The portal REST API uses oauth flow. Pass the provide configuration a username and password and it will handle authentication with the REST API from there.

An example of initialization of the provider with input variables is shown:

```
variable "username" {
  description = "Burwood portal username"
  type        = string
} 

variable "password" {
  description = "Burwood portal password"
  type        = string
}


provider "burwoodportal" {
    username = var.username
    password = var.password
}

```

## Data Sources 

### burwoodportal_hierarchy
This data source does not require parameters, and will output an object with this structure:

```

list(object({
    groupname = string
    groupid = string
    departments = list(object({
       departmentname = string
       departmentid = string
       projects = list(object({
          projectname = string
          projectid = string
       }))
    }))
  }))

```

The main purpose of this data source object is to serve as a way to retrieve relations between groups, departments, and projects as well as retrieve the unique IDs for these objects.

## Resources

### burwoodportal_projects

This resource will configure a new project in the portal and give it an initial active budget.
It will be automatically assigned to the given billing account.
If the project already exists, the existing project will be updated with any given fields.

This resource supports creation of a new budget. Simply define a subblock called latestbudget and pass it the fields shown.

A basic example of project configuration with an active budget is shown here,


```
resource "burwoodportal_projects" "YOUR-GCP-PROJECT-ID-EXAMPLE1" { 
   projectid = "YOUR-GCP-PROJECT-ID-EXAMPLE1"
   departmentid  = "DEPARTMENTID" 
   latestbudget {
      ponumber = "12345"
      grant = "grantnum"
      amount = 1337
      state = "Active"
      billingaccountid = "ABCDEF-ABCDEF-ABCDEF" 
    }
}
```

## burwoodportal_projects Inputs
| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| projectid | Project id. Should match the project ID in GCP. | `string` | n/a | yes |
| projectname | Name of the project in the portal. | `string` | n/a | no |
| primarycontactemail |  Project primary contact email. | `string` | n/a | no |
| billingcontactemail |  Project billing contact email. | `string` | n/a | no |
| aftercredits |  After credits behavior. VALID VALUES: 'Suspend' or 'Bill'. Only set to 'Bill' if you want the billing account to spend freely! | `string` | Suspend | no |
| aftercreditsaccount |  Billing account ID for after credits spend. Only relevant to projects set to aftercredits 'Bill'.  | `string` | n/a | no |
| aftercreditspo | PO to assign after credits spend to. | `string` | n/a | no |
| recurringbudget | Whether the budget should be recurring. Boolean true or false.| `boolean` | false | no |
| departmentid | ID of the department to place the project into. The ID can be retrieved via the hierarchy object. [Example here](https://github.com/Burwood/burwood-portal-public-trerraform-provider/blob/8b17e282fa08920820d2cec1aa6da32f5385f23c/examples/provider-overview/projects-and-budgets.tf)| `string` | n/a | yes |

## latestbudget Inputs
*** note that this is a nested block inside Portal Resource config
| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| ponumber | PO number to assign to this budget. | `string` | n/a | no |
| grant | Grant number to assign to this budget. | `string` | n/a | no |
| amount|  Budget amount in dollars. | `number` | n/a | yes |
| billingaccountid |  Billing account ID to consume this budget on. | `string` | n/a | yes |
| expirationdate | YYYY-MM-DD date to expire the budget on. | `string` | n/a | no |
| state| Desired state of the budget. Either 'Active' to activate the budget immediately, or 'Future' to set a future budget. WARNING: Setting an active budget will change the project's active billing account! | `string` | Future | yes |

