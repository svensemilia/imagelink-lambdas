package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/svensemilia/util"
)

var (
	svc      *ec2.EC2
	input    *ec2.DescribeInstancesInput
	response *util.LambdaApiResponse
	state    State
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

	response = &util.LambdaApiResponse{Code: 200, Headers: map[string]string{"hello": "world"}}
	state = State{}
}

type State struct {
	State *ec2.InstanceState `json:"state"`
	IP    *string            `json:"ip,omitempty"`
	Error string             `json:"error,omitempty"`
}

func HandleRequest() (util.LambdaApiResponse, error) {
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
		response.Code = 500
		state.Error = err.Error()
		return util.StringifyBody(response, state), nil
	}
	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		response.Code = 404
		state.Error = "No instances found for your request - please contact the administrator"
		return util.StringifyBody(response, state), nil
	}
	state.IP = result.Reservations[0].Instances[0].PublicIpAddress
	state.State = result.Reservations[0].Instances[0].State
	return util.StringifyBody(response, state), nil
}

func main() {
	lambda.Start(HandleRequest)
}
