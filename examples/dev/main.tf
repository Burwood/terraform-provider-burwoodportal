# Barebons config file for dev purposes
terraform {
  required_providers {
    burwoodportal = {
      version = "0.2.0"
      source  = "burwood.com/portal/burwoodportal"
    }
  }
}

# # Enter portal credentials as inputs
# variable "username" {
#   description = "Burwood portal username"
#   type        = string
# } 

# variable "password" {
#   description = "Burwood portal password"
#   type        = string
# }


provider "burwoodportal" {
    # Pass a username and password when initializing the provider.
    # Needed to make authenticated requests to the portal.
    host = "http://localhost:5000"
    username = "dpalencia@burwood.com"
    password = "ExXgsHrkZG34Uv5"
}


data "burwoodportal_hierarchy" "hierarchy" {}


resource "burwoodportal_projects" "acme-corp-export" { 
   projectid = "acme-corp-export"
   aftercredits = "Bill"
   aftercreditsaccount = "018213-2CBF7F-0EBB8B"
   departmentid  = flatten([
    for groupIndex, _ in data.burwoodportal_hierarchy.hierarchy.groups : [
      for departmentIndex, departmentValue in data.burwoodportal_hierarchy.hierarchy.groups[groupIndex].departments : {
          departmentid = departmentValue.departmentid
      } 
    if departmentValue.departmentname == "CloudOps"  ] 
  ])[0].departmentid
   latestbudget {
      ponumber = "MYPO123"
      grant = "MYGRANT456"
      amount = 133337
      state = "Future"
      billingaccountid = "018213-2CBF7F-0EBB8B"
    }
}
