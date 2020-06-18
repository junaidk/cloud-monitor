package cloudman

import (
	"context"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"log"
	"strconv"
	"strings"
)

type GCPCloud struct {
	credentials GCPCredentials
}

// first login using gcloud cli
// use autogenerated credentials from cli to authenicate as acount level
func NewGCPCloud(c GCPCredentials) *GCPCloud {

	out := GCPCloud{
		credentials: c,
	}

	return &out

}

type GCPCredentials struct {
	CredentialFilePath string
}

func (c *GCPCloud) GetInstanceListAllRegions() ([]InstanceListResponse, error) {

	rList := c.GetProjects()

	var out []InstanceListResponse

	recvChan := make(chan []InstanceListResponse, len(rList))

	for _, project := range rList {
		go func(reg string) {
			res, _ := c.GetInstanceList(reg)
			recvChan <- res
		}(project)
	}

	for j := 0; j < len(rList); j++ {
		list := <-recvChan
		out = append(out, list...)
	}

	return out, nil
}

func (c *GCPCloud) GetInstanceList(project string) ([]InstanceListResponse, error) {

	computeService, err := compute.NewService(context.Background(), option.WithCredentialsFile(c.credentials.CredentialFilePath))
	if err != nil {
		log.Fatal(err)
	}

	var res []InstanceListResponse
	reqC := computeService.Instances.AggregatedList(project)
	if err := reqC.Pages(context.Background(), func(page *compute.InstanceAggregatedList) error {
		for _, instance := range page.Items {
			// TODO: Change code below to process each `project` resource:
			for _, ins := range instance.Instances {
				//fmt.Printf("%#v\n", ins.Name)
				obj := InstanceListResponse{
					Name:        ins.Name,
					Id:          strconv.Itoa(int(ins.Id)),
					Status:      ins.Status,
					LaunchDate:  ins.CreationTimestamp,
					Region:      c.extract(ins.Zone),
					MachineType: c.extract(ins.MachineType),
					Project:     project,
				}
				res = append(res, obj)
			}

		}
		return nil
	}); err != nil {
		//log.Println(err)
	}

	return res, nil
}

// split and extract last part or uri
// e.g https://www.googleapis.com/compute/v1/projects/cloudplex-infrastructure/zones/us-central1-b
func (c *GCPCloud) extract(in string) string {

	out := strings.Split(in, "/")
	if len(out) > 0 {
		return out[len(out)-1]
	} else {
		return in
	}
}

func (c *GCPCloud) GetInstanceDetails() {

}

// try to cache it
func (c *GCPCloud) GetProjects() []string {

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(context.Background(), option.WithCredentialsFile(c.credentials.CredentialFilePath))
	if err != nil {
		log.Fatal(err)
	}

	var projList []string
	req := cloudresourcemanagerService.Projects.List()
	if err := req.Pages(context.Background(), func(page *cloudresourcemanager.ListProjectsResponse) error {
		for _, project := range page.Projects {
			// TODO: Change code below to process each `project` resource:
			projList = append(projList, project.ProjectId)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	return projList
}