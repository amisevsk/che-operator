//
// Copyright (c) 2012-2019 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//
package deploy

import (
	"encoding/json"
	"fmt"

	orgv1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	"github.com/eclipse/che-operator/pkg/util"
	"github.com/sirupsen/logrus"
)

type DevFileRegistryConfigMap struct {
	CheDevfileImagesRegistryURL          string `json:"CHE_DEVFILE_IMAGES_REGISTRY_URL"`
	CheDevfileImagesRegistryOrganization string `json:"CHE_DEVFILE_IMAGES_REGISTRY_ORGANIZATION"`
	CheDevfileRegistryURL                string `json:"CHE_DEVFILE_REGISTRY_URL"`
}

const (
	DevfileRegistry = "devfile-registry"
)

/**
 * Create devfile registry resources unless an external registry is used.
 */
func SyncDevfileRegistryToCluster(deployContext *DeployContext) (bool, error) {
	devfileRegistryURL := deployContext.CheCluster.Spec.Server.DevfileRegistryUrl
	if !deployContext.CheCluster.Spec.Server.ExternalDevfileRegistry {
		var host string
		if !util.IsOpenShift {
			ingress, err := SyncIngressToCluster(deployContext, DevfileRegistry, "", DevfileRegistry, 8080)
			if !util.IsTestMode() {
				if ingress == nil {
					logrus.Infof("Waiting on ingress '%s' to be ready", DevfileRegistry)
					if err != nil {
						logrus.Error(err)
					}
					return false, err
				}
			}

			ingressStrategy := util.GetValue(deployContext.CheCluster.Spec.K8s.IngressStrategy, DefaultIngressStrategy)
			if ingressStrategy == "multi-host" {
				host = DevfileRegistry + "-" + deployContext.CheCluster.Namespace + "." + deployContext.CheCluster.Spec.K8s.IngressDomain
			} else {
				host = deployContext.CheCluster.Spec.K8s.IngressDomain + "/" + DevfileRegistry
			}
		} else {
			route, err := SyncRouteToCluster(deployContext, DevfileRegistry, "", DevfileRegistry, 8080)
			if !util.IsTestMode() {
				if route == nil {
					logrus.Infof("Waiting on route '%s' to be ready", DevfileRegistry)
					if err != nil {
						logrus.Error(err)
					}

					return false, err
				}
			}

			if !util.IsTestMode() {
				host = route.Spec.Host
			}
		}

		if devfileRegistryURL == "" {
			if deployContext.CheCluster.Spec.Server.TlsSupport {
				devfileRegistryURL = "https://" + host
			} else {
				devfileRegistryURL = "http://" + host
			}
		}

		configMapData := getDevfileRegistryConfigMapData(deployContext.CheCluster, devfileRegistryURL)
		configMapSpec, err := GetSpecConfigMap(deployContext, DevfileRegistry, configMapData)
		if err != nil {
			return false, err
		}

		configMap, err := SyncConfigMapToCluster(deployContext, configMapSpec)
		if configMap == nil {
			return false, err
		}

		// Create a new registry service
		registryLabels := GetLabels(deployContext.CheCluster, DevfileRegistry)
		serviceStatus := SyncServiceToCluster(deployContext, DevfileRegistry, []string{"http"}, []int32{8080}, registryLabels)
		if !util.IsTestMode() {
			if !serviceStatus.Continue {
				logrus.Info("Waiting on service '" + DevfileRegistry + "' to be ready")
				if serviceStatus.Err != nil {
					logrus.Error(serviceStatus.Err)
				}

				return false, serviceStatus.Err
			}
		}

		// Deploy devfile registry
		deploymentStatus := SyncDevfileRegistryDeploymentToCluster(deployContext)
		if !util.IsTestMode() {
			if !deploymentStatus.Continue {
				logrus.Info("Waiting on deployment '" + DevfileRegistry + "' to be ready")
				if deploymentStatus.Err != nil {
					logrus.Error(deploymentStatus.Err)
				}

				return false, deploymentStatus.Err
			}
		}
	}

	if devfileRegistryURL != deployContext.CheCluster.Status.DevfileRegistryURL {
		deployContext.CheCluster.Status.DevfileRegistryURL = devfileRegistryURL
		if err := UpdateCheCRStatus(deployContext, "status: Devfile Registry URL", devfileRegistryURL); err != nil {
			return false, err
		}
	}

	return true, nil
}

func getDevfileRegistryConfigMapData(cr *orgv1.CheCluster, endpoint string) map[string]string {
	devfileRegistryEnv := make(map[string]string)
	data := &DevFileRegistryConfigMap{
		CheDevfileImagesRegistryURL:          cr.Spec.Server.AirGapContainerRegistryHostname,
		CheDevfileImagesRegistryOrganization: cr.Spec.Server.AirGapContainerRegistryOrganization,
		CheDevfileRegistryURL:                endpoint,
	}

	out, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(out, &devfileRegistryEnv)
	if err != nil {
		fmt.Println(err)
	}
	return devfileRegistryEnv
}
