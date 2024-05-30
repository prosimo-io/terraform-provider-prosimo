package prosimo

import (
	"context"
	"fmt"
	"log"
	"os"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

func flattenStringList(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func validateFilePath(v interface{}, k string) (warnings []string, errors []error) {
	operation := v.(string)
	log.Printf("[DEBUG] Validating Path %s", operation)
	if _, err := os.Stat(operation); os.IsNotExist(err) {
		return returnError("Path does not exist", fmt.Errorf("[ERROR] Invalid Path"))
	}
	return nil, nil
}
func returnError(message string, err error) (warnings []string, errors []error) {
	var errorVar []error
	var warningVar []string
	return append(warningVar, message), append(errorVar, err)
}

func retryUntilTaskComplete(ctx context.Context, d *schema.ResourceData, meta interface{}, taskID string) resource.RetryFunc {
	prosimoClient := meta.(*client.ProsimoClient)
	return func() *resource.RetryError {
		getTaskStatus, err := prosimoClient.GetTaskStatus(ctx, taskID)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if getTaskStatus.TaskDetails.Status == "IN-PROGRESS" {
			return resource.RetryableError(fmt.Errorf("task %s is not completed yet", taskID))
		} else if getTaskStatus.TaskDetails.Status == "FAILURE" {
			for _, subtask := range getTaskStatus.ItemList {
				if subtask.Status == "FAILURE" {
					return resource.NonRetryableError(fmt.Errorf("task %s has failed at step %s", taskID, subtask.Name))
				}
			}

		}
		return nil
	}
}

func retryUntilTaskComplete_edge_onboard(ctx context.Context, d *schema.ResourceData, meta interface{}, taskID string) resource.RetryFunc {
	prosimoClient := meta.(*client.ProsimoClient)
	return func() *resource.RetryError {
		getTaskStatus, err := prosimoClient.GetTaskStatus(ctx, taskID)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if getTaskStatus.TaskDetails.Status == "IN-PROGRESS" {
			return resource.RetryableError(fmt.Errorf("task %s is not completed yet", taskID))
		} else if getTaskStatus.TaskDetails.Status == "FAILURE" {
			for _, subtask := range getTaskStatus.ItemList {
				if subtask.Status == "FAILURE" {
					log.Printf("[ERROR]: task %s has failed at step %s, rolling back", taskID, subtask.Name)
				}
			}
			log.Println("[DEBUG]: decommisioning edge with id", d.Id())
			appResponseData, err := prosimoClient.DeleteAppDeployment(ctx, d.Id())
			if err != nil {
				return resource.NonRetryableError(err)
			}
			log.Printf("[DEBUG] Waiting for task id %s to complete", appResponseData.ResourceData.ID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskComplete_edge_Offboard(ctx, d, meta, appResponseData.ResourceData.ID))
			if err != nil {
				return resource.NonRetryableError(err)
			}
			log.Printf("[INFO] task %s is successful", appResponseData.ResourceData.ID)

		}
		return nil
	}
}

func retryUntilTaskComplete_edge_Offboard(ctx context.Context, d *schema.ResourceData, meta interface{}, taskID string) resource.RetryFunc {
	prosimoClient := meta.(*client.ProsimoClient)
	return func() *resource.RetryError {
		getTaskStatus, err := prosimoClient.GetTaskStatus(ctx, taskID)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if getTaskStatus.TaskDetails.Status == "IN-PROGRESS" {
			return resource.RetryableError(fmt.Errorf("task %s is not completed yet", taskID))
		} else if getTaskStatus.TaskDetails.Status == "FAILURE" {
			for _, subtask := range getTaskStatus.ItemList {
				if subtask.Status == "FAILURE" {
					log.Printf("[ERROR]: task %s has failed at step %s, rolling back", taskID, subtask.Name)
				}
			}
			log.Println("[DEBUG]: force decommisioning edge with id", d.Id())
			appResponseData, err := prosimoClient.ForceDeleteAppDeployment(ctx, d.Id())
			if err != nil {
				return resource.NonRetryableError(err)
			}
			log.Printf("[DEBUG] Waiting for task id %s to complete", appResponseData.ResourceData.ID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskComplete(ctx, d, meta, appResponseData.ResourceData.ID))
			if err != nil {
				return resource.NonRetryableError(err)
			}
			log.Printf("[INFO] task %s is successful", appResponseData.ResourceData.ID)

		}
		return nil
	}
}

func retryUntilTaskCompleteNetworkOnboard(ctx context.Context, d *schema.ResourceData, meta interface{}, taskID string, networkOnboardops *client.NetworkOnboardoptns) resource.RetryFunc {
	prosimoClient := meta.(*client.ProsimoClient)
	return func() *resource.RetryError {
		getTaskStatus, err := prosimoClient.GetTaskStatus(ctx, taskID)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if getTaskStatus.TaskDetails.Status == "IN-PROGRESS" {
			return resource.RetryableError(fmt.Errorf("task %s is not completed yet", taskID))
		} else if getTaskStatus.TaskDetails.Status == "FAILURE" {
			for _, subtask := range getTaskStatus.ItemList {
				if subtask.Status == "FAILURE" {
					log.Printf("[ERROR]: task %s has failed at step %s, rolling back", taskID, subtask.Name)
				}
			}
			// resourceNetworkOnboardingRead(ctx, d, meta)
			log.Println("[DEBUG]: offboarding network", d.Id())
			networkOnboardops.ID = d.Id()
			res, err := prosimoClient.OffboardNetworkDeployment(ctx, networkOnboardops.ID)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			log.Printf("[DEBUG] Waiting for task id %s to complete", res.NetworkDeploymentResops.TaskID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskCompleteNetworkOffboard(ctx, d, meta, res.NetworkDeploymentResops.TaskID, networkOnboardops))
			if err != nil {
				return resource.NonRetryableError(err)
			}
			log.Printf("[INFO] task %s is successful", res.NetworkDeploymentResops.TaskID)

		}
		return nil
	}
}

func retryUntilTaskCompleteNetworkOffboard(ctx context.Context, d *schema.ResourceData, meta interface{}, taskID string, networkOnboardops *client.NetworkOnboardoptns) resource.RetryFunc {
	prosimoClient := meta.(*client.ProsimoClient)
	return func() *resource.RetryError {
		getTaskStatus, err := prosimoClient.GetTaskStatus(ctx, taskID)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if getTaskStatus.TaskDetails.Status == "IN-PROGRESS" {
			return resource.RetryableError(fmt.Errorf("task %s is not completed yet", taskID))
		} else if getTaskStatus.TaskDetails.Status == "FAILURE" {
			for _, subtask := range getTaskStatus.ItemList {
				if subtask.Status == "FAILURE" {
					log.Printf("[ERROR]: task %s has failed at step %s, rolling back", taskID, subtask.Name)
				}
			}
			// resourceNetworkOnboardingRead(ctx, d, meta)
			if d.Get("force_offboard").(bool) {
				log.Println("[DEBUG]: force offboarding network", d.Id())
				networkOnboardops.ID = d.Id()
				res, err := prosimoClient.ForceOffboardNetworkDeployment(ctx, networkOnboardops.ID)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				log.Printf("[DEBUG] Waiting for task id %s to complete", res.NetworkDeploymentResops.TaskID)
				err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
					retryUntilTaskComplete(ctx, d, meta, res.NetworkDeploymentResops.TaskID))
				if err != nil {
					return resource.NonRetryableError(err)
				}
			}

		}
		return nil
	}
}

func retryUntilTaskCompleteSharedService(ctx context.Context, d *schema.ResourceData, meta interface{}, taskID string) resource.RetryFunc {
	prosimoClient := meta.(*client.ProsimoClient)
	return func() *resource.RetryError {
		getTaskStatus, err := prosimoClient.GetTaskStatus(ctx, taskID)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if getTaskStatus.TaskDetails.Status == "IN-PROGRESS" {
			return resource.RetryableError(fmt.Errorf("task %s is not completed yet", taskID))
		} else if getTaskStatus.TaskDetails.Status == "FAILURE" {
			for _, subtask := range getTaskStatus.ItemList {
				if subtask.Status == "FAILURE" {
					log.Printf("[ERROR]: task %s has failed at step %s, rolling back", taskID, subtask.Name)
				}
			}
			// resourceNetworkOnboardingRead(ctx, d, meta)
			log.Println("[DEBUG]: offboarding shared service", d.Id())
			_, err := prosimoClient.DecomSharedService(ctx, d.Id())
			if err != nil {
				return resource.NonRetryableError(err)
			}

		}
		return nil
	}
}
func retryUntilTaskCompleteManagedFirewall(ctx context.Context, d *schema.ResourceData, meta interface{}, taskID string) resource.RetryFunc {
	prosimoClient := meta.(*client.ProsimoClient)
	return func() *resource.RetryError {
		getTaskStatus, err := prosimoClient.GetTaskStatus(ctx, taskID)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if getTaskStatus.TaskDetails.Status == "IN-PROGRESS" {
			return resource.RetryableError(fmt.Errorf("task %s is not completed yet", taskID))
		} else if getTaskStatus.TaskDetails.Status == "FAILURE" {
			for _, subtask := range getTaskStatus.ItemList {
				if subtask.Status == "FAILURE" {
					log.Printf("[ERROR]: task %s has failed at step %s, rolling back", taskID, subtask.Name)
				}
			}
			// resourceNetworkOnboardingRead(ctx, d, meta)
			log.Println("[DEBUG]: offboarding shared service", d.Id())
			_, err := prosimoClient.DecomFirewall(ctx, d.Id())
			if err != nil {
				return resource.NonRetryableError(err)
			}

		}
		return nil
	}
}
