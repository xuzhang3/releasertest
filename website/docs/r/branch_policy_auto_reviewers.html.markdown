---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_branch_policy_auto_reviewers"
description: |-
  Manages required reviewer policy branch policy within Azure DevOps project.
---

# azuredevops_branch_policy_auto_reviewers

Manages required reviewer policy branch policy within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "p" {
  name = "Sample Project"
}

resource "azuredevops_git_repository" "r" {
  project_id = azuredevops_project.p.id
  name       = "Sample Repo"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_user_entitlement" "user" {
  principal_name       = "mail@email.com"
  account_license_type = "basic"
}

resource "azuredevops_branch_policy_auto_reviewers" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
    auto_reviewer_ids  = [azuredevops_user_entitlement.user.id]
    submitter_can_vote = false
    message            = "Auto reviewer"
    path_filters       = ["*/src/*.ts"]

    scope {
      repository_id  = azuredevops_git_repository.r.id
      repository_ref = azuredevops_git_repository.r.default_branch
      match_type     = "Exact"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project in which the policy will be created.
- `enabled` - (Optional) A flag indicating if the policy should be enabled. Defaults to `true`.
- `blocking` - (Optional) A flag indicating if the policy should be blocking. Defaults to `true`.
- `settings` - (Required) Configuration for the policy. This block must be defined exactly once.

`settings` block supports the following:

- `auto_reviewer_ids` - (Required) Required reviewers ids. Supports multiples user Ids.
- `path_filters` - (Optional) Filter path(s) on which the policy is applied. Supports absolute paths, wildcards and multiple paths. Example: /WebApp/Models/Data.cs, /WebApp/* or *.cs,/WebApp/Models/Data.cs;ClientApp/Models/Data.cs.
- `submitter_can_vote` - (Optional) Controls whether or not the submitter's vote counts. Defaults to `false`.
- `message` - (Optional) Activity feed message, Message will appear in the activity feed of pull requests with automatically added reviewers.
- `scope` (Required) Controls which repositories and branches the policy will be enabled for. This block must be defined at least once.

  `scope` block supports the following:

  - `repository_id` - (Optional) The repository ID. Needed only if the scope of the policy will be limited to a single repository.
  - `repository_ref` - (Optional) The ref pattern to use for the match. If `match_type` is `Exact`, this should be a qualified ref such as `refs/heads/master`. If `match_type` is `Prefix`, this should be a ref path such as `refs/heads/releases`.
  - `match_type` (Optional) The match type to use when applying the policy. Supported values are `Exact` (default) or `Prefix`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of branch policy configuration.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/create?view=azure-devops-rest-5.1)

## Import

Azure DevOps Branch Policies can be imported using the project ID and policy configuration ID:

```sh
terraform import azuredevops_branch_policy_auto_reviewers.p aa4a9756-8a86-4588-86d7-b3ee2d88b033/60
```
