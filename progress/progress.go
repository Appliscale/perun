// Copyright 2018 Appliscale
//
// Maintainers and contributors are listed in README file inside repository.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package progress provides displaying of progress e.g during stack creation.
package progress

import (
	"encoding/json"
	"errors"
	"os/user"
	"strings"
	"time"

	"github.com/Appliscale/perun/context"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/olekukonko/tablewriter"
)

// Connection contains elements need to get connection.
type Connection struct {
	context *context.Context

	SqsClient          *sqs.SQS
	sqsQueueOutput     *sqs.CreateQueueOutput
	sqsQueueAttributes *sqs.GetQueueAttributesOutput

	snsClient *sns.SNS
	TopicArn  *string
}

var sinkName = "perun-sink-"

const awsTimestampLayout = "2006-01-02T15:04:05.000Z"

// Configure AWS Resources needed for progress monitoring.
func ConfigureRemoteSink(context *context.Context) (err error) {
	conn := initRemoteConnection(context)
	snsTopicExists, sqsQueueExists, err := conn.verifyRemoteSinkConfigured()

	if snsTopicExists && sqsQueueExists {
		context.Logger.Info("Remote sink has already been configured")
		return
	} else {
		shouldSNSTopicBeRemoved := snsTopicExists && !sqsQueueExists
		err = conn.deleteRemainingSinkResources(shouldSNSTopicBeRemoved, false)
		if err != nil {
			context.Logger.Error("error deleting up remote sink: " + err.Error())
			return
		}

		if !sqsQueueExists {
			err = conn.setUpSQSQueue()
			if err != nil {
				context.Logger.Error("Error creating sqs queue: " + err.Error())
				return
			}
		}

		if shouldSNSTopicBeRemoved || !snsTopicExists { // SNS Topic has been removed or does not exist
			err = conn.setUpSNSNotification()
			if err != nil {
				context.Logger.Error("Error creating sqs queue: " + err.Error())
				return
			}
		}

		if err == nil {
			context.Logger.Info("Remote sink configuration successful")
			context.Logger.Warning("It's configuration may take up to a minute, wait before using Perun with flag --progress")
		}
		return
	}
}

// Remove all AWS Resources created for stack monitoring.
func DestroyRemoteSink(context *context.Context) (conn Connection, err error) {
	conn = initRemoteConnection(context)
	snsTopicExists, sqsQueueExists, err := conn.verifyRemoteSinkConfigured()
	if err != nil {
		context.Logger.Error("error verifying: " + err.Error())
	}

	if !(snsTopicExists && sqsQueueExists) {
		err = errors.New("remote sink has not been configured or has already been deleted")
		return
	} else {
		err = conn.deleteRemainingSinkResources(snsTopicExists, sqsQueueExists)
		if err != nil {
			return
		}
		context.Logger.Info("Remote sink deconstruction successful.")
		return
	}
}

// Get configuration of created AWS Resources.
func GetRemoteSink(context *context.Context) (conn Connection, err error) {
	conn = initMessageService(context)
	snsTopicExists, sqsQueueExists, err := conn.verifyRemoteSinkConfigured()
	if !(snsTopicExists && sqsQueueExists) {
		err = errors.New("remote sink has not been configured, run 'perun setup-remote-sink' first. If You done it already, wait for aws sink configuration")
		return
	}
	return
}

func initRemoteConnection(context *context.Context) Connection {
	context.InitializeAwsAPI()
	return initMessageService(context)

}
func initMessageService(context *context.Context) (conn Connection) {
	currentUser, userError := user.Current()
	if userError != nil {
		context.Logger.Error("error reading currentUser")
	}
	sinkName += currentUser.Username + "-" + currentUser.Uid
	conn.context = context
	return
}

func (conn *Connection) verifyRemoteSinkConfigured() (snsTopicExists bool, sqsQueueExists bool, err error) {
	snsTopicExists, err = conn.getSnsTopicAttributes()
	if err != nil {
		conn.context.Logger.Error("Error getting sns topic configuration: " + err.Error())
	}
	sqsQueueExists, err = conn.getSqsQueueAttributes()
	if err != nil {
		conn.context.Logger.Error("Error getting sqs queue configuration: " + err.Error())
	}
	return
}

// Message - struct with elements of message.
type Message struct {
	Type             string
	MessageId        string
	TopicArn         string
	Subject          string
	Message          string
	Timestamp        string
	SignatureVersion string
	Signature        string
	SigningCertURL   string
	UnsubscribeURL   string
}

// Monitor queue, that delivers messages sent by cloud formation stack progress.
func (conn *Connection) MonitorStackQueue() {
	waitTimeSeconds := int64(3)
	receiveMessageInput := sqs.ReceiveMessageInput{
		QueueUrl:        conn.sqsQueueOutput.QueueUrl,
		WaitTimeSeconds: &waitTimeSeconds,
	}

	pw, table := initStackTableWriter()

	tolerance, err := time.ParseDuration("1s")
	if err != nil {
		conn.context.Logger.Error(err.Error())
	}
	startReadingMessagesTime := time.Now().Add(-tolerance)

	AnyMessagesLeft := true
	for AnyMessagesLeft {
		receivedMessages, err := conn.SqsClient.ReceiveMessage(&receiveMessageInput)
		if err != nil {
			conn.context.Logger.Error("Error reading messages: " + err.Error())
		}
		for e := range receivedMessages.Messages {
			v := Message{}
			jsonBlob := []byte(*receivedMessages.Messages[e].Body)
			err = json.Unmarshal(jsonBlob, &v)
			if err != nil {
				conn.context.Logger.Error("error reading json message" + err.Error())
			}

			// DELETE READ MESSAGE (to prevent reading the same message multiple times)
			conn.SqsClient.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      conn.sqsQueueOutput.QueueUrl,
				ReceiptHandle: receivedMessages.Messages[e].ReceiptHandle,
			})

			// Parse property message
			splittedMessage := strings.FieldsFunc(v.Message, func(r rune) bool { return r == '\n' })
			messageMap := map[string]string{}
			for messageNum := range splittedMessage {
				messages := strings.FieldsFunc(splittedMessage[messageNum], func(r rune) bool { return r == '=' })
				messageMap[messages[0]] = messages[1]
			}
			// Parse timestamp of message
			messageArrivedTime, err := time.Parse(awsTimestampLayout, v.Timestamp)
			if err != nil {
				conn.context.Logger.Error(err.Error())
			}

			if startReadingMessagesTime.Before(messageArrivedTime) {
				table.Append([]string{v.Timestamp, messageMap["ResourceStatus"], messageMap["ResourceType"], messageMap["LogicalResourceId"], messageMap["ResourceStatusReason"]})
				pw.returnWritten()
				table.Render()
				// Check if the message has been the last one (status COMPLETE for current stack resource)
				if strings.Contains(messageMap["LogicalResourceId"], *conn.context.CliArguments.Stack) &&
					strings.Contains(messageMap["ResourceStatus"], "COMPLETE") {
					AnyMessagesLeft = false
				}
			}
		}
	}
}
func initStackTableWriter() (*ParseWriter, *tablewriter.Table) {
	pw := NewParseWriter()
	table := tablewriter.NewWriter(pw)
	table.SetHeader([]string{"Time", "Status", "Type", "LogicalID", "Status Reason"})
	table.SetBorder(false)
	// Set Border to false
	table.SetColumnColor(tablewriter.Colors{tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgWhiteColor})
	return pw, table
}

func (conn *Connection) setUpSNSNotification() (err error) {
	//CREATE SNS TOPIC
	conn.snsClient = sns.New(conn.context.CurrentSession)
	topicInput := sns.CreateTopicInput{
		Name: &sinkName,
	}
	topicOutput, _ := conn.snsClient.CreateTopic(&topicInput)
	conn.TopicArn = topicOutput.TopicArn

	//SET UP POLICY
	err = conn.setUpSqsPolicy()

	protocolSQS := "sqs"
	subscribeInput := sns.SubscribeInput{
		Endpoint: conn.sqsQueueAttributes.Attributes[sqs.QueueAttributeNameQueueArn],
		Protocol: &protocolSQS,
		TopicArn: conn.TopicArn,
	}
	conn.snsClient.Subscribe(&subscribeInput)

	conn.context.Logger.Info("Set up SNS Notification topic: " + sinkName)
	return
}
func (conn *Connection) setUpSQSQueue() (err error) {
	conn.SqsClient = sqs.New(conn.context.CurrentSession)

	sixtySec := "60"
	sqsInput := sqs.CreateQueueInput{
		QueueName: &sinkName,
		Attributes: map[string]*string{
			"MessageRetentionPeriod": &sixtySec,
		},
	}
	conn.sqsQueueOutput, err = conn.SqsClient.CreateQueue(&sqsInput)
	if err != nil {
		return
	}

	arnAttribute := sqs.QueueAttributeNameAll
	queueAttributesInput := sqs.GetQueueAttributesInput{
		AttributeNames: []*string{&arnAttribute},
		QueueUrl:       conn.sqsQueueOutput.QueueUrl,
	}
	conn.sqsQueueAttributes, err = conn.SqsClient.GetQueueAttributes(&queueAttributesInput)
	if err != nil {
		return
	}
	conn.context.Logger.Info("Set up SQS Notification Queue: " + sinkName)
	return
}

func (conn *Connection) setUpSqsPolicy() (err error) {
	jsonStringPolicy, err := conn.createJsonPolicy()
	if err != nil {
		conn.context.Logger.Error("error creating json: " + err.Error())
	}

	queueAttributes := sqs.SetQueueAttributesInput{
		QueueUrl: conn.sqsQueueOutput.QueueUrl,
		Attributes: map[string]*string{
			sqs.QueueAttributeNamePolicy: &jsonStringPolicy,
		},
	}
	conn.SqsClient.SetQueueAttributes(&queueAttributes)

	conn.context.Logger.Info("Created SQS access policy for SNS Topic: " + sinkName)
	return
}

type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}
type StatementEntry struct {
	Sid       string
	Effect    string
	Action    []string
	Resource  string
	Condition Condition
	Principal string
}
type Condition struct {
	StringEquals map[string]string
}

func (conn *Connection) createJsonPolicy() (jsonStringPolicy string, err error) {
	policy := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			{
				Effect: "Allow",
				Action: []string{
					"SQS:*",
				},
				Resource: *conn.sqsQueueAttributes.Attributes[sqs.QueueAttributeNameQueueArn],
				Condition: Condition{
					StringEquals: map[string]string{"aws:SourceArn": *conn.TopicArn},
				},
				Principal: "*",
			},
		},
	}

	jsonPolicy, err := json.Marshal(policy)
	jsonStringPolicy = string(jsonPolicy)
	return
}
func (conn *Connection) getSnsTopicAttributes() (topicExists bool, err error) {
	conn.snsClient = sns.New(conn.context.CurrentSession)

	topicExists = false
	listTopicsInput := sns.ListTopicsInput{}
	err = conn.snsClient.ListTopicsPages(&listTopicsInput,
		func(output *sns.ListTopicsOutput, lastPage bool) bool {
			for topicNum := range output.Topics {
				if strings.Contains(*output.Topics[topicNum].TopicArn, sinkName) {
					topicExists = true
					conn.TopicArn = output.Topics[topicNum].TopicArn
					return false
				}
			}
			return true
		})
	return
}
func (conn *Connection) getSqsQueueAttributes() (queueExists bool, err error) {
	conn.SqsClient = sqs.New(conn.context.CurrentSession)

	queueExists = false
	listQueuesInput := sqs.ListQueuesInput{
		QueueNamePrefix: &sinkName,
	}

	listQueuesOutput, err := conn.SqsClient.ListQueues(&listQueuesInput)

	for queueNum := range listQueuesOutput.QueueUrls {
		if strings.Contains(*listQueuesOutput.QueueUrls[queueNum], sinkName) {
			queueExists = true
			conn.sqsQueueOutput = &sqs.CreateQueueOutput{
				QueueUrl: listQueuesOutput.QueueUrls[queueNum],
			}

			arnAttribute := sqs.QueueAttributeNameQueueArn
			queueAttributesInput := sqs.GetQueueAttributesInput{
				AttributeNames: []*string{&arnAttribute},
				QueueUrl:       conn.sqsQueueOutput.QueueUrl,
			}
			conn.sqsQueueAttributes, err = conn.SqsClient.GetQueueAttributes(&queueAttributesInput)
			if err != nil {
				return
			}
			return
		}
	}
	return
}

func (conn *Connection) deleteSnsTopic() (err error) {
	deleteTopicInput := sns.DeleteTopicInput{
		TopicArn: conn.TopicArn,
	}
	_, err = conn.snsClient.DeleteTopic(&deleteTopicInput)
	return
}
func (conn *Connection) deleteSqsQueue() (err error) {
	deleteQueueInput := sqs.DeleteQueueInput{
		QueueUrl: conn.sqsQueueOutput.QueueUrl,
	}
	_, err = conn.SqsClient.DeleteQueue(&deleteQueueInput)
	return
}

func (conn *Connection) deleteRemainingSinkResources(deleteSnsTopic bool, deleteSqsQueue bool) (err error) {
	if deleteSnsTopic {
		err = conn.deleteSnsTopic()
		if err != nil {
			conn.context.Logger.Error("Error deleting sns Topic: " + err.Error())
			return
		}
		conn.context.Logger.Info("Deleting SNS Topic")
	}
	if deleteSqsQueue {
		err = conn.deleteSqsQueue()
		if err != nil {
			conn.context.Logger.Error("Error deleting sqs Queue: " + err.Error())
			return
		}
		conn.context.Logger.Info("Deleting SQS Queue")
	}
	return
}
