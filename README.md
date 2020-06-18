### Intro 
A simple tool to list all the vm instances run in a cloud account.
Supported clouds are AWS, Azure, GCP and DigitalOcean.


### Build

`go build -o bin/cloud-mon main.go`


### Config file

See `scratch/sample-conf.json` for sample config file.
Fill the necessary details in the config file.
Provide config file path with -c parameter.

`bin/cloud-mon -c /path/to/file.json`

If you want to exclude any cloud, remove it from config file.