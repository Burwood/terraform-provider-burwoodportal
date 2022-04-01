
terraform {
  required_providers {
    burwoodportal = {
      version = "0.2.0"
      source  = "burwood.com/portal/burwoodportal"
    }
  }
}

# Enter portal credentials as inputs
variable "username" {
  description = "Burwood portal username"
  type        = string
} 

variable "password" {
  description = "Burwood portal password"
  type        = string
}


provider "burwoodportal" {
    # Pass a username and password when initializing the provider.
    # Needed to make authenticated requests to the portal.
    username = var.username
    password = var.password
}

data "burwoodportal_hierarchy" "hierarchy" {}


# This block will configure a new project in the portal and give it an initial active budget.
# It will be automatically assigned to the given billing account.
# If the project already exists, 
resource "burwoodportal_projects" "YOUR-GCP-PROJECT-ID-EXAMPLE1" { 
   projectid = "YOUR-GCP-PROJECT-ID-EXAMPLE1"
   aftercredits = "Suspend"  # The after credits behavior will be set to "Suspend" by default, or can be given explicitly like here.
   departmentid  = "DEPARTMENTID"# This is a unique identifier for the desired department for the project. 
   recurringbudget = false
   # This will create a new budget or append a new budget to the project being configured.
   # If set to "Active", the new budget WILL replace the current active budget and change the billing account on the project!
   latestbudget {
      ponumber = "12345"
      grant = "grantnum"
      amount = 1337
      state = "Active"
      billingaccountid = "ABCDEF-ABCDEF-ABCDEF" # This must be a valid billing account configured in the portal for your organization.
    }
}

# Same as above, except this project will be set to bill after credits and given a future budget.
# It is recommended to always give a project an active billing account, either by specifying an active budget
# or by explicitly specifying an active account.
# Beware of this setup--it will allow the project to spend freely!
resource "burwoodportal_projects" "YOUR-GCP-PROJECT-EXAMPLE2" { 
   projectid = "YOUR-GCP-PROJECT-EXAMPLE2"
   paidbillingaccount = "ABCDEF-ABCDEF-ABCDEF" # This must be a valid billing account configured in the portal for your organization.
   aftercredits = "Bill" 
   aftercreditsaccount = "ABCDEF-ABCDEF-ABCDEF" 

   # This example pulls the department id for a department called "Test Department" using the hierarchy.
   departmentid  = flatten([
    for groupIndex, _ in data.burwoodportal_hierarchy.hierarchy.groups : [
      for departmentIndex, departmentValue in data.burwoodportal_hierarchy.hierarchy.groups[groupIndex].departments : {
          departmentid = departmentValue.departmentid
      } 
    if departmentValue.departmentname == "Test Department"  ] 
  ])[0].departmentid
  
   recurringbudget = false
   # This will create a new budget or append a new budget to the project being configured.
   latestbudget {
      ponumber = "MYPO123"
      grant = "MYGRANT456"
      amount = 13337
      state = "Future"
      billingaccountid = "ABCDEF-ABCDEF-ABCDEF"
    }
}







