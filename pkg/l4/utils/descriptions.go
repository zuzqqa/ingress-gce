const (
	L4ILBServiceDescKey     = "networking.gke.io/service-name"
	L4LBSharedResourcesDesc = "This resource is shared by all L4 %s Services using ExternalTrafficPolicy: Cluster."
)

// L4LBType indicates if L4 LoadBalancer is Internal or External
type L4LBType int

const (
	ILB L4LBType = iota
	XLB
)

func (lbType L4LBType) ToString() string {
	if lbType == ILB {
		return "ILB"
	}
	return "XLB"
}

// L4LBResourceDescription stores the description fields for L4 ILB or NetLB resources.
// This is useful to identify which resources correspond to which L4 LB service.
type L4LBResourceDescription struct {
	// ServiceName indicates the name of the service the resource is for.
	ServiceName string `json:"networking.gke.io/service-name"`
	// APIVersion stores the version og the compute API used to create this resource.
	APIVersion          meta.Version `json:"networking.gke.io/api-version,omitempty"`
	ServiceIP           string       `json:"networking.gke.io/service-ip,omitempty"`
	ResourceDescription string       `json:"networking.gke.io/resource-description,omitempty"`
}

// Marshal returns the description as a JSON-encoded string.
func (d *L4LBResourceDescription) Marshal() (string, error) {
	out, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(out), err
}

// Unmarshal converts the JSON-encoded description string into the struct.
func (d *L4LBResourceDescription) Unmarshal(desc string) error {
	return json.Unmarshal([]byte(desc), d)
}

func MakeL4LBFirewallDescription(svcName, ip string, version meta.Version, shared bool) (string, error) {
	if shared {
		return (&L4LBResourceDescription{APIVersion: version, ResourceDescription: fmt.Sprintf(L4LBSharedResourcesDesc, "")}).Marshal()
	}
	return (&L4LBResourceDescription{ServiceName: svcName, ServiceIP: ip, APIVersion: version}).Marshal()
}

func MakeL4LBServiceDescription(svcName, ip string, version meta.Version, shared bool, lbType L4LBType) (string, error) {
	if shared {
		return (&L4LBResourceDescription{APIVersion: version, ResourceDescription: fmt.Sprintf(L4LBSharedResourcesDesc, lbType.ToString())}).Marshal()
	}
	return (&L4LBResourceDescription{ServiceName: svcName, ServiceIP: ip, APIVersion: version}).Marshal()
}

func MakeL4IPv6ForwardingRuleDescription(service *api_v1.Service) (string, error) {
	return (&L4LBResourceDescription{ServiceName: ServiceKeyFunc(service.Namespace, service.Name)}).Marshal()
}