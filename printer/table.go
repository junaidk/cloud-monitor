package printer

import (
	"cloud-monitor/cloudman"
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
	"strconv"
	"strings"
)

func PrintTable(data []cloudman.InstanceListResponse, cloud string) {

	// sort on state
	sort.Slice(data, func(i, j int) bool {
		return data[i].Status > data[j].Status
	})

	running := 0
	var out [][]string
	for _, item := range data {

		if strings.ToLower(item.Status) == "running" || strings.ToLower(item.Status) == "succeeded" || strings.ToLower(item.Status) == "active" {
			running += 1
		}
		row := []string{
			item.Name,
			item.Status,
			item.Region,
			item.MachineType,
			item.Id,
			item.LaunchDate,
			item.Project,
		}
		out = append(out, row)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Status", "Region", "Type", "Id", "Time", "Project"})
	table.SetFooter([]string{"Cloud", cloud, "", "", "", "Running", strconv.Itoa(running)})
	//table.SetBorder(false) // Set Border to false
	table.AppendBulk(out) // Add Bulk Data
	table.Render()

}
