import (
	"k8s.io/ingress-gce/pkg/utils/common"
	"k8s.io/ingress-gce/pkg/utils/slice"
)

// IsLegacyL4ILBService returns true if the given LoadBalancer service is managed by service controller.
func IsLegacyL4ILBService(svc *api_v1.Service) bool {
	if svc.Spec.LoadBalancerClass != nil {
		return l4annotations.HasLoadBalancerClass(svc, l4annotations.LegacyRegionalInternalLoadBalancerClass)
	}
	return HasLegacyL4ILBFinalizerV1(svc)
}

// HasLegacyL4ILBFinalizerV1 returns true if the given Service has ILBFinalizerV1
func HasLegacyL4ILBFinalizerV1(svc *api_v1.Service) bool {
	return slice.ContainsString(svc.ObjectMeta.Finalizers, common.LegacyILBFinalizer, nil)
}

// IsSubsettingL4ILBService returns true if the given LoadBalancer service is managed by NEG and L4 controller.
func IsSubsettingL4ILBService(svc *api_v1.Service) bool {
	if svc.Spec.LoadBalancerClass != nil {
		return l4annotations.HasLoadBalancerClass(svc, l4annotations.RegionalInternalLoadBalancerClass)
	}
	return HasL4ILBFinalizerV2(svc)
}

// HasL4ILBFinalizerV2 returns true if the given Service has ILBFinalizerV2
func HasL4ILBFinalizerV2(svc *api_v1.Service) bool {
	return slice.ContainsString(svc.ObjectMeta.Finalizers, common.ILBFinalizerV2, nil)
}

// HasL4NetLBFinalizerV2 returns true if the given Service has NetLBFinalizerV2
func HasL4NetLBFinalizerV2(svc *api_v1.Service) bool {
	return slice.ContainsString(svc.ObjectMeta.Finalizers, common.NetLBFinalizerV2, nil)
}

// HasL4NetLBFinalizerV3 returns true if the given Service has NetLBFinalizerV3
func HasL4NetLBFinalizerV3(svc *api_v1.Service) bool {
	return slice.ContainsString(svc.ObjectMeta.Finalizers, common.NetLBFinalizerV3, nil)
}