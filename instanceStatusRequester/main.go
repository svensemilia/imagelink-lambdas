package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	svc   *ec2.EC2
	input *ec2.DescribeInstancesInput
)

func init() {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	svc = ec2.New(sess)

	instanceId := os.Getenv("instance_id")
	fmt.Println("Instance name", instanceId)
	input = &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}
}

type StatusResult struct {
	State *ec2.InstanceState `json:"state"`
	IP    *string            `json:"ip"`
}

func HandleRequest() (StatusResult, error) {
	output := StatusResult{}
	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return output, err
	}
	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return output, fmt.Errorf("No instances found for your request - please contact the administrator")
	}
	output.IP = result.Reservations[0].Instances[0].PublicIpAddress
	output.State = result.Reservations[0].Instances[0].State
	return output, nil
}

func main() {
	lambda.Start(HandleRequest)
}
