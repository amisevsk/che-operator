// +build !ignore_autogenerated

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/eclipse/che-operator/pkg/apis/org/v1.CheCluster":            schema_pkg_apis_org_v1_CheCluster(ref),
		"github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpec":        schema_pkg_apis_org_v1_CheClusterSpec(ref),
		"github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecAuth":    schema_pkg_apis_org_v1_CheClusterSpecAuth(ref),
		"github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecDB":      schema_pkg_apis_org_v1_CheClusterSpecDB(ref),
		"github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecK8SOnly": schema_pkg_apis_org_v1_CheClusterSpecK8SOnly(ref),
		"github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecServer":  schema_pkg_apis_org_v1_CheClusterSpecServer(ref),
		"github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecStorage": schema_pkg_apis_org_v1_CheClusterSpecStorage(ref),
	}
}

func schema_pkg_apis_org_v1_CheCluster(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "The `CheCluster` custom resource allows defining and managing a Che server installation",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Description: "Desired configuration of the Che installation. Based on these settings, the operator automatically creates and maintains several config maps that will contain the appropriate environment variables the various components of the Che installation. These generated config maps should NOT be updated manually.",
							Ref:         ref("github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "CheClusterStatus defines the observed state of Che installation",
							Ref:         ref("github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpec", "github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_org_v1_CheClusterSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Desired configuration of the Che installation. Based on these settings, the operator automatically creates and maintains several config maps that will contain the appropriate environment variables the various components of the Che installation. These generated config maps should NOT be updated manually.",
				Properties: map[string]spec.Schema{
					"server": {
						SchemaProps: spec.SchemaProps{
							Description: "General configuration settings related to the Che server and the plugin and devfile registries",
							Ref:         ref("github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecServer"),
						},
					},
					"database": {
						SchemaProps: spec.SchemaProps{
							Description: "Configuration settings related to the database used by the Che installation.",
							Ref:         ref("github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecDB"),
						},
					},
					"auth": {
						SchemaProps: spec.SchemaProps{
							Description: "Configuration settings related to the Authentication used by the Che installation.",
							Ref:         ref("github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecAuth"),
						},
					},
					"storage": {
						SchemaProps: spec.SchemaProps{
							Description: "Configuration settings related to the persistent storage used by the Che installation.",
							Ref:         ref("github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecStorage"),
						},
					},
					"metrics": {
						SchemaProps: spec.SchemaProps{
							Description: "Configuration settings related to the metrics collection used by the Che installation.",
							Ref:         ref("github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecMetrics"),
						},
					},
					"k8s": {
						SchemaProps: spec.SchemaProps{
							Description: "Configuration settings specific to Che installations made on upstream Kubernetes.",
							Ref:         ref("github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecK8SOnly"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecAuth", "github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecDB", "github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecK8SOnly", "github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecMetrics", "github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecServer", "github.com/eclipse/che-operator/pkg/apis/org/v1.CheClusterSpecStorage"},
	}
}

func schema_pkg_apis_org_v1_CheClusterSpecAuth(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Configuration settings related to the Authentication used by the Che installation.",
				Properties: map[string]spec.Schema{
					"externalIdentityProvider": {
						SchemaProps: spec.SchemaProps{
							Description: "Instructs the operator on whether or not to deploy a dedicated Identity Provider (Keycloak or RH SSO instance). By default a dedicated Identity Provider server is deployed as part of the Che installation. But if `externalIdentityProvider` is `true`, then no dedicated identity provider will be deployed by the operator and you might need to provide details about the external identity provider you want to use. See also all the other fields starting with: `identityProvider`.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"identityProviderURL": {
						SchemaProps: spec.SchemaProps{
							Description: "Public URL of the Identity Provider server (Keycloak / RH SSO server). You should set it ONLY if you use an external Identity Provider (see the `externalIdentityProvider` field). By default this will be automatically calculated and set by the operator.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"identityProviderAdminUserName": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the name of the Identity Provider admin user. Defaults to `admin`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"identityProviderPassword": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the password of Keycloak admin user. This is useful to override it ONLY if you use an external Identity Provider (see the `externalIdentityProvider` field). If omitted or left blank, it will be set to an auto-generated password.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"identityProviderSecret": {
						SchemaProps: spec.SchemaProps{
							Description: "The secret that contains `user` and `password` for Identity Provider. If the secret is defined then `identityProviderAdminUserName` and `identityProviderPassword` are ignored. If the value is omitted or left blank then there are two scenarios: 1. `identityProviderAdminUserName` and `identityProviderPassword` are defined, then they will be used. 2. `identityProviderAdminUserName` or `identityProviderPassword` are not defined, then a new secret with the name `che-identity-secret` will be created with default value `admin` for `user` and with an auto-generated value for `password`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"identityProviderRealm": {
						SchemaProps: spec.SchemaProps{
							Description: "Name of a Identity provider (Keycloak / RH SSO) realm that should be used for Che. This is useful to override it ONLY if you use an external Identity Provider (see the `externalIdentityProvider` field). If omitted or left blank, it will be set to the value of the `flavour` field.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"identityProviderClientId": {
						SchemaProps: spec.SchemaProps{
							Description: "Name of a Identity provider (Keycloak / RH SSO) `client-id` that should be used for Che. This is useful to override it ONLY if you use an external Identity Provider (see the `externalIdentityProvider` field). If omitted or left blank, it will be set to the value of the `flavour` field suffixed with `-public`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"identityProviderPostgresPassword": {
						SchemaProps: spec.SchemaProps{
							Description: "Password for The Identity Provider (Keycloak / RH SSO) to connect to the database. This is useful to override it ONLY if you use an external Identity Provider (see the `externalIdentityProvider` field). If omitted or left blank, it will be set to an auto-generated password.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"identityProviderPostgresSecret": {
						SchemaProps: spec.SchemaProps{
							Description: "The secret that contains `password` for The Identity Provider (Keycloak / RH SSO) to connect to the database. If the secret is defined then `identityProviderPostgresPassword` will be ignored. If the value is omitted or left blank then there are two scenarios: 1. `identityProviderPostgresPassword` is defined, then it will be used to connect to the database. 2. `identityProviderPostgresPassword` is not defined, then a new secret with the name `che-identity-postgres-secret` will be created with an auto-generated value for `password`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"updateAdminPassword": {
						SchemaProps: spec.SchemaProps{
							Description: "Forces the default `admin` Che user to update password on first login. Defaults to `false`.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"openShiftoAuth": {
						SchemaProps: spec.SchemaProps{
							Description: "Enables the integration of the identity provider (Keycloak / RHSSO) with OpenShift OAuth. Enabled by default on OpenShift. This will allow users to directly login with their Openshift user through the Openshift login, and have their workspaces created under personal OpenShift namespaces. WARNING: the `kubeadmin` user is NOT supported, and logging through it will NOT allow accessing the Che Dashboard.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"oAuthClientName": {
						SchemaProps: spec.SchemaProps{
							Description: "Name of the OpenShift `OAuthClient` resource used to setup identity federation on the OpenShift side. Auto-generated if left blank. See also the `OpenShiftoAuth` field.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"oAuthSecret": {
						SchemaProps: spec.SchemaProps{
							Description: "Name of the secret set in the OpenShift `OAuthClient` resource used to setup identity federation on the OpenShift side. Auto-generated if left blank. See also the `OAuthClientName` field.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"identityProviderImage": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the container image used in the Identity Provider (Keycloak / RH SSO) deployment. This includes the image tag. Omit it or leave it empty to use the defaut container image provided by the operator.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"identityProviderImagePullPolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the image pull policy used in the Identity Provider (Keycloak / RH SSO) deployment. Default value is `Always` for `nightly` or `latest` images, and `IfNotPresent` in other cases.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_org_v1_CheClusterSpecDB(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Configuration settings related to the database used by the Che installation.",
				Properties: map[string]spec.Schema{
					"externalDb": {
						SchemaProps: spec.SchemaProps{
							Description: "Instructs the operator on whether or not to deploy a dedicated database. By default a dedicated Postgres database is deployed as part of the Che installation. But if `externalDb` is `true`, then no dedicated database will be deployed by the operator and you might need to provide connection details to the external DB you want to use. See also all the fields starting with: `chePostgres`.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"chePostgresHostName": {
						SchemaProps: spec.SchemaProps{
							Description: "Postgres Database hostname that the Che server uses to connect to. Defaults to postgres. This value should be overridden ONLY when using an external database (see field `externalDb`). In the default case it will be automatically set by the operator.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"chePostgresPort": {
						SchemaProps: spec.SchemaProps{
							Description: "Postgres Database port that the Che server uses to connect to. Defaults to 5432. This value should be overridden ONLY when using an external database (see field `externalDb`). In the default case it will be automatically set by the operator.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"chePostgresUser": {
						SchemaProps: spec.SchemaProps{
							Description: "Postgres user that the Che server should use to connect to the DB. Defaults to `pgche`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"chePostgresPassword": {
						SchemaProps: spec.SchemaProps{
							Description: "Postgres password that the Che server should use to connect to the DB. If omitted or left blank, it will be set to an auto-generated value.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"chePostgresDb": {
						SchemaProps: spec.SchemaProps{
							Description: "Postgres database name that the Che server uses to connect to the DB. Defaults to `dbche`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"chePostgresSecret": {
						SchemaProps: spec.SchemaProps{
							Description: "The secret that contains Postgres `user` and `password` that the Che server should use to connect to the DB. If the secret is defined then `chePostgresUser` and `chePostgresPassword` are ignored. If the value is omitted or left blank then there are two scenarios: 1. `chePostgresUser` and `chePostgresPassword` are defined, then they will be used to connect to the DB. 2. `chePostgresUser` or `chePostgresPassword` are not defined, then a new secret with the name `che-postgres-secret` will be created with default value of `pgche` for `user` and with an auto-generated value for `password`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"postgresImage": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the container image used in the Postgres database deployment. This includes the image tag. Omit it or leave it empty to use the defaut container image provided by the operator.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"postgresImagePullPolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the image pull policy used in the Postgres database deployment. Default value is `Always` for `nightly` or `latest` images, and `IfNotPresent` in other cases.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_org_v1_CheClusterSpecK8SOnly(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Configuration settings specific to Che installations made on upstream Kubernetes.",
				Properties: map[string]spec.Schema{
					"ingressDomain": {
						SchemaProps: spec.SchemaProps{
							Description: "Global ingress domain for a K8S cluster. This MUST be explicitly specified: there are no defaults.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"ingressStrategy": {
						SchemaProps: spec.SchemaProps{
							Description: "Strategy for ingress creation. This can be `multi-host` (host is explicitly provided in ingress), `single-host` (host is provided, path-based rules) and `default-host.*`(no host is provided, path-based rules). Defaults to `\"multi-host`",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"ingressClass": {
						SchemaProps: spec.SchemaProps{
							Description: "Ingress class that will define the which controler will manage ingresses. Defaults to `nginx`. NB: This drives the `is kubernetes.io/ingress.class` annotation on Che-related ingresses.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"tlsSecretName": {
						SchemaProps: spec.SchemaProps{
							Description: "Name of a secret that will be used to setup ingress TLS termination if TLS is enabled. See also the `tlsSupport` field.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"securityContextFsGroup": {
						SchemaProps: spec.SchemaProps{
							Description: "FSGroup the Che pod and Workspace pods containers should run in. Defaults to `1724`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"securityContextRunAsUser": {
						SchemaProps: spec.SchemaProps{
							Description: "ID of the user the Che pod and Workspace pods containers should run as. Default to `1724`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_org_v1_CheClusterSpecServer(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "General configuration settings related to the Che server and the plugin and devfile registries.",
				Properties: map[string]spec.Schema{
					"airGapContainerRegistryHostname": {
						SchemaProps: spec.SchemaProps{
							Description: "Optional hostname (or url) to an alternate container registry to pull images from. This value overrides the container registry hostname defined in all the default container images involved in a Che deployment. This is particularly useful to install Che in an air-gapped environment.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"airGapContainerRegistryOrganization": {
						SchemaProps: spec.SchemaProps{
							Description: "Optional repository name of an alternate container registry to pull images from. This value overrides the container registry organization defined in all the default container images involved in a Che deployment. This is particularly useful to install Che in an air-gapped environment.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"cheImage": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the container image used in Che deployment. This does NOT include the container image tag. Omit it or leave it empty to use the defaut container image provided by the operator.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"cheImageTag": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the tag of the container image used in Che deployment. Omit it or leave it empty to use the defaut image tag provided by the operator.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"cheImagePullPolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the image pull policy used in Che deployment. Default value is `Always` for `nightly` or `latest` images, and `IfNotPresent` in other cases.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"cheFlavor": {
						SchemaProps: spec.SchemaProps{
							Description: "Flavor of the installation. This is either `che` for upstream Che installations, or `codeready` for CodeReady Workspaces installation. In most cases the default value should not be overridden.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"cheHost": {
						SchemaProps: spec.SchemaProps{
							Description: "Public hostname of the installed Che server. If value is omitted then it will be automatically set by the operator. (see the `cheHostTLSSecret` field).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"cheHostTLSSecret": {
						SchemaProps: spec.SchemaProps{
							Description: "Name of a secret containing certificates to secure ingress/route for the custom hostname of the installed Che server. (see the `cheHost` field).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"cheLogLevel": {
						SchemaProps: spec.SchemaProps{
							Description: "Log level for the Che server: `INFO` or `DEBUG`. Defaults to `INFO`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"cheDebug": {
						SchemaProps: spec.SchemaProps{
							Description: "Enables the debug mode for Che server. Defaults to `false`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"cheWorkspaceClusterRole": {
						SchemaProps: spec.SchemaProps{
							Description: "Custom cluster role bound to the user for the Che workspaces. The default roles are used if this is omitted or left blank.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"workspaceNamespaceDefault": {
						SchemaProps: spec.SchemaProps{
							Description: "Defines Kubernetes default namespace in which user's workspaces are created if user does not override it. It's possible to use <username>, <userid> and <workspaceid> placeholders (e.g.: che-workspace-<username>). In that case, new namespace will be created for each user (or workspace). Is used by OpenShift infra as well to specify Project",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"allowUserDefinedWorkspaceNamespaces": {
						SchemaProps: spec.SchemaProps{
							Description: "Defines if a user is able to specify Kubernetes namespace (or OpenShift project) different from the default. It's NOT RECOMMENDED to configured true without OAuth configured. This property is also used by the OpenShift infra.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"selfSignedCert": {
						SchemaProps: spec.SchemaProps{
							Description: "Deprecated. The value of this flag is ignored. Che operator will automatically detect if router certificate is self-signed. If so it will be propagated to Che server and some other components.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"serverTrustStoreConfigMapName": {
						SchemaProps: spec.SchemaProps{
							Description: "Name of the config-map with public certificates to add to Java trust store of the Che server. This is usually required when adding the OpenShift OAuth provider which has https endpoint signed with self-signed cert. So, Che server must be aware of its CA cert to be able to request it. This is disabled by default.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"gitSelfSignedCert": {
						SchemaProps: spec.SchemaProps{
							Description: "If enabled, then the certificate from `che-git-self-signed-cert` config map will be propagated to the Che components and provide particular configuration for Git.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"tlsSupport": {
						SchemaProps: spec.SchemaProps{
							Description: "Deprecated. Instructs the operator to deploy Che in TLS mode. This is enabled by default. Disabling TLS may cause malfunction of some Che components.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"devfileRegistryUrl": {
						SchemaProps: spec.SchemaProps{
							Description: "Public URL of the Devfile registry, that serves sample, ready-to-use devfiles. You should set it ONLY if you use an external devfile registry (see the `externalDevfileRegistry` field). By default this will be automatically calculated by the operator.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"devfileRegistryImage": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the container image used in the Devfile registry deployment. This includes the image tag. Omit it or leave it empty to use the defaut container image provided by the operator.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"devfileRegistryPullPolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the image pull policy used in the Devfile registry deployment. Default value is `Always` for `nightly` or `latest` images, and `IfNotPresent` in other cases.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"devfileRegistryMemoryLimit": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the memory limit used in the Devfile registry deployment. Defaults to 256Mi.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"devfileRegistryMemoryRequest": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the memory request used in the Devfile registry deployment. Defaults to 16Mi.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"externalDevfileRegistry": {
						SchemaProps: spec.SchemaProps{
							Description: "Instructs the operator on whether or not to deploy a dedicated Devfile registry server. By default a dedicated devfile registry server is started. But if `externalDevfileRegistry` is `true`, then no such dedicated server will be started by the operator and you will have to manually set the `devfileRegistryUrl` field",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"pluginRegistryUrl": {
						SchemaProps: spec.SchemaProps{
							Description: "Public URL of the Plugin registry, that serves sample ready-to-use devfiles. You should set it ONLY if you use an external devfile registry (see the `externalPluginRegistry` field). By default this will be automatically calculated by the operator.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"pluginRegistryImage": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the container image used in the Plugin registry deployment. This includes the image tag. Omit it or leave it empty to use the defaut container image provided by the operator.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"pluginRegistryPullPolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the image pull policy used in the Plugin registry deployment. Default value is `Always` for `nightly` or `latest` images, and `IfNotPresent` in other cases.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"pluginRegistryMemoryLimit": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the memory limit used in the Plugin registry deployment. Defaults to 256Mi.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"pluginRegistryMemoryRequest": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the memory request used in the Plugin registry deployment. Defaults to 16Mi.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"externalPluginRegistry": {
						SchemaProps: spec.SchemaProps{
							Description: "Instructs the operator on whether or not to deploy a dedicated Plugin registry server. By default a dedicated plugin registry server is started. But if `externalPluginRegistry` is `true`, then no such dedicated server will be started by the operator and you will have to manually set the `pluginRegistryUrl` field.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"customCheProperties": {
						SchemaProps: spec.SchemaProps{
							Description: "Map of additional environment variables that will be applied in the generated `che` config map to be used by the Che server, in addition to the values already generated from other fields of the `CheCluster` custom resource (CR). If `customCheProperties` contains a property that would be normally generated in `che` config map from other CR fields, then the value defined in the `customCheProperties` will be used instead.",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"proxyURL": {
						SchemaProps: spec.SchemaProps{
							Description: "URL (protocol+hostname) of the proxy server. This drives the appropriate changes in the `JAVA_OPTS` and `https(s)_proxy` variables in the Che server and workspaces containers. Only use when configuring a proxy is required. Operator respects OpenShift cluster wide proxy configuration and no additional configuration is required, but defining `proxyUrl` in a custom resource leads to overrides the cluster proxy configuration with fields `proxyUrl`, `proxyPort`, `proxyUser` and `proxyPassword` from the custom resource. (see the doc https://docs.openshift.com/container-platform/4.4/networking/enable-cluster-wide-proxy.html) (see also the `proxyPort` and `nonProxyHosts` fields).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"proxyPort": {
						SchemaProps: spec.SchemaProps{
							Description: "Port of the proxy server. Only use when configuring a proxy is required. (see also the `proxyURL` and `nonProxyHosts` fields).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"nonProxyHosts": {
						SchemaProps: spec.SchemaProps{
							Description: "List of hosts that should not use the configured proxy. Use `|`` as delimiter, eg `localhost|my.host.com|123.42.12.32` Only use when configuring a proxy is required. Operator respects OpenShift cluster wide proxy configuration and no additional configuration is required, but defining `nonProxyHosts` in a custom resource leads to merging non proxy hosts lists from the cluster proxy configuration and ones defined in the custom resources. (see the doc https://docs.openshift.com/container-platform/4.4/networking/enable-cluster-wide-proxy.html) (see also the `proxyURL` fields).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"proxyUser": {
						SchemaProps: spec.SchemaProps{
							Description: "User name of the proxy server. Only use when configuring a proxy is required (see also the `proxyURL`, `proxyPassword` and `proxySecret` fields).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"proxyPassword": {
						SchemaProps: spec.SchemaProps{
							Description: "Password of the proxy server Only use when proxy configuration is required (see also the `proxyURL`, `proxyUser` and `proxySecret` fields).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"proxySecret": {
						SchemaProps: spec.SchemaProps{
							Description: "The secret that contains `user` and `password` for a proxy server. If the secret is defined then `proxyUser` and `proxyPassword` are ignored",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"serverMemoryRequest": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the memory request used in the Che server deployment. Defaults to 512Mi.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"serverMemoryLimit": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the memory limit used in the Che server deployment. Defaults to 1Gi.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_org_v1_CheClusterSpecStorage(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Configuration settings related to the persistent storage used by the Che installation.",
				Properties: map[string]spec.Schema{
					"pvcStrategy": {
						SchemaProps: spec.SchemaProps{
							Description: "Persistent volume claim strategy for the Che server. This Can be:`common` (all workspaces PVCs in one volume), `per-workspace` (one PVC per workspace for all declared volumes) and `unique` (one PVC per declared volume). Defaults to `common`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"pvcClaimSize": {
						SchemaProps: spec.SchemaProps{
							Description: "Size of the persistent volume claim for workspaces. Defaults to `1Gi`",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"preCreateSubPaths": {
						SchemaProps: spec.SchemaProps{
							Description: "Instructs the Che server to launch a special pod to pre-create a subpath in the Persistent Volumes. Defaults to `false`, however it might need to enable it according to the configuration of your K8S cluster.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"pvcJobsImage": {
						SchemaProps: spec.SchemaProps{
							Description: "Overrides the container image used to create sub-paths in the Persistent Volumes. This includes the image tag. Omit it or leave it empty to use the defaut container image provided by the operator. See also the `preCreateSubPaths` field.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"postgresPVCStorageClassName": {
						SchemaProps: spec.SchemaProps{
							Description: "Storage class for the Persistent Volume Claim dedicated to the Postgres database. If omitted or left blank, default storage class is used.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"workspacePVCStorageClassName": {
						SchemaProps: spec.SchemaProps{
							Description: "Storage class for the Persistent Volume Claims dedicated to the Che workspaces. If omitted or left blank, default storage class is used.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{},
	}
}
