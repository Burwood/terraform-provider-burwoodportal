# Burwood portal terraform provider
This is a custom terraform provider written to interface with the cloudbilling portal REST API, meant to automate the configuration of projects in the portal.

## Examples
Examples of the objects defined by this provider are in the examples directory.

## Authentication
The portal REST API uses oauth flow. All you need to do is pass the configuration a username and password, and the provider will handle authentication from there.

An example of initialization of the provider with authentication is shown:

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
This data source does not require parameters, and will return an object with this structure:

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

This block will configure a new project in the portal and give it an initial active budget.
It will be automatically assigned to the given billing account.
If the project already exists, the existing project will be updated with any given fields.

This resource supports creation of a new budget. Simply define a subblock called latestbudget and pass it the fields shown.

A basic example of project configuration with a budget is shown here,


```
resource "burwoodportal_projects" "YOUR-GCP-PROJECT-ID-EXAMPLE1" { 
   projectid = "YOUR-GCP-PROJECT-ID-EXAMPLE1"
   departmentid  = "DEPARTMENTID" 
   recurringbudget = false
   latestbudget {
      ponumber = "12345"
      grant = "grantnum"
      amount = 1337
      state = "Active"
      billingaccountid = "ABCDEF-ABCDEF-ABCDEF" 
    }
}
```

