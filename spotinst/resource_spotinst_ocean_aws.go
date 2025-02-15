package spotinst

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/spotinst-sdk-go/spotinst/client"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/ocean_aws"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/ocean_aws_auto_scaling"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/ocean_aws_instance_types"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/ocean_aws_launch_configuration"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/ocean_aws_logging"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/ocean_aws_scheduling"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/ocean_aws_strategy"
)

func resourceSpotinstOceanAWS() *schema.Resource {
	setupClusterAWSResource()

	return &schema.Resource{
		CreateContext: resourceSpotinstClusterAWSCreate,
		ReadContext:   resourceSpotinstClusterAWSRead,
		UpdateContext: resourceSpotinstClusterAWSUpdate,
		DeleteContext: resourceSpotinstClusterAWSDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: commons.OceanAWSResource.GetSchemaMap(),
	}
}

func setupClusterAWSResource() {
	fieldsMap := make(map[commons.FieldName]*commons.GenericField)

	ocean_aws.Setup(fieldsMap)
	ocean_aws_auto_scaling.Setup(fieldsMap)
	ocean_aws_instance_types.Setup(fieldsMap)
	ocean_aws_launch_configuration.Setup(fieldsMap)
	ocean_aws_strategy.Setup(fieldsMap)
	ocean_aws_scheduling.Setup(fieldsMap)
	ocean_aws_logging.Setup(fieldsMap)

	commons.OceanAWSResource = commons.NewOceanAWSResource(fieldsMap)
}

func resourceSpotinstClusterAWSCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf(string(commons.ResourceOnCreate),
		commons.OceanAWSResource.GetName())

	cluster, err := commons.OceanAWSResource.OnCreate(resourceData, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID, err := createAWSCluster(resourceData, cluster, meta.(*Client))
	if err != nil {
		return diag.FromErr(err)
	}

	resourceData.SetId(spotinst.StringValue(clusterID))

	log.Printf("===> Cluster created successfully: %s <===", resourceData.Id())
	return resourceSpotinstClusterAWSRead(ctx, resourceData, meta)
}

func createAWSCluster(resourceData *schema.ResourceData, cluster *aws.Cluster, spotinstClient *Client) (*string, error) {
	if json, err := commons.ToJson(cluster); err != nil {
		return nil, err
	} else {
		log.Printf("===> Cluster create configuration: %s", json)
	}

	if v, ok := resourceData.Get(string(ocean_aws_launch_configuration.IAMInstanceProfile)).(string); ok && v != "" {
		// Wait for IAM instance profile to be ready.
		time.Sleep(10 * time.Second)
	}

	var resp *aws.CreateClusterOutput = nil
	err := resource.RetryContext(context.Background(), time.Minute, func() *resource.RetryError {
		input := &aws.CreateClusterInput{Cluster: cluster}
		r, err := spotinstClient.ocean.CloudProviderAWS().CreateCluster(context.Background(), input)
		if err != nil {
			// Checks whether we should retry cluster creation.
			if errs, ok := err.(client.Errors); ok && len(errs) > 0 {
				for _, err := range errs {
					if err.Code == "InvalidParamterValue" &&
						strings.Contains(err.Message, "Invalid IAM Instance Profile") {
						return resource.NonRetryableError(err)
					}
				}
			}
			// Some other error, report it.
			return resource.NonRetryableError(err)
		}
		resp = r
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to create cluster: %s", err)
	}
	return resp.Cluster.ID, nil
}

const ErrCodeClusterNotFound = "CLUSTER_DOESNT_EXIST"

func resourceSpotinstClusterAWSRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceOnRead),
		commons.OceanAWSResource.GetName(), id)

	input := &aws.ReadClusterInput{ClusterID: spotinst.String(id)}
	resp, err := meta.(*Client).ocean.CloudProviderAWS().ReadCluster(context.Background(), input)

	if err != nil {
		// If the cluster was not found, return nil so that we can show
		// that the cluster does not exist
		if errs, ok := err.(client.Errors); ok && len(errs) > 0 {
			for _, err := range errs {
				if err.Code == ErrCodeClusterNotFound {
					resourceData.SetId("")
					return nil
				}
			}
		}

		// Some other error, report it.
		return diag.Errorf("failed to read cluster: %s", err)
	}

	// if nothing was found, return no state
	clusterResponse := resp.Cluster
	if clusterResponse == nil {
		resourceData.SetId("")
		return nil
	}

	if err := commons.OceanAWSResource.OnRead(clusterResponse, resourceData, meta); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("===> Cluster read successfully: %s <===", id)
	return nil
}

func resourceSpotinstClusterAWSUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceOnUpdate),
		commons.OceanAWSResource.GetName(), id)

	shouldUpdate, changesRequiredRoll, tagsChanged, cluster, err := commons.OceanAWSResource.OnUpdate(resourceData, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if shouldUpdate {
		cluster.SetId(spotinst.String(id))
		if err := updateAWSCluster(cluster, resourceData, meta, changesRequiredRoll, tagsChanged); err != nil {
			return diag.FromErr(err)
		}
	}
	log.Printf("===> Cluster updated successfully: %s <===", id)
	return resourceSpotinstClusterAWSRead(ctx, resourceData, meta)
}

func updateAWSCluster(cluster *aws.Cluster, resourceData *schema.ResourceData, meta interface{}, changesRequiredRoll bool, tagsChanged bool) error {
	var input = &aws.UpdateClusterInput{
		Cluster: cluster,
	}

	var shouldRoll = false
	var conditionedRoll = false
	var autoApplyTags = false
	clusterID := resourceData.Id()
	if updatePolicy, exists := resourceData.GetOkExists(string(ocean_aws.UpdatePolicy)); exists {
		list := updatePolicy.([]interface{})
		if len(list) > 0 && list[0] != nil {
			m := list[0].(map[string]interface{})

			if roll, ok := m[string(ocean_aws.ShouldRoll)].(bool); ok && roll {
				shouldRoll = roll
			}

			if condRoll, ok := m[string(ocean_aws.ConditionedRoll)].(bool); ok && condRoll {
				conditionedRoll = condRoll
			}

			if aat, ok := m[string(ocean_aws.AutoApplyTags)].(bool); ok && aat {
				autoApplyTags = aat
			}
		}
	}

	if json, err := commons.ToJson(cluster); err != nil {
		return err
	} else {
		log.Printf("===> Cluster update configuration: %s", json)
	}

	if _, err := meta.(*Client).ocean.CloudProviderAWS().UpdateCluster(context.Background(), input); err != nil {
		return fmt.Errorf("[ERROR] Failed to update cluster [%v]: %v", clusterID, err)
	} else if shouldRoll {
		if !conditionedRoll || changesRequiredRoll || (!autoApplyTags && tagsChanged) {
			if err := rollOceanAWSCluster(resourceData, meta); err != nil {
				log.Printf("[ERROR] Cluster [%v] roll failed, error: %v", clusterID, err)
				return err
			}
		}
	} else {
		log.Printf("onRoll() -> Field [%v] is false, skipping cluster roll", string(ocean_aws.ShouldRoll))
	}

	return nil
}

func rollOceanAWSCluster(resourceData *schema.ResourceData, meta interface{}) error {
	clusterID := resourceData.Id()

	updatePolicy, exists := resourceData.GetOkExists(string(ocean_aws.UpdatePolicy))
	if !exists {
		return fmt.Errorf("ocean/aws: missing update policy for cluster %q", clusterID)
	}

	list := updatePolicy.([]interface{})
	if len(list) > 0 && list[0] != nil {
		updateClusterSchema := list[0].(map[string]interface{})

		rollConfig, ok := updateClusterSchema[string(ocean_aws.RollConfig)]
		if !ok || rollConfig == nil {
			return fmt.Errorf("ocean/aws: missing roll configuration, "+
				"skipping roll for cluster %q", clusterID)
		}

		rollSpec, err := expandOceanAWSClusterRollConfig(rollConfig, clusterID)
		if err != nil {
			return fmt.Errorf("ocean/aws: failed expanding roll "+
				"configuration for cluster %q, error: %v", clusterID, err)
		}

		rollJSON, err := commons.ToJson(rollConfig)
		if err != nil {
			return fmt.Errorf("ocean/aws: failed marshaling roll "+
				"configuration for cluster %q, error: %v", clusterID, err)
		}

		log.Printf("onRoll() -> Rolling cluster [%v] with configuration %s", clusterID, rollJSON)
		rollInput := &aws.CreateRollInput{Roll: rollSpec}
		if _, err = meta.(*Client).ocean.CloudProviderAWS().CreateRoll(context.TODO(), rollInput); err != nil {
			return fmt.Errorf("onRoll() -> Roll failed for cluster [%v], error: %v", clusterID, err)
		}
		log.Printf("onRoll() -> Successfully rolled cluster [%v]", clusterID)
	}

	return nil
}

func resourceSpotinstClusterAWSDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceOnDelete),
		commons.OceanAWSResource.GetName(), id)

	if err := deleteAWSCluster(resourceData, meta); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("===> Cluster deleted successfully: %s <===", resourceData.Id())
	resourceData.SetId("")
	return nil
}

func deleteAWSCluster(resourceData *schema.ResourceData, meta interface{}) error {
	clusterID := resourceData.Id()
	input := &aws.DeleteClusterInput{
		ClusterID: spotinst.String(clusterID),
	}

	if json, err := commons.ToJson(input); err != nil {
		return err
	} else {
		log.Printf("===> Cluster delete configuration: %s", json)
	}

	if _, err := meta.(*Client).ocean.CloudProviderAWS().DeleteCluster(context.Background(), input); err != nil {
		return fmt.Errorf("[ERROR] onDelete() -> Failed to delete cluster: %s", err)
	}
	return nil
}

func expandOceanAWSClusterRollConfig(data interface{}, clusterID string) (*aws.RollSpec, error) {
	list := data.([]interface{})
	spec := &aws.RollSpec{
		ClusterID: spotinst.String(clusterID),
	}

	if list != nil && list[0] != nil {
		m := list[0].(map[string]interface{})

		if v, ok := m[string(ocean_aws.BatchSizePercentage)].(int); ok {
			spec.BatchSizePercentage = spotinst.Int(v)
		}

		if v, ok := m[string(ocean_aws.LaunchSpecIDs)].([]string); ok {
			spec.LaunchSpecIDs = expandOceanAWSLaunchSpecIDs(v)
		}

		if v, ok := m[string(ocean_aws.BatchMinHealthyPercentage)].(int); ok && v > 0 {
			spec.BatchMinHealthyPercentage = spotinst.Int(v)
		}

		if v, ok := m[string(ocean_aws.RespectPDB)].(bool); ok {
			spec.RespectPDB = spotinst.Bool(v)
		}
	}

	return spec, nil
}

func expandOceanAWSLaunchSpecIDs(data interface{}) []string {
	list := data.([]interface{})
	result := make([]string, 0, len(list))

	for _, v := range list {
		if ls, ok := v.(string); ok && ls != "" {
			result = append(result, ls)
		}
	}

	return result
}
