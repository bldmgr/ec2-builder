package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var m = make(map[string]bool)
var a = []string{}

func add(s string) {
	if m[s] {
		return // Already in the map
	}
	a = append(a, s)
	m[s] = true
}

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Long:  "List all available resources of the specified type",
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, _ := cmd.Flags().GetString("type")

		if resourceType == "" {
			resourceType = "all"
		}
		listAMI()
	},
}

func listAMI() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create EC2 service client
	svc := ec2.New(sess)

	runningInstances, err := GetRunningInstances(svc)
	if err != nil {
		fmt.Printf("Couldn't retrieve running instances: %v", err)
		return

	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "InstanceId", "PublicDnsName", "ImageId", "InstanceType", "LaunchTime", "Architecture", "PlatformDetails", "VpcId", "SubnetId"})
	dataShuttle := ""
	instanceName := ""
	for _, reservation := range runningInstances.Reservations {
		for _, instance := range reservation.Instances {
			dataShuttle = fmt.Sprintf("%s:%s", *instance.PublicDnsName, *instance.ImageId)

			add(dataShuttle)
			for i := 0; i < len(instance.Tags); i++ {
				if *instance.Tags[i].Key == "Name" {
					instanceName = *instance.Tags[i].Value
				}
			}

			t.AppendRows([]table.Row{{instanceName, *instance.InstanceId, *instance.PublicDnsName, *instance.ImageId, *instance.InstanceType, *instance.LaunchTime, *instance.Architecture, *instance.PlatformDetails, *instance.VpcId, *instance.SubnetId}})

		}
	}
	t.Render()
}
