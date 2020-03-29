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
	"github.com/svensemilia/util"
)

const (
	actionStart = "start"
	actionStop  = "stop"
)

var (
	svc        *ec2.EC2
	inputStart *ec2.StartInstancesInput
	inputStop  *ec2.StopInstancesInput
	response   *util.LambdaApiResponse
	state      *State
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

	response = &util.LambdaApiResponse{Code: 200, Headers: map[string]string{"hello": "world"}}
	state = &State{}
}

type InstanceAction struct {
	Action string `json:"action"`
}

type State struct {
	Action string `json:"action,omitempty"`
	Error  string `json:"error,omitempty"`
}

func HandleRequest(request util.LambdaApiRequest) (util.LambdaApiResponse, error) {
	var bodyJson InstanceAction
	var err error

	err = json.Unmarshal([]byte(request.Body), &bodyJson)
	if err != nil {
		response.Code = 500
		state.Error = err.Error()
		return util.StringifyBody(response, *state), nil
	}

	if bodyJson.Action == actionStart {
		fmt.Println("Starting instance...")
		err = startInstance()
	} else if bodyJson.Action == actionStop {
		err = stopInstance()
		fmt.Println("Stopping instance...")
	} else {
		fmt.Println("Action type not recognized:", bodyJson.Action)
		err = fmt.Errorf("Action type not recognized: %s", bodyJson.Action)
	}
	if err != nil {
		response.Code = 500
		state.Error = err.Error()
	}
	state.Action = bodyJson.Action
	return util.StringifyBody(response, *state), nil
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
