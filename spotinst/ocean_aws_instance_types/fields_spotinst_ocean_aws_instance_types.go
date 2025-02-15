package ocean_aws_instance_types

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[Whitelist] = commons.NewGenericField(
		commons.OceanAWSInstanceTypes,
		Whitelist,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AWSClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			var result []string = nil
			if cluster.Compute != nil && cluster.Compute.InstanceTypes != nil &&
				cluster.Compute.InstanceTypes.Whitelist != nil {
				result = cluster.Compute.InstanceTypes.Whitelist
			}
			if err := resourceData.Set(string(Whitelist), result); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Whitelist), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AWSClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			if v, ok := resourceData.GetOk(string(Whitelist)); ok {
				if whitelist, err := expandInstanceTypeList(v); err != nil {
					return err
				} else {
					cluster.Compute.InstanceTypes.SetWhitelist(whitelist)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AWSClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			if v, ok := resourceData.GetOk(string(Whitelist)); ok {
				if whitelist, err := expandInstanceTypeList(v); err != nil {
					return err
				} else {
					cluster.Compute.InstanceTypes.SetWhitelist(whitelist)
				}
			} else {
				cluster.Compute.InstanceTypes.SetWhitelist(nil)
			}

			return nil
		},
		nil,
	)

	fieldsMap[Blacklist] = commons.NewGenericField(
		commons.OceanAWSInstanceTypes,
		Blacklist,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AWSClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			var result []string = nil
			if cluster.Compute != nil && cluster.Compute.InstanceTypes != nil &&
				cluster.Compute.InstanceTypes.Blacklist != nil {
				result = cluster.Compute.InstanceTypes.Blacklist
			}
			if err := resourceData.Set(string(Blacklist), result); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Blacklist), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AWSClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			if v, ok := resourceData.GetOk(string(Blacklist)); ok {
				if blacklist, err := expandInstanceTypeList(v); err != nil {
					return err
				} else {
					cluster.Compute.InstanceTypes.SetBlacklist(blacklist)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AWSClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			if v, ok := resourceData.GetOk(string(Blacklist)); ok {
				if blacklist, err := expandInstanceTypeList(v); err != nil {
					return err
				} else {
					cluster.Compute.InstanceTypes.SetBlacklist(blacklist)
				}
			} else {
				cluster.Compute.InstanceTypes.SetBlacklist(nil)
			}

			return nil
		},
		nil,
	)

	fieldsMap[Filters] = commons.NewGenericField(
		commons.OceanAWSInstanceTypes,
		Filters,
		&schema.Schema{
			Type:          schema.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{string(Blacklist), string(Whitelist)},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					string(Architectures): {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},

					string(Categories): {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},

					string(DiskTypes): {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},

					string(ExcludeFamilies): {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},

					string(ExcludeMetal): {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},

					string(Hypervisor): {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},

					string(IncludeFamilies): {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},

					string(IsEnaSupported): {
						Type:     schema.TypeBool,
						Optional: true,
					},

					string(MaxGpu): {
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(MaxMemoryGiB): {
						Type:     schema.TypeFloat,
						Optional: true,
					},

					string(MaxNetworkPerformance): {
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(MaxVcpu): {
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(MinEnis): {
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(MinGpu): {
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(MinMemoryGiB): {
						Type:     schema.TypeFloat,
						Optional: true,
					},

					string(MinNetworkPerformance): {
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(MinVcpu): {
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(RootDeviceTypes): {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},

					string(VirtualizationTypes): {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AWSClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			var result []interface{} = nil

			if cluster != nil && cluster.Compute != nil && cluster.Compute.InstanceTypes != nil &&
				cluster.Compute.InstanceTypes.Filters != nil {
				result = flattenFilters(cluster.Compute.InstanceTypes.Filters)
			}
			if len(result) > 0 {
				if err := resourceData.Set(string(Filters), result); err != nil {
					return fmt.Errorf(commons.FailureFieldReadPattern, string(Filters), err)
				}
			}
			return nil
		},

		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AWSClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			if v, ok := resourceData.GetOk(string(Filters)); ok {
				if filters, err := expandFilters(v, false); err != nil {
					return err
				} else {
					cluster.Compute.InstanceTypes.SetFilters(filters)
				}
			}
			return nil
		},

		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AWSClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			var value *aws.Filters = nil

			if v, ok := resourceData.GetOk(string(Filters)); ok {
				if filters, err := expandFilters(v, true); err != nil {
					return err
				} else {
					value = filters
				}
			}
			if cluster.Compute.InstanceTypes == nil {
				cluster.Compute.InstanceTypes = &aws.InstanceTypes{}
			}
			cluster.Compute.InstanceTypes.SetFilters(value)
			return nil
		},
		nil,
	)
}

func expandFilters(data interface{}, nullify bool) (*aws.Filters, error) {
	filters := &aws.Filters{}
	list := data.([]interface{})
	if list == nil || list[0] == nil {
		return filters, nil
	}
	m := list[0].(map[string]interface{})

	if v, ok := m[string(Architectures)]; ok {
		architectures, err := expandInstanceTypeFiltersList(v)
		if err != nil {
			return nil, err
		}
		if architectures != nil && len(architectures) > 0 {
			filters.SetArchitectures(architectures)
		} else {
			if nullify {
				filters.SetArchitectures(nil)
			}
		}
	}

	if v, ok := m[string(Categories)]; ok {
		categories, err := expandInstanceTypeFiltersList(v)
		if err != nil {
			return nil, err
		}
		if categories != nil && len(categories) > 0 {
			filters.SetCategories(categories)
		} else {
			if nullify {
				filters.SetCategories(nil)
			}
		}
	}

	if v, ok := m[string(DiskTypes)]; ok {
		diskTypes, err := expandInstanceTypeFiltersList(v)
		if err != nil {
			return nil, err
		}
		if diskTypes != nil && len(diskTypes) > 0 {
			filters.SetDiskTypes(diskTypes)
		} else {
			if nullify {
				filters.SetDiskTypes(nil)
			}
		}
	}

	if v, ok := m[string(ExcludeFamilies)]; ok {
		excludeFamilies, err := expandInstanceTypeFiltersList(v)
		if err != nil {
			return nil, err
		}
		if excludeFamilies != nil && len(excludeFamilies) > 0 {
			filters.SetExcludeFamilies(excludeFamilies)
		} else {
			if nullify {
				filters.SetExcludeFamilies(nil)
			}
		}
	}

	if v, ok := m[string(Hypervisor)]; ok {
		hypervisor, err := expandInstanceTypeFiltersList(v)
		if err != nil {
			return nil, err
		}
		if hypervisor != nil && len(hypervisor) > 0 {
			filters.SetHypervisor(hypervisor)
		} else {
			if nullify {
				filters.SetHypervisor(nil)
			}
		}
	}

	if v, ok := m[string(ExcludeMetal)].(bool); ok {
		filters.SetExcludeMetal(spotinst.Bool(v))
	}

	if v, ok := m[string(IncludeFamilies)]; ok {
		includeFamilies, err := expandInstanceTypeFiltersList(v)
		if err != nil {
			return nil, err
		}
		if includeFamilies != nil && len(includeFamilies) > 0 {
			filters.SetIncludeFamilies(includeFamilies)
		} else {
			if nullify {
				filters.SetIncludeFamilies(nil)
			}
		}
	}

	if v, ok := m[string(IsEnaSupported)].(bool); ok {
		filters.SetIsEnaSupported(spotinst.Bool(v))
	}

	if v, ok := m[string(MaxGpu)].(int); ok && v > 0 {
		filters.SetMaxGpu(spotinst.Int(v))
	}

	if v, ok := m[string(MaxMemoryGiB)].(float64); ok && v > 0 {
		filters.SetMaxMemoryGiB(spotinst.Float64(v))
	}

	if v, ok := m[string(MaxNetworkPerformance)].(int); ok && v > 0 {
		filters.SetMaxNetworkPerformance(spotinst.Int(v))
	}

	if v, ok := m[string(MaxVcpu)].(int); ok && v > 0 {
		filters.SetMaxVcpu(spotinst.Int(v))
	}

	if v, ok := m[string(MinEnis)].(int); ok && v > 0 {
		filters.SetMinEnis(spotinst.Int(v))
	}

	if v, ok := m[string(MinGpu)].(int); ok && v > 0 {
		filters.SetMinGpu(spotinst.Int(v))
	}

	if v, ok := m[string(MinMemoryGiB)].(float64); ok && v > 0 {
		filters.SetMinMemoryGiB(spotinst.Float64(v))
	}

	if v, ok := m[string(MinNetworkPerformance)].(int); ok && v > 0 {
		filters.SetMinNetworkPerformance(spotinst.Int(v))
	}

	if v, ok := m[string(MinVcpu)].(int); ok && v > 0 {
		filters.SetMinVcpu(spotinst.Int(v))
	}

	if v, ok := m[string(RootDeviceTypes)]; ok {
		rootDevicetypes, err := expandInstanceTypeFiltersList(v)
		if err != nil {
			return nil, err
		}
		if rootDevicetypes != nil && len(rootDevicetypes) > 0 {
			filters.SetRootDeviceTypes(rootDevicetypes)
		} else {
			if nullify {
				filters.SetRootDeviceTypes(nil)
			}
		}
	}

	if v, ok := m[string(VirtualizationTypes)]; ok {
		virtualizationtypes, err := expandInstanceTypeFiltersList(v)
		if err != nil {
			return nil, err
		}
		if virtualizationtypes != nil && len(virtualizationtypes) > 0 {
			filters.SetVirtualizationTypes(virtualizationtypes)
		} else {
			if nullify {
				filters.SetVirtualizationTypes(nil)
			}
		}
	}

	return filters, nil
}

func expandInstanceTypeList(data interface{}) ([]string, error) {
	list := data.([]interface{})
	result := make([]string, 0, len(list))

	for _, v := range list {
		if instanceTypeList, ok := v.(string); ok && instanceTypeList != "" {
			result = append(result, instanceTypeList)
		}
	}
	return result, nil
}

func expandInstanceTypeFiltersList(data interface{}) ([]string, error) {
	list := data.(*schema.Set).List()
	result := make([]string, 0, len(list))

	for _, v := range list {
		if instanceTypeList, ok := v.(string); ok && instanceTypeList != "" {
			result = append(result, instanceTypeList)
		}
	}
	return result, nil
}

func flattenFilters(filters *aws.Filters) []interface{} {
	var out []interface{}

	if filters != nil {
		result := make(map[string]interface{})

		result[string(ExcludeMetal)] = spotinst.BoolValue(filters.ExcludeMetal)
		result[string(IsEnaSupported)] = spotinst.BoolValue(filters.IsEnaSupported)
		result[string(MaxGpu)] = spotinst.IntValue(filters.MaxGpu)
		result[string(MinGpu)] = spotinst.IntValue(filters.MinGpu)
		result[string(MaxMemoryGiB)] = spotinst.Float64Value(filters.MaxMemoryGiB)
		result[string(MinMemoryGiB)] = spotinst.Float64Value(filters.MinMemoryGiB)
		result[string(MaxVcpu)] = spotinst.IntValue(filters.MaxVcpu)
		result[string(MinVcpu)] = spotinst.IntValue(filters.MinVcpu)
		result[string(MaxNetworkPerformance)] = spotinst.IntValue(filters.MaxNetworkPerformance)
		result[string(MinNetworkPerformance)] = spotinst.IntValue(filters.MinNetworkPerformance)
		result[string(MinEnis)] = spotinst.IntValue(filters.MinEnis)

		if filters.Architectures != nil {
			result[string(Architectures)] = filters.Architectures
		}

		if filters.Categories != nil {
			result[string(Categories)] = filters.Categories
		}

		if filters.DiskTypes != nil {
			result[string(DiskTypes)] = filters.DiskTypes
		}

		if filters.ExcludeFamilies != nil {
			result[string(ExcludeFamilies)] = filters.ExcludeFamilies
		}
		if filters.Hypervisor != nil {
			result[string(Hypervisor)] = filters.Hypervisor
		}

		if filters.IncludeFamilies != nil {
			result[string(IncludeFamilies)] = filters.IncludeFamilies
		}

		if filters.RootDeviceTypes != nil {
			result[string(RootDeviceTypes)] = filters.RootDeviceTypes
		}

		if filters.VirtualizationTypes != nil {
			result[string(VirtualizationTypes)] = filters.VirtualizationTypes
		}

		if len(result) > 0 {
			out = append(out, result)
		}
	}

	return out
}
