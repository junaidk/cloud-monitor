package cloudman

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
)

type AWSCloud struct {
	credentials AWSCredentials
	session     *session.Session
}

func NewAWSCloud(c AWSCredentials) *AWSCloud {

	creds := credentials.NewStaticCredentials(c.AccessKey, c.AccessSecret, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
	})
	if err != nil {
		log.Println(err)
	}
	out := AWSCloud{
		credentials: c,
		session:     sess,
	}

	return &out

}

type AWSCredentials struct {
	AccessKey    string
	AccessSecret string
}

func (c *AWSCloud) GetInstanceListAllRegions() ([]InstanceListResponse, error) {

	rList := []string{
		"ca-central-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-south-1",
		"eu-west-3",
		"eu-north-1",
		"me-south-1",
		"sa-east-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"us-east-2",
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"af-south-1",
		"ap-east-1",
		"ap-south-1",
	}

	var out []InstanceListResponse

	recvChan := make(chan []InstanceListResponse, len(rList))

	for _, region := range rList {
		go func(reg string) {
			res, _ := c.GetInstanceList(reg)
			recvChan <- res
		}(region)
	}

	for j := 0; j < len(rList); j++ {
		list := <-recvChan
		out = append(out, list...)
	}

	return out, nil
}

func (c *AWSCloud) GetInstanceList(region string) ([]InstanceListResponse, error) {

	sess := c.session.Copy(&aws.Config{
		Region: aws.String(region),
	})
	client := ec2.New(sess)

	input := &ec2.DescribeInstancesInput{}

	result, err := client.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return nil, err
	}

	var out []InstanceListResponse
	for _, res := range result.Reservations {
		for _, ins := range res.Instances {
			name := c.getInstanceName(ins.Tags)
			obj := InstanceListResponse{
				Name:        name,
				Id:          *ins.InstanceId,
				Status:      *ins.State.Name,
				LaunchDate:  ins.LaunchTime.String(),
				Region:      region,
				MachineType: *ins.InstanceType,
			}
			out = append(out, obj)
		}

	}

	return out, nil
}

func (c *AWSCloud) getInstanceName(t []*ec2.Tag) string {
	name := ""
	for _, tag := range t {
		if *tag.Key == "Name" {
			name = *tag.Value
		}
	}
	return name
}

func (c *AWSCloud) GetInstanceDetails() {

}
