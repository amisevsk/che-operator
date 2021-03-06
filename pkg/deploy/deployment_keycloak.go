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
	"context"
	"regexp"
	"strconv"
	"strings"

	orgv1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	"github.com/eclipse/che-operator/pkg/util"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	KeycloakDeploymentName   = "keycloak"
	selectSslRequiredCommand = "OUT=$(psql keycloak -tAc \"SELECT 1 FROM REALM WHERE id = 'master'\"); " +
		"if [ $OUT -eq 1 ]; then psql keycloak -tAc \"SELECT ssl_required FROM REALM WHERE id = 'master'\"; fi"
	updateSslRequiredCommand = "psql keycloak -c \"update REALM set ssl_required='NONE' where id = 'master'\""
)

var (
	trustpass              = util.GeneratePasswd(12)
	keycloakCustomDiffOpts = cmp.Options{
		cmp.Comparer(func(x, y appsv1.Deployment) bool {
			return x.Annotations["che.self-signed-certificate.version"] == y.Annotations["che.self-signed-certificate.version"] &&
				x.Annotations["che.openshift-api-crt.version"] == y.Annotations["che.openshift-api-crt.version"] &&
				x.Annotations["che.keycloak-ssl-required-updated"] == y.Annotations["che.keycloak-ssl-required-updated"]
		}),
	}
	keycloakAdditionalDeploymentMerge = func(specDeployment *appsv1.Deployment, clusterDeployment *appsv1.Deployment) *appsv1.Deployment {
		clusterDeployment.Spec = specDeployment.Spec
		clusterDeployment.Annotations["che.self-signed-certificate.version"] = specDeployment.Annotations["che.self-signed-certificate.version"]
		clusterDeployment.Annotations["che.openshift-api-crt.version"] = specDeployment.Annotations["che.openshift-api-crt.version"]
		clusterDeployment.Annotations["che.keycloak-ssl-required-updated"] = specDeployment.Annotations["che.keycloak-ssl-required-updated"]
		return clusterDeployment
	}
)

func SyncKeycloakDeploymentToCluster(deployContext *DeployContext) DeploymentProvisioningStatus {
	clusterDeployment, err := getClusterDeployment(KeycloakDeploymentName, deployContext.CheCluster.Namespace, deployContext.ClusterAPI.Client)
	if err != nil {
		return DeploymentProvisioningStatus{
			ProvisioningStatus: ProvisioningStatus{Err: err},
		}
	}

	specDeployment, err := getSpecKeycloakDeployment(deployContext, clusterDeployment)
	if err != nil {
		return DeploymentProvisioningStatus{
			ProvisioningStatus: ProvisioningStatus{Err: err},
		}
	}

	return SyncDeploymentToCluster(deployContext, specDeployment, clusterDeployment, keycloakCustomDiffOpts, keycloakAdditionalDeploymentMerge)
}

func getSpecKeycloakDeployment(
	deployContext *DeployContext,
	clusterDeployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	optionalEnv := true
	labels := GetLabels(deployContext.CheCluster, KeycloakDeploymentName)
	cheFlavor := DefaultCheFlavor(deployContext.CheCluster)
	keycloakImage := util.GetValue(deployContext.CheCluster.Spec.Auth.IdentityProviderImage, DefaultKeycloakImage(deployContext.CheCluster))
	pullPolicy := corev1.PullPolicy(util.GetValue(string(deployContext.CheCluster.Spec.Auth.IdentityProviderImagePullPolicy), DefaultPullPolicyFromDockerImage(keycloakImage)))
	jbossDir := "/opt/eap"
	if cheFlavor == "che" {
		// writable dir in the upstream Keycloak image
		jbossDir = "/scripts"
	}

	if clusterDeployment != nil {
		env := clusterDeployment.Spec.Template.Spec.Containers[0].Env
		for _, e := range env {
			// To be compatible with prev deployments when "TRUSTPASS" env was used
			if "TRUSTPASS" == e.Name || "SSO_TRUSTSTORE_PASSWORD" == e.Name {
				trustpass = e.Value
				break
			}
		}
	}

	terminationGracePeriodSeconds := int64(30)
	cheCertSecretVersion := getSecretResourceVersion("self-signed-certificate", deployContext.CheCluster.Namespace, deployContext.ClusterAPI)
	openshiftApiCertSecretVersion := getSecretResourceVersion("openshift-api-crt", deployContext.CheCluster.Namespace, deployContext.ClusterAPI)

	// add various certificates to Java trust store so that Keycloak can connect to OpenShift API
	// certificate that OpenShift router uses (for 4.0 only)
	addRouterCrt := "if [ ! -z \"${CHE_SELF__SIGNED__CERT}\" ]; then echo \"${CHE_SELF__SIGNED__CERT}\" > " + jbossDir + "/che.crt && " +
		"keytool -importcert -alias ROUTERCRT" +
		" -keystore " + jbossDir + "/openshift.jks" +
		" -file " + jbossDir + "/che.crt -storepass " + trustpass + " -noprompt; fi"
	// certificate retrieved from http call to OpenShift API endpoint
	addOpenShiftAPICrt := "if [ ! -z \"${OPENSHIFT_SELF__SIGNED__CERT}\" ]; then echo \"${OPENSHIFT_SELF__SIGNED__CERT}\" > " + jbossDir + "/openshift.crt && " +
		"keytool -importcert -alias OPENSHIFTAPI" +
		" -keystore " + jbossDir + "/openshift.jks" +
		" -file " + jbossDir + "/openshift.crt -storepass " + trustpass + " -noprompt; fi"
	// certificate mounted into container /var/run/secrets
	addMountedCrt := " keytool -importcert -alias MOUNTEDCRT" +
		" -keystore " + jbossDir + "/openshift.jks" +
		" -file /var/run/secrets/kubernetes.io/serviceaccount/ca.crt -storepass " + trustpass + " -noprompt"
	addMountedServiceCrt := "if [ -f /var/run/secrets/kubernetes.io/serviceaccount/service-ca.crt ]; then " +
		"keytool -importcert -alias MOUNTEDSERVICECRT" +
		" -keystore " + jbossDir + "/openshift.jks" +
		" -file /var/run/secrets/kubernetes.io/serviceaccount/service-ca.crt -storepass " + trustpass + " -noprompt; fi"
	importJavaCacerts := "keytool -importkeystore -srckeystore /etc/pki/ca-trust/extracted/java/cacerts" +
		" -destkeystore " + jbossDir + "/openshift.jks" +
		" -srcstorepass changeit -deststorepass " + trustpass

	customPublicCertsDir := "/public-certs"
	customPublicCertsVolumeSource := corev1.VolumeSource{}
	if deployContext.CheCluster.Spec.Server.ServerTrustStoreConfigMapName != "" {
		customPublicCertsVolumeSource = corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: deployContext.CheCluster.Spec.Server.ServerTrustStoreConfigMapName,
				},
			},
		}
	}
	customPublicCertsVolume := corev1.Volume{
		Name:         "che-public-certs",
		VolumeSource: customPublicCertsVolumeSource,
	}
	customPublicCertsVolumeMount := corev1.VolumeMount{
		Name:      "che-public-certs",
		MountPath: customPublicCertsDir,
	}
	addCustomPublicCertsCommand := "if [[ -d \"" + customPublicCertsDir + "\" && -n \"$(find " + customPublicCertsDir + " -type f)\" ]]; then " +
		"for certfile in " + customPublicCertsDir + "/* ; do " +
		"keytool -importcert -alias CERT_$(basename $certfile) -keystore " + jbossDir + "/openshift.jks -file $certfile -storepass " + trustpass + " -noprompt; " +
		"done; fi"

	addCertToTrustStoreCommand := addRouterCrt + " && " + addOpenShiftAPICrt + " && " + addMountedCrt + " && " + addMountedServiceCrt + " && " + importJavaCacerts + " && " + addCustomPublicCertsCommand

	// upstream Keycloak has a bit different mechanism of adding jks
	changeConfigCommand := "echo Installing certificates into Keycloak && " +
		"echo -e \"embed-server --server-config=standalone.xml --std-out=echo \n" +
		"/subsystem=keycloak-server/spi=truststore/:add \n" +
		"/subsystem=keycloak-server/spi=truststore/provider=file/:add(properties={file => " +
		"\"" + jbossDir + "/openshift.jks\", password => \"" + trustpass + "\", disabled => \"false\" },enabled=true) \n" +
		"stop-embedded-server\" > /scripts/add_openshift_certificate.cli && " +
		"/opt/jboss/keycloak/bin/jboss-cli.sh --file=/scripts/add_openshift_certificate.cli"

	addProxyCliCommand := ""
	applyProxyCliCommand := ""
	proxyEnvVars := []corev1.EnvVar{}

	if deployContext.Proxy.HttpProxy != "" {
		proxyEnvVars = []corev1.EnvVar{
			corev1.EnvVar{
				Name:  "HTTP_PROXY",
				Value: deployContext.Proxy.HttpProxy,
			},
			corev1.EnvVar{
				Name:  "HTTPS_PROXY",
				Value: deployContext.Proxy.HttpsProxy,
			},
			corev1.EnvVar{
				Name:  "NO_PROXY",
				Value: deployContext.Proxy.NoProxy,
			},
		}

		quotedNoProxy := ""
		for _, noProxyHost := range strings.Split(deployContext.Proxy.NoProxy, ",") {
			if len(quotedNoProxy) != 0 {
				quotedNoProxy += ","
			}
			quotedNoProxy += "\"" + strings.ReplaceAll(regexp.QuoteMeta(noProxyHost), "\\", "\\\\\\") + ";NO_PROXY\""
		}

		jbossCli := "/opt/jboss/keycloak/bin/jboss-cli.sh"
		serverConfig := "standalone.xml"
		if cheFlavor == "codeready" {
			jbossCli = "/opt/eap/bin/jboss-cli.sh"
			serverConfig = "standalone-openshift.xml"
		}
		addProxyCliCommand = " && echo Configuring Proxy && " +
			"echo -e 'embed-server --server-config=" + serverConfig + " --std-out=echo \n" +
			"/subsystem=keycloak-server/spi=connectionsHttpClient/provider=default:write-attribute(name=properties.proxy-mappings,value=[" + quotedNoProxy + ",\".*;" + deployContext.Proxy.HttpProxy + "\"]) \n" +
			"stop-embedded-server' > " + jbossDir + "/setup-http-proxy.cli"

		applyProxyCliCommand = " && " + jbossCli + " --file=" + jbossDir + "/setup-http-proxy.cli"
		if cheFlavor == "codeready" {
			applyProxyCliCommand = " && mkdir -p " + jbossDir + "/extensions && echo '#!/bin/bash\n" +
				"" + jbossDir + "/bin/jboss-cli.sh --file=" + jbossDir + "/setup-http-proxy.cli' > " + jbossDir + "/extensions/postconfigure.sh && " +
				"chmod a+x " + jbossDir + "/extensions/postconfigure.sh "
		}
	}

	keycloakEnv := []corev1.EnvVar{
		{
			Name:  "PROXY_ADDRESS_FORWARDING",
			Value: "true",
		},
		{
			Name:  "DB_VENDOR",
			Value: "POSTGRES",
		},
		{
			Name:  "POSTGRES_PORT_5432_TCP_ADDR",
			Value: util.GetValue(deployContext.CheCluster.Spec.Database.ChePostgresHostName, DefaultChePostgresHostName),
		},
		{
			Name:  "POSTGRES_PORT_5432_TCP_PORT",
			Value: util.GetValue(deployContext.CheCluster.Spec.Database.ChePostgresPort, DefaultChePostgresPort),
		},
		{
			Name:  "POSTGRES_PORT",
			Value: util.GetValue(deployContext.CheCluster.Spec.Database.ChePostgresPort, DefaultChePostgresPort),
		},
		{
			Name:  "POSTGRES_ADDR",
			Value: util.GetValue(deployContext.CheCluster.Spec.Database.ChePostgresHostName, DefaultChePostgresHostName),
		},
		{
			Name:  "POSTGRES_DATABASE",
			Value: "keycloak",
		},
		{
			Name:  "POSTGRES_USER",
			Value: "keycloak",
		},
		{
			Name:  "SSO_TRUSTSTORE",
			Value: "openshift.jks",
		},
		{
			Name:  "SSO_TRUSTSTORE_DIR",
			Value: jbossDir,
		},
		{
			Name:  "SSO_TRUSTSTORE_PASSWORD",
			Value: trustpass,
		},
		{
			Name: "CHE_SELF__SIGNED__CERT",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "ca.crt",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: CheTLSSelfSignedCertificateSecretName,
					},
					Optional: &optionalEnv,
				},
			},
		},
		{
			Name: "OPENSHIFT_SELF__SIGNED__CERT",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "ca.crt",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "openshift-api-crt",
					},
					Optional: &optionalEnv,
				},
			},
		},
	}

	identityProviderPostgresSecret := deployContext.CheCluster.Spec.Auth.IdentityProviderPostgresSecret
	if len(identityProviderPostgresSecret) > 0 {
		keycloakEnv = append(keycloakEnv, corev1.EnvVar{
			Name: "POSTGRES_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "password",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: identityProviderPostgresSecret,
					},
				},
			},
		})
	} else {
		keycloakEnv = append(keycloakEnv, corev1.EnvVar{
			Name:  "POSTGRES_PASSWORD",
			Value: deployContext.CheCluster.Spec.Auth.IdentityProviderPostgresPassword,
		})
	}

	identityProviderSecret := deployContext.CheCluster.Spec.Auth.IdentityProviderSecret
	if len(identityProviderSecret) > 0 {
		keycloakEnv = append(keycloakEnv, corev1.EnvVar{
			Name: "KEYCLOAK_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "password",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: identityProviderSecret,
					},
				},
			},
		},
			corev1.EnvVar{
				Name: "KEYCLOAK_USER",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						Key: "user",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: identityProviderSecret,
						},
					},
				},
			})
	} else {
		keycloakEnv = append(keycloakEnv, corev1.EnvVar{
			Name:  "KEYCLOAK_PASSWORD",
			Value: deployContext.CheCluster.Spec.Auth.IdentityProviderPassword,
		},
			corev1.EnvVar{
				Name:  "KEYCLOAK_USER",
				Value: deployContext.CheCluster.Spec.Auth.IdentityProviderAdminUserName,
			})
	}

	if cheFlavor == "codeready" {
		keycloakEnv = []corev1.EnvVar{
			{
				Name:  "PROXY_ADDRESS_FORWARDING",
				Value: "true",
			},
			{
				Name:  "DB_SERVICE_PREFIX_MAPPING",
				Value: "keycloak-postgresql=DB",
			},
			{
				Name:  "KEYCLOAK_POSTGRESQL_SERVICE_HOST",
				Value: util.GetValue(deployContext.CheCluster.Spec.Database.ChePostgresHostName, DefaultChePostgresHostName),
			},
			{
				Name:  "KEYCLOAK_POSTGRESQL_SERVICE_PORT",
				Value: util.GetValue(deployContext.CheCluster.Spec.Database.ChePostgresPort, DefaultChePostgresPort),
			},
			{
				Name:  "DB_DATABASE",
				Value: "keycloak",
			},
			{
				Name:  "DB_USERNAME",
				Value: "keycloak",
			},
			{
				Name:  "DB_VENDOR",
				Value: "POSTGRES",
			},
			{
				Name:  "SSO_TRUSTSTORE",
				Value: "openshift.jks",
			},
			{
				Name:  "SSO_TRUSTSTORE_DIR",
				Value: jbossDir,
			},
			{
				Name:  "SSO_TRUSTSTORE_PASSWORD",
				Value: trustpass,
			},
			{
				Name: "CHE_SELF__SIGNED__CERT",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						Key: "ca.crt",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: CheTLSSelfSignedCertificateSecretName,
						},
						Optional: &optionalEnv,
					},
				},
			},
			{
				Name: "OPENSHIFT_SELF__SIGNED__CERT",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						Key: "ca.crt",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "openshift-api-crt",
						},
						Optional: &optionalEnv,
					},
				},
			},
		}

		identityProviderPostgresSecret := deployContext.CheCluster.Spec.Auth.IdentityProviderPostgresSecret
		if len(identityProviderPostgresSecret) > 0 {
			keycloakEnv = append(keycloakEnv, corev1.EnvVar{
				Name: "DB_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						Key: "password",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: identityProviderPostgresSecret,
						},
					},
				},
			})
		} else {
			keycloakEnv = append(keycloakEnv, corev1.EnvVar{
				Name:  "DB_PASSWORD",
				Value: deployContext.CheCluster.Spec.Auth.IdentityProviderPostgresPassword,
			})
		}

		identityProviderSecret := deployContext.CheCluster.Spec.Auth.IdentityProviderSecret
		if len(identityProviderSecret) > 0 {
			keycloakEnv = append(keycloakEnv, corev1.EnvVar{
				Name: "SSO_ADMIN_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						Key: "password",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: identityProviderSecret,
						},
					},
				},
			},
				corev1.EnvVar{
					Name: "SSO_ADMIN_USERNAME",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							Key: "user",
							LocalObjectReference: corev1.LocalObjectReference{
								Name: identityProviderSecret,
							},
						},
					},
				})
		} else {
			keycloakEnv = append(keycloakEnv, corev1.EnvVar{
				Name:  "SSO_ADMIN_PASSWORD",
				Value: deployContext.CheCluster.Spec.Auth.IdentityProviderPassword,
			},
				corev1.EnvVar{
					Name:  "SSO_ADMIN_USERNAME",
					Value: deployContext.CheCluster.Spec.Auth.IdentityProviderAdminUserName,
				})
		}
	}

	for _, envvar := range proxyEnvVars {
		keycloakEnv = append(keycloakEnv, envvar)
	}

	command := addCertToTrustStoreCommand + addProxyCliCommand + applyProxyCliCommand + " && " + changeConfigCommand +
		" && /opt/jboss/docker-entrypoint.sh -b 0.0.0.0 -c standalone.xml"
	command += " -Dkeycloak.profile.feature.token_exchange=enabled -Dkeycloak.profile.feature.admin_fine_grained_authz=enabled"
	if cheFlavor == "codeready" {
		command = addCertToTrustStoreCommand + addProxyCliCommand + applyProxyCliCommand +
			" && echo \"feature.token_exchange=enabled\nfeature.admin_fine_grained_authz=enabled\" > /opt/eap/standalone/configuration/profile.properties  " +
			" && sed -i 's/WILDCARD/ANY/g' /opt/eap/bin/launch/keycloak-spi.sh && /opt/eap/bin/openshift-launch.sh -b 0.0.0.0"
	}

	sslRequiredUpdatedForMasterRealm := isSslRequiredUpdatedForMasterRealm(deployContext)
	if sslRequiredUpdatedForMasterRealm {
		// update command to restart pod
		command = "echo \"ssl_required WAS UPDATED for master realm.\" && " + command
	}

	args := []string{"-c", command}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      KeycloakDeploymentName,
			Namespace: deployContext.CheCluster.Namespace,
			Labels:    labels,
			Annotations: map[string]string{
				"che.self-signed-certificate.version": cheCertSecretVersion,
				"che.openshift-api-crt.version":       openshiftApiCertSecretVersion,
				"che.keycloak-ssl-required-updated":   strconv.FormatBool(sslRequiredUpdatedForMasterRealm),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						customPublicCertsVolume,
					},
					Containers: []corev1.Container{
						{
							Name:            KeycloakDeploymentName,
							Image:           keycloakImage,
							ImagePullPolicy: pullPolicy,
							Command: []string{
								"/bin/sh",
							},
							Args: args,
							Ports: []corev1.ContainerPort{
								{
									Name:          KeycloakDeploymentName,
									ContainerPort: 8080,
									Protocol:      "TCP",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("2Gi"),
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "auth/js/keycloak.js",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: int32(8080),
										},
										Scheme: corev1.URISchemeHTTP,
									},
								},
								InitialDelaySeconds: 25,
								FailureThreshold:    10,
								TimeoutSeconds:      5,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
							},
							Env: keycloakEnv,
							VolumeMounts: []corev1.VolumeMount{
								customPublicCertsVolumeMount,
							},
						},
					},
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					RestartPolicy:                 "Always",
				},
			},
		},
	}

	if !util.IsTestMode() {
		err := controllerutil.SetControllerReference(deployContext.CheCluster, deployment, deployContext.ClusterAPI.Scheme)
		if err != nil {
			return nil, err
		}
	}

	return deployment, nil
}

func getSecretResourceVersion(name string, namespace string, clusterAPI ClusterAPI) string {
	secret := &corev1.Secret{}
	err := clusterAPI.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, secret)
	if err != nil {
		if !errors.IsNotFound(err) {
			logrus.Errorf("Failed to get %s secret: %s", name, err)
		}
		return ""
	}
	return secret.ResourceVersion
}

func isSslRequiredUpdatedForMasterRealm(deployContext *DeployContext) bool {
	if deployContext.CheCluster.Spec.Database.ExternalDb {
		return false
	}

	if util.IsTestMode() {
		return false
	}

	clusterDeployment, _ := getClusterDeployment(KeycloakDeploymentName, deployContext.CheCluster.Namespace, deployContext.ClusterAPI.Client)
	if clusterDeployment == nil {
		return false
	}

	value, err := strconv.ParseBool(clusterDeployment.ObjectMeta.Annotations["che.keycloak-ssl-required-updated"])
	if err == nil && value {
		return true
	}

	dbValue, _ := getSslRequiredForMasterRealm(deployContext.CheCluster)
	return dbValue == "NONE"
}

func getSslRequiredForMasterRealm(checluster *orgv1.CheCluster) (string, error) {
	podName, err := util.K8sclient.GetDeploymentPod(PostgresDeploymentName, checluster.Namespace)
	if err != nil {
		return "", err
	}

	stdout, err := util.K8sclient.ExecIntoPod(podName, selectSslRequiredCommand, "", checluster.Namespace)
	return strings.TrimSpace(stdout), err
}

func updateSslRequiredForMasterRealm(checluster *orgv1.CheCluster) error {
	podName, err := util.K8sclient.GetDeploymentPod(PostgresDeploymentName, checluster.Namespace)
	if err != nil {
		return err
	}

	_, err = util.K8sclient.ExecIntoPod(podName, updateSslRequiredCommand, "Update ssl_required to NONE", checluster.Namespace)
	return err
}

func ProvisionKeycloakResources(deployContext *DeployContext) error {
	if !deployContext.CheCluster.Spec.Database.ExternalDb {
		value, err := getSslRequiredForMasterRealm(deployContext.CheCluster)
		if err != nil {
			return err
		}

		if value != "NONE" {
			err := updateSslRequiredForMasterRealm(deployContext.CheCluster)
			if err != nil {
				return err
			}
		}
	}

	keycloakProvisionCommand := GetKeycloakProvisionCommand(deployContext.CheCluster)
	podToExec, err := util.K8sclient.GetDeploymentPod(KeycloakDeploymentName, deployContext.CheCluster.Namespace)
	if err != nil {
		logrus.Errorf("Failed to retrieve pod name. Further exec will fail")
	}

	_, err = util.K8sclient.ExecIntoPod(podToExec, keycloakProvisionCommand, "create realm, client and user", deployContext.CheCluster.Namespace)
	return err
}
