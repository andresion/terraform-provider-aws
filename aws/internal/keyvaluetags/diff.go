package keyvaluetags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service"
)

// Meta is the interface implemented by the CRUD handlers' `meta` parameter.
// This is a local copy of the definition from `internal/provider/meta` to avoid circular imports.
type Meta interface {
	// GetDefaultTagsConfig returns the provider's `default_tags` configuration.
	GetDefaultTagsConfig() *DefaultConfig

	// GetIgnoreTagsConfig returns the provider's `ignore_tags` configuration.
	GetIgnoreTagsConfig() *IgnoreConfig

	// GetServicePackage returns the ServicePackage for the specified service ID.
	GetServicePackage(id string) service.ServicePackage
}

func fromMeta(v interface{}) (*DefaultConfig, *IgnoreConfig) {
	return v.(Meta).GetDefaultTagsConfig(), v.(Meta).GetIgnoreTagsConfig()
}

// SetTagsDiff sets the new plan difference with the result of
// merging resource tags on to those defined at the provider-level;
// returns an error if unsuccessful or if the resource tags are identical
// to those configured at the provider-level to avoid non-empty plans
// after resource READ operations as resource and provider-level tags
// will be indistinguishable when returned from an AWS API.
func SetTagsDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	defaultTagsConfig, ignoreTagsConfig := fromMeta(meta)

	resourceTags := New(diff.Get("tags").(map[string]interface{}))

	if defaultTagsConfig.TagsEqual(resourceTags) {
		return fmt.Errorf(`"tags" are identical to those in the "default_tags" configuration block of the provider: please de-duplicate and try again`)
	}

	allTags := defaultTagsConfig.MergeTags(resourceTags).IgnoreConfig(ignoreTagsConfig)

	// To ensure "tags_all" is correctly computed, we explicitly set the attribute diff
	// when the merger of resource-level tags onto provider-level tags results in n > 0 tags,
	// otherwise we mark the attribute as "Computed" only when their is a known diff (excluding an empty map)
	// or a change for "tags_all".
	// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/18366
	// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/19005
	if len(allTags) > 0 {
		if err := diff.SetNew("tags_all", allTags.Map()); err != nil {
			return fmt.Errorf("error setting new tags_all diff: %w", err)
		}
	} else if len(diff.Get("tags_all").(map[string]interface{})) > 0 {
		if err := diff.SetNewComputed("tags_all"); err != nil {
			return fmt.Errorf("error setting tags_all to computed: %w", err)
		}
	} else if diff.HasChange("tags_all") {
		if err := diff.SetNewComputed("tags_all"); err != nil {
			return fmt.Errorf("error setting tags_all to computed: %w", err)
		}
	}

	return nil
}
