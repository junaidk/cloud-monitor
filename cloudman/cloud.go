package cloudman

type Cloud interface {
	//SetClientCredentials()
	GetInstanceListAllRegions() ([]InstanceListResponse, error)
	GetInstanceList(region string) ([]InstanceListResponse, error)
	GetInstanceDetails()
}

type InstanceListResponse struct {
	Name        string
	Id          string
	Status      string
	LaunchDate  string
	Region      string
	MachineType string
	Project     string
}
