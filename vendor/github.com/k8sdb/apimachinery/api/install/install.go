package install

import (
	aci "github.com/k8sdb/apimachinery/api"
	"k8s.io/kubernetes/pkg/apimachinery/announced"
	"k8s.io/kubernetes/pkg/util/sets"
)

func init() {
	if err := announced.NewGroupMetaFactory(
		&announced.GroupMetaFactoryArgs{
			GroupName:                  aci.GroupName,
			VersionPreferenceOrder:     []string{aci.V1alpha1SchemeGroupVersion.Version},
			ImportPrefix:               "github.com/k8sdb/apimachinery/api",
			RootScopedKinds:            sets.NewString("PodSecurityPolicy", "ThirdPartyResource"),
			AddInternalObjectsToScheme: aci.AddToScheme,
		},
		announced.VersionToSchemeFunc{
			aci.V1alpha1SchemeGroupVersion.Version: aci.V1alpha1AddToScheme,
		},
	).Announce().RegisterAndEnable(); err != nil {
		panic(err)
	}
}
