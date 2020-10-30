// +build all resource_branchpolicy_acceptance_test policy
// +build !exclude_resource_branchpolicy_acceptance_test !exclude_policy

package acceptancetests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccBranchPolicyMinReviewers_CreateAndUpdate(t *testing.T) {
	minReviewerTfNode := "azuredevops_branch_policy_min_reviewers.p"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getMinReviewersHcl(true, true, 1, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(minReviewerTfNode, "id"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "blocking", "true"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "enabled", "true"),
				),
			}, {
				Config: getMinReviewersHcl(false, false, 2, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(minReviewerTfNode, "id"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "blocking", "false"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "enabled", "false"),
				),
			}, {
				ResourceName:      minReviewerTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(minReviewerTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getMinReviewersHcl(enabled bool, blocking bool, reviewers int, submitterCanVote bool) string {
	settings := fmt.Sprintf(
		`
		reviewer_count     = %d
		submitter_can_vote = %t
		`, reviewers, submitterCanVote,
	)

	return getBranchPolicyHcl("azuredevops_branch_policy_min_reviewers", enabled, blocking, settings)
}

func TestAccBranchPolicyAutoReviewers_CreateAndUpdate(t *testing.T) {
	autoReviewerTfNode := "azuredevops_branch_policy_auto_reviewers.p"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, &[]string{"AZDO_TEST_AAD_USER_EMAIL"}) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getAutoReviewersHcl(true, true, false, "auto reviewer", fmt.Sprintf("\"%s\",\"%s\"", "*/API*.cs", "README.md")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(autoReviewerTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "blocking", "true"),
				),
			}, {
				Config: getAutoReviewersHcl(false, false, true, "new auto reviewer", fmt.Sprintf("\"%s\",\"%s\"", "*/API*.cs", "README.md")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(autoReviewerTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "blocking", "false"),
				),
			}, {
				ResourceName:      autoReviewerTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(autoReviewerTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getAutoReviewersHcl(enabled bool, blocking bool, submitterCanVote bool, message string, pathFilters string) string {
	settings := fmt.Sprintf(
		`
		auto_reviewer_ids  = [azuredevops_user_entitlement.user.id]
		submitter_can_vote = %t
		message 		   = "%s"
		path_filters       = [%s]
		`, submitterCanVote, message, pathFilters,
	)
	userPrincipalName := os.Getenv("AZDO_TEST_AAD_USER_EMAIL")
	userEntitlement := testutils.HclUserEntitlementResource(userPrincipalName)

	return strings.Join(
		[]string{
			userEntitlement,
			getBranchPolicyHcl("azuredevops_branch_policy_auto_reviewers", enabled, blocking, settings),
		},
		"\n",
	)
}

func TestAccBranchPolicyBuildValidation_CreateAndUpdate(t *testing.T) {
	buildValidationTfNode := "azuredevops_branch_policy_build_validation.p"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getBuildValidationHcl(true, true, "build validation", 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(buildValidationTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.filename_patterns.#", "3"),
				),
			}, {
				Config: getBuildValidationHcl(false, false, "build validation rename", 720),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(buildValidationTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.filename_patterns.#", "3"),
				),
			}, {
				ResourceName:      buildValidationTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(buildValidationTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getBuildValidationHcl(enabled bool, blocking bool, displayName string, validDuration int) string {
	settings := fmt.Sprintf(
		`
		display_name = "%s"
		valid_duration = %d
		build_definition_id = azuredevops_build_definition.build.id
		filename_patterns =  [
			"/WebApp/*",
			"!/WebApp/Tests/*",
			"*.cs"
		]
		`, displayName, validDuration,
	)

	return getBranchPolicyHcl("azuredevops_branch_policy_build_validation", enabled, blocking, settings)
}

func TestAccBranchPolicyWorkItemLinking_CreateAndUpdate(t *testing.T) {
	resourceName := "azuredevops_branch_policy_work_item_linking"
	workItemLinkingTfNode := fmt.Sprintf("%s.p", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getBranchPolicyHcl(resourceName, true, true, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(workItemLinkingTfNode, "enabled", "true"),
				),
			}, {
				Config: getBranchPolicyHcl(resourceName, false, false, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(workItemLinkingTfNode, "enabled", "false"),
				),
			}, {
				ResourceName:      workItemLinkingTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(workItemLinkingTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBranchPolicyCommentResolution_CreateAndUpdate(t *testing.T) {
	resourceName := "azuredevops_branch_policy_comment_resolution"
	workItemLinkingTfNode := fmt.Sprintf("%s.p", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getBranchPolicyHcl(resourceName, true, true, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(workItemLinkingTfNode, "enabled", "true"),
				),
			}, {
				Config: getBranchPolicyHcl(resourceName, false, false, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(workItemLinkingTfNode, "enabled", "false"),
				),
			}, {
				ResourceName:      workItemLinkingTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(workItemLinkingTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getBranchPolicyHcl(resourceName string, enabled bool, blocking bool, settings string) string {
	branchPolicy := fmt.Sprintf(`
	resource "%s" "p" {
		project_id = azuredevops_project.project.id
		enabled  = %t
		blocking = %t
		settings {
			%s
			scope {
				repository_id  = azuredevops_git_repository.repository.id
				repository_ref = azuredevops_git_repository.repository.default_branch
				match_type     = "exact"
			}
		}
	}
	`, resourceName, enabled, blocking, settings)
	projectAndRepo := testutils.HclGitRepoResource(testutils.GenerateResourceName(), testutils.GenerateResourceName(), "Clean")
	buildDef := testutils.HclBuildDefinitionResource(
		"Sample Build Definition",
		`\\`,
		"TfsGit",
		"${azuredevops_git_repository.repository.id}",
		"master",
		"path/to/yaml",
		"")

	return strings.Join(
		[]string{
			branchPolicy,
			projectAndRepo,
			buildDef,
		},
		"\n",
	)
}
