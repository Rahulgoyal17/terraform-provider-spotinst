package stateful_node_azure_extension

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	azurev3 "github.com/spotinst/spotinst-sdk-go/service/stateful/providers/azure"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[Extensions] = commons.NewGenericField(
		commons.StatefulNodeAzureExtensions,
		Extensions,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(Publisher): {
						Type:     schema.TypeString,
						Required: true,
					},

					string(APIVersion): {
						Type:     schema.TypeString,
						Required: true,
					},

					string(MinorVersionAutoUpgrade): {
						Type:     schema.TypeBool,
						Required: true,
					},

					string(Name): {
						Type:     schema.TypeString,
						Required: true,
					},

					string(Type): {
						Type:     schema.TypeString,
						Required: true,
					},

					string(ProtectedSettings): {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
					},

					string(PublicSettings): {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			stWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			st := stWrapper.GetStatefulNode()
			var result []interface{} = nil
			if st != nil && st.Compute != nil && st.Compute.LaunchSpecification != nil && st.Compute.LaunchSpecification.Extensions != nil {
				extensions := st.Compute.LaunchSpecification.Extensions
				result = flattenExtensions(extensions)
			}

			if result != nil {
				if err := resourceData.Set(string(Extensions), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Extensions), err)
				}
			}

			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			stWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			st := stWrapper.GetStatefulNode()
			var value []*azurev3.Extension = nil

			if v, ok := resourceData.GetOk(string(Extensions)); ok {
				var extensions []*azurev3.Extension

				if st != nil && st.Compute != nil && st.Compute.LaunchSpecification != nil {
					if st.Compute.LaunchSpecification.Extensions != nil {
						extensions = st.Compute.LaunchSpecification.Extensions
					}

					if ext, err := expandExtensions(v, extensions); err != nil {
						return err
					} else {
						value = ext
					}

					st.Compute.LaunchSpecification.SetExtensions(value)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			stWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			st := stWrapper.GetStatefulNode()
			var value []*azurev3.Extension = nil

			if v, ok := resourceData.GetOk(string(Extensions)); ok {
				//create new image object in case st did not get it from previous import step.
				var extensions []*azurev3.Extension

				if st != nil && st.Compute != nil && st.Compute.LaunchSpecification != nil {

					if st.Compute.LaunchSpecification.Extensions != nil {
						extensions = st.Compute.LaunchSpecification.Extensions
					}

					if extensions, err := expandExtensions(v, extensions); err != nil {
						return err
					} else {
						value = extensions
					}

					st.Compute.LaunchSpecification.SetExtensions(value)
				}
			}
			return nil
		},
		nil,
	)
}

func flattenExtensions(extensions []*azurev3.Extension) []interface{} {
	result := make([]interface{}, 0, len(extensions))

	for _, extension := range extensions {
		m := make(map[string]interface{})
		m[string(APIVersion)] = spotinst.StringValue(extension.APIVersion)
		m[string(Name)] = spotinst.StringValue(extension.Name)
		m[string(Publisher)] = spotinst.StringValue(extension.Publisher)
		m[string(Type)] = spotinst.StringValue(extension.Type)
		m[string(MinorVersionAutoUpgrade)] = spotinst.BoolValue(extension.MinorVersionAutoUpgrade)
		m[string(ProtectedSettings)] = extension.ProtectedSettings
		m[string(PublicSettings)] = extension.PublicSettings

		result = append(result, m)
	}
	return result
}

func expandExtensions(data interface{}, extensions []*azurev3.Extension) ([]*azurev3.Extension, error) {
	list := data.(*schema.Set).List()

	for _, v := range list {
		ext, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		extension := &azurev3.Extension{}

		if v, ok := ext[string(APIVersion)].(string); ok && v != "" {
			extension.SetAPIVersion(spotinst.String(v))
		}
		if v, ok := ext[string(Name)].(string); ok && v != "" {
			extension.SetName(spotinst.String(v))
		}
		if v, ok := ext[string(Publisher)].(string); ok && v != "" {
			extension.SetPublisher(spotinst.String(v))
		}
		if v, ok := ext[string(Type)].(string); ok && v != "" {
			extension.SetType(spotinst.String(v))
		}
		if v, ok := ext[string(MinorVersionAutoUpgrade)].(bool); ok {
			extension.SetMinorVersionAutoUpgrade(spotinst.Bool(v))
		}
		if v, ok := ext[string(ProtectedSettings)].(map[string]interface{}); ok {
			extension.SetProtectedSettings(v)
		}
		if v, ok := ext[string(PublicSettings)].(map[string]interface{}); ok {
			extension.SetPublicSettings(v)
		}

		extensions = append(extensions, extension)
	}

	return extensions, nil
}
