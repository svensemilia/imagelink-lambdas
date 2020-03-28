package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	svc      *ec2.EC2
	input    *ec2.DescribeInstancesInput
	response *StatusResponse
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

	response = &StatusResponse{Code: 200, Headers: map[string]string{"hello": "world"}}
	state = State{}
}

type StatusResponse struct {
	Code    int               `json:"statusCode"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Base64  bool              `json:"isBase64Encoded"`
}

type State struct {
	State *ec2.InstanceState `json:"state"`
	IP    *string            `json:"ip,omitempty"`
	Error error              `json:"error,omitempty"`
}

func HandleRequest() (StatusResponse, error) {
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
		state.Error = err
		return stringifyJson(response, state), err
	}
	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		response.Code = 404
		return stringifyJson(response, state), fmt.Errorf("No instances found for your request - please contact the administrator")
	}
	state.IP = result.Reservations[0].Instances[0].PublicIpAddress
	state.State = result.Reservations[0].Instances[0].State
	return stringifyJson(response, state), nil
}

func stringifyJson(response *StatusResponse, state State) StatusResponse {
	b, _ := json.Marshal(state)
	response.Body = string(b)
	return *response
}

func main() {
	lambda.Start(HandleRequest)
}
