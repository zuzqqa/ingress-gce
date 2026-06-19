import (
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/cloud-provider-gcp/providers/gce"
	"k8s.io/ingress-gce/pkg/utils/common"
)

func TestIsLegacyL4ILBService(t *testing.T) {
	t.Parallel()
	svc := &api_v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:        "testsvc",
			Namespace:   "default",
			Annotations: map[string]string{gce.ServiceAnnotationLoadBalancerType: string(gce.LBTypeInternal)},
			Finalizers:  []string{common.LegacyILBFinalizer},
		},
		Spec: api_v1.ServiceSpec{
			Type: api_v1.ServiceTypeLoadBalancer,
			Ports: []api_v1.ServicePort{
				{Name: "testport", Port: int32(80)},
			},
		},
	}
	if !IsLegacyL4ILBService(svc) {
		t.Errorf("Expected True for Legacy service %s, got False", svc.Name)
	}

	// Remove the finalizer and ensure the check returns False.
	svc.ObjectMeta.Finalizers = nil
	if IsLegacyL4ILBService(svc) {
		t.Errorf("Expected False for Legacy service %s, got True", svc.Name)
	}
}

func TestLBBasedOnFinalizer(t *testing.T) {
	type wants struct {
		IsLegacyL4ILBService      bool
		IsSubsettingL4ILBService  bool
		HasLegacyL4ILBFinalizerV1 bool
		HasL4ILBFinalizerV2       bool
		HasL4NetLBFinalizerV2     bool
		HasL4NetLBFinalizerV3     bool
	}

	testCases := []struct {
		finalizer string
		want      wants
	}{
		{
			finalizer: common.LegacyILBFinalizer,
			want: wants{
				IsLegacyL4ILBService:      true,
				HasLegacyL4ILBFinalizerV1: true,
			},
		},
		{
			finalizer: common.ILBFinalizerV2,
			want: wants{
				IsSubsettingL4ILBService: true,
				HasL4ILBFinalizerV2:      true,
			},
		},
		{
			finalizer: common.NetLBFinalizerV2,
			want: wants{
				HasL4NetLBFinalizerV2: true,
			},
		},
		{
			finalizer: common.NetLBFinalizerV3,
			want: wants{
				HasL4NetLBFinalizerV3: true,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.finalizer, func(t *testing.T) {
			svc := &api_v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{tC.finalizer},
				},
			}

			got := wants{
				IsLegacyL4ILBService:      IsLegacyL4ILBService(svc),
				IsSubsettingL4ILBService:  IsSubsettingL4ILBService(svc),
				HasLegacyL4ILBFinalizerV1: HasLegacyL4ILBFinalizerV1(svc),
				HasL4ILBFinalizerV2:       HasL4ILBFinalizerV2(svc),

				HasL4NetLBFinalizerV2: HasL4NetLBFinalizerV2(svc),
				HasL4NetLBFinalizerV3: HasL4NetLBFinalizerV3(svc),
			}

			if diff := cmp.Diff(tC.want, got); diff != "" {
				t.Errorf("got != want, diff(-tc.want +got) = %s", diff)
			}
		})
	}
}
