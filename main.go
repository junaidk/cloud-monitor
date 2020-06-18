package main

import (
	"cloud-monitor/cloudman"
	"cloud-monitor/config"
	"cloud-monitor/driver"
	"cloud-monitor/printer"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) == 3 && os.Args[1] == "-c" {
		config.ConfPath = os.Args[2]
	} else {
		log.Println("provide config file path with -c parameter")
		return
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	d, err := driver.NewDriver()
	if err != nil {
		log.Println(err)
		return
	}

	type printStruct struct {
		Resp  []cloudman.InstanceListResponse
		Cloud string
	}
	printChan := make(chan printStruct, len(d.CloudList))
	for _, cloud := range d.CloudList {

		go func(c cloudman.Cloud) {
			list, _ := c.GetInstanceListAllRegions()
			printChan <- printStruct{
				Resp:  list,
				Cloud: strings.Replace(fmt.Sprintf("%T", c), "*cloudman.", "", -1),
			}
		}(cloud)

	}

	for j := 0; j < len(d.CloudList); j++ {
		// table printer
		data := <-printChan
		printer.PrintTable(data.Resp, data.Cloud)

	}

}
