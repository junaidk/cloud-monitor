package cloudman

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"log"
)

type AzureCloud struct {
	credentials AzureCredentials
	authorizer  *autorest.BearerAuthorizer
}

func NewAzureCloud(c AzureCredentials) *AzureCloud {

	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, c.Tenant)
	if err != nil {
		log.Println(err)
	}

	spt, err := adal.NewServicePrincipalToken(*oauthConfig, c.ID, c.Key, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Println(err)
	}

	authorizer := autorest.NewBearerAuthorizer(spt)

	out := AzureCloud{
		credentials: c,
		authorizer:  authorizer,
	}

	return &out

}

type AzureCredentials struct {
	ID           string
	Key          string
	Tenant       string
	Subscription string
}

func (c *AzureCloud) GetInstanceListAllRegions() ([]InstanceListResponse, error) {

	out, err := c.GetInstanceList("")

	return out, err
}

func (c *AzureCloud) GetInstanceList(project string) ([]InstanceListResponse, error) {

	vmClient := compute.NewVirtualMachinesClient(c.credentials.Subscription)
	vmClient.Authorizer = c.authorizer

	res, err := vmClient.List(context.Background(), "application-0jkqy6")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	vmList := res.Values()

	for err := res.NextWithContext(context.Background()); err != nil; {
		log.Println(err)
		vmList = append(res.Values())
	}

	var out []InstanceListResponse
	for _, vm := range vmList {
		//log.Println(*vm.Name)

		obj := InstanceListResponse{
			Name:        *vm.Name,
			Id:          *vm.ID,
			Status:      *vm.ProvisioningState,
			LaunchDate:  "",
			Region:      *vm.Location,
			MachineType: string(vm.HardwareProfile.VMSize),
		}
		out = append(out, obj)
	}

	return out, nil
}

func (c *AzureCloud) GetInstanceDetails() {

}
