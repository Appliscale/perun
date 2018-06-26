package progress

import (
	"encoding/json"
	"errors"
	"os/user"
	"strings"
	"time"

	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/mysession"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/olekukonko/tablewriter"
)

type Connection struct {
	context *context.Context
	session *session.Session

	SqsClient          *sqs.SQS
	sqsQueueOutput     *sqs.CreateQueueOutput
	sqsQueueAttributes *sqs.GetQueueAttributesOutput

	snsClient *sns.SNS
	TopicArn  *string
}

var sinkName = "perun-sink-"

const awsTimestampLayout = "2006-01-02T15:04:05.000Z"

// Configure AWS Resources needed for pregress monitoring
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
			context.Logger.Warning("It's configuration may take up to a minute, wait before calling 'create-stack' with flag --progress")
		}
		return
	}
}

// Remove all AWS Resources created for stack monitoring
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

// Get configuration of created AWS Resources
func GetRemoteSink(context *context.Context, session *session.Session) (conn Connection, err error) {
	conn = initMessageService(context, session)
	snsTopicExists, sqsQueueExists, err := conn.verifyRemoteSinkConfigured()
	if !(snsTopicExists && sqsQueueExists) {
		err = errors.New("remote sink has not been configured, run 'perun setup-remote-sink' first. If You done it already, wait for aws sink configuration")
		return
	}
	return
}

func initRemoteConnection(context *context.Context) Connection {
	currentSession := initSession(context)
	return initMessageService(context, currentSession)

}
func initMessageService(context *context.Context, currentSession *session.Session) (conn Connection) {
	currentUser, userError := user.Current()
	if userError != nil {
		context.Logger.Error("error reading currentUser")
	}
	sinkName += currentUser.Username + "-" + currentUser.Uid
	conn.session = currentSession
	conn.context = context
	return
}

func initSession(context *context.Context) *session.Session {
	tokenError := mysession.UpdateSessionToken(context.Config.DefaultProfile, context.Config.DefaultRegion, context.Config.DefaultDurationForMFA, context)
	if tokenError != nil {
		context.Logger.Error(tokenError.Error())
	}
	currentSession, createSessionError := mysession.CreateSession(context, context.Config.DefaultProfile, &context.Config.DefaultRegion)
	if createSessionError != nil {
		context.Logger.Error(createSessionError.Error())
	}
	return currentSession
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

// Monitor queue, that delivers messages sent by cloud formation stack progress
func (conn *Connection) MonitorQueue() {
	waitTimeSeconds := int64(3)
	receiveMessageInput := sqs.ReceiveMessageInput{
		QueueUrl:        conn.sqsQueueOutput.QueueUrl,
		WaitTimeSeconds: &waitTimeSeconds,
	}

	pw, table := initTableWriter()

	tolerance, err := time.ParseDuration("1s")
	if err != nil {
		conn.context.Logger.Error(err.Error())
	}
	startReadingMessagesTime := time.Now().Add(-tolerance)

	receivedAllMessages := true
	for receivedAllMessages {
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
			}
			// Check if the message has been the last one (status COMPLETE for current stack resource)
			if strings.Contains(messageMap["LogicalResourceId"], *conn.context.CliArguments.Stack) &&
				strings.Contains(messageMap["ResourceStatus"], "COMPLETE") {
				receivedAllMessages = false
			}
		}
	}
}
func initTableWriter() (*parseWriter, *tablewriter.Table) {
	pw := newParseWriter()
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
	conn.snsClient = sns.New(conn.session)
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
	conn.SqsClient = sqs.New(conn.session)

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
	conn.snsClient = sns.New(conn.session)

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
	conn.SqsClient = sqs.New(conn.session)

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
