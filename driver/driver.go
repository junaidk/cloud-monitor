package driver

import (
	"cloud-monitor/cloudman"
	config "cloud-monitor/config"
)

type Driver struct {
	CloudList []cloudman.Cloud
}

func NewDriver() (*Driver, error) {

	d := &Driver{}

	conf, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	// AWS
	if len(conf.AWS.AccessKey) > 0 && len(conf.AWS.AccessSecret) > 0 {
		awsCreds := cloudman.AWSCredentials{
			AccessKey:    conf.AWS.AccessKey,
			AccessSecret: conf.AWS.AccessSecret,
		}
		aws := cloudman.NewAWSCloud(awsCreds)
		d.CloudList = append(d.CloudList, aws)
	}

	// GCP
	if len(conf.GCP.CredentialFilePath) > 0 {
		gpcCreds := cloudman.GCPCredentials{
			CredentialFilePath: conf.GCP.CredentialFilePath,
		}
		gcp := cloudman.NewGCPCloud(gpcCreds)
		d.CloudList = append(d.CloudList, gcp)
	}

	//Azure
	if len(conf.Azure.ID) > 0 &&
		len(conf.Azure.Key) > 0 &&
		len(conf.Azure.Subscription) > 0 &&
		len(conf.Azure.Tenant) > 0 {

		azureCreds := cloudman.AzureCredentials{
			ID:           conf.Azure.ID,
			Key:          conf.Azure.Key,
			Tenant:       conf.Azure.Tenant,
			Subscription: conf.Azure.Subscription,
		}
		do := cloudman.NewAzureCloud(azureCreds)
		d.CloudList = append(d.CloudList, do)
	}

	//DO
	if len(conf.DO.AccessToken) > 0 {
		doCreds := cloudman.DOCredentials{
			AccessToken: conf.DO.AccessToken,
		}
		do := cloudman.NewDOCloud(doCreds)
		d.CloudList = append(d.CloudList, do)
	}

	// IBM
	if len(conf.IBM.APIKey) > 0 {
		ibmCreds := cloudman.IBMCredentials{
			APIKey: conf.IBM.APIKey,
		}
		do := cloudman.NewIBMCloud(ibmCreds)
		d.CloudList = append(d.CloudList, do)
	}

	return d, nil
}
