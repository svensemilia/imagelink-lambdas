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

const(
	actionStart = "start"
	actionStop = "stop" 
)

var (
	svc   *ec2.EC2
	inputStart *ec2.StartInstancesInput
	inputStop *ec2.StopInstancesInput
)

func init() {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	svc = ec2.New(sess)

	instanceId := os.Getenv("instance_id")
	fmt.Println("Instance name", instanceId)
	inputStart = &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}
	inputStop = &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}
}

type InstanceAction struct {
	Action string `json:"action"`
}

func HandleRequest(action InstanceAction) error {
	var err error
	if action.Action == actionStart {
		fmt.Println("Starting instance...")
		err = startInstance()
	} else if action.Action == actionStop {
		err = stopInstance()
		fmt.Println("Stopping instance...")
	} else {
		fmt.Println("Action type not recognized:", action.Action)
	}
	return err
}

func startInstance() error {
	_, err := svc.StartInstances(inputStart)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	return err
}

func stopInstance() error {
	_, err := svc.StopInstances(inputStop)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	return err
}

func main() {
	lambda.Start(HandleRequest)
}