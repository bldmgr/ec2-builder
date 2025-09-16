package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	b64 "encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new resource",
	Long:  "Create a new resource with the specified parameters",
	Run: func(cmd *cobra.Command, args []string) {
		instanceType, _ := cmd.Flags().GetString("type")
		imageId, _ := cmd.Flags().GetString("image")
		name, _ := cmd.Flags().GetString("name")

		if instanceType == "" {
			fmt.Println("Error: resource type is required")
			cmd.Help()
			os.Exit(1)
		}

		if name == "" {
			fmt.Println("Error: name is required")
			cmd.Help()
			os.Exit(1)
		}

		amiID := lanuchEc2(imageId, name, instanceType)
		fmt.Printf("Successfully created %s: %s\n", instanceType, name)
		getAMI(amiID)

	},
}

func lanuchEc2(aws_ami string, nameTag string, instanceType string) (aimID string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create EC2 service client
	svc := ec2.New(sess)

	userDataScript := `#cloud-config
system_info:
  default_user:
    name: "circleci-admin"
ssh_authorized_keys:
  - "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDXk086YfVp2jr5JS0Vytb68if3WCcnTlAFZqUts3Yzpjn3VOf8bee6BmU1SVt6uYt97KLgBQ3rTNE0zzhQvfpNz6Gf1ji5d3J6PYWvExQlzCtGdb59b/UuxKA9QyMAwU6D9DcQnFNEmjSBIFVUYfcMZhPRo63Zho3CTt21lQMQZEZEzTdTWjLT8XoM611lC4cuKhSKfB13kEMDSlgEQuTkPNoMh4w4xD5q4JVbSvBzGiEPt3ArnyEK6RjIybx+ucjYeoc+hj2j0w5TgGpJkuEt61tIdhl0GZSDfTTThPelhLYw8ym3gvAKYSM06UROI/dPR84kQ0MG+XwrZHgYFBQAbpH0on13LR9g3bZEBY8uwGO5LhhNg2o5WcP51U6d44FWCq0KANB/XXiEljQ9dN4d1rk2yjp3r9GIQT8h63Vu05FfFNYeTJwU/biA19BYaqlPlejUZrCZF2/LwBoSWohhxTdirVzf9lvHp57Sgvur5g4ivVgGnQ+Y1yjH3H60hboFZrYXiGs5MAz3dg7MMZ4gkwNOKNNHGEur55s5IwlExD/1j+eUqc/hQlltw3yT2GKjdOU/a3GiG7J/GtF6la+7XHd1oTuIfhLOYdb4vqSqGFcW6YuTphLhyFgrp1lw/NevtVpPCNhpmfNXQy6H9fw5hgmFdXU/Ct/SO+tOTGzA0Q== circleci-admin"
bootcmd:
  - gpasswd -d circleci sudo`

	sEnc := b64.StdEncoding.EncodeToString([]byte(userDataScript))
	userData := aws.String(sEnc)

	// Specify the details of the inls tance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(aws_ami),
		InstanceType: aws.String(instanceType),
		KeyName:      aws.String("us-east-1"),
		UserData:     userData,
		SecurityGroupIds: []*string{
			aws.String("sg-0bc572cd216ea3f54"),
		},
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(200),
				},
			},
		},

		MinCount: aws.Int64(1),
		MaxCount: aws.Int64(1),
	})

	if err != nil {
		fmt.Println("Could not create instance", err)
		return
	}

	fmt.Println("Created instance", *runResult.Instances[0].InstanceId)

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(nameTag),
			},
		},
	})
	if errtag != nil {
		log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		return
	}

	fmt.Println("Successfully tagged instance")

	return *runResult.Instances[0].InstanceId
}

func getAMI(imageId string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create EC2 service client
	svc := ec2.New(sess)

	time.Sleep(90 * time.Second)

	runningInstances, err := GetRunningInstances(svc)
	if err != nil {
		fmt.Printf("Couldn't retrieve running instances: %v", err)
		return
	}

	for _, reservation := range runningInstances.Reservations {
		for _, instance := range reservation.Instances {
			if imageId == *instance.InstanceId {
				fmt.Printf("Found running instance: %s\n", *instance.PublicDnsName)
				fmt.Printf(" ssh -i \"belkin\" circleci-admin@%s\n", *instance.PublicDnsName)
			}
		}
	}
}

func GetRunningInstances(client *ec2.EC2) (*ec2.DescribeInstancesOutput, error) {
	result, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return result, err
}
