# Barebons config file for dev purposes
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


resource "burwoodportal_projects" "acme-corp-export" { 
   projectid = "acme-corp-export"
   aftercredits = "Bill"
   aftercreditsaccount = "ACCOUNTID"
   departmentid  = flatten([
    for groupIndex, _ in data.burwoodportal_hierarchy.hierarchy.groups : [
      for departmentIndex, departmentValue in data.burwoodportal_hierarchy.hierarchy.groups[groupIndex].departments : {
          departmentid = departmentValue.departmentid
      } 
    if departmentValue.departmentname == "CloudOps"  ] 
  ])[0].departmentid
   recurringbudget = true
   latestbudget {
      ponumber = "MYPO123"
      grant = "MYGRANT456"
      amount = 13337
      state = "Future"
      recurring = true
      billingaccountid = "ACCOUNTID"
    }
}
