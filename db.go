package ddbsync

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joshmyers/ddbsync/models"
)

type database struct {
	client    AWSDynamoer
	tableName string
}

func NewDatabase(tableName string, region string, endpoint string, disableSSL bool) DBer {
	return &database{
		client: dynamodb.New(session.New(&aws.Config{
			Endpoint:   &endpoint,
			Region:     &region,
			DisableSSL: &disableSSL,
		})),
		tableName: tableName,
	}
}

var _ DBer = (*database)(nil) // Forces compile time checking of the interface

var _ AWSDynamoer = (*dynamodb.DynamoDB)(nil) // Forces compile time checking of the interface

type DBer interface {
	Put(string, int64, int64) error
	Get(string) (*models.Item, error)
	Delete(string) error
}

type AWSDynamoer interface {
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
}

func (db *database) Put(name string, created int64, ttl int64) error {
	i := map[string]*dynamodb.AttributeValue{
		"Name": {
			S: aws.String(name),
		},
		"Created": {
			N: aws.String(strconv.FormatInt(created, 10)),
		},
		"TTL": {
			N: aws.String(strconv.FormatInt(ttl, 10)),
		},
	}

	e := map[string]*dynamodb.ExpectedAttributeValue{
		"Name": {
			Exists: aws.Bool(false),
		},
	}

	pit := &dynamodb.PutItemInput{
		TableName: aws.String(db.tableName),
		Item:      i,
		Expected:  e,
	}

	_, err := db.client.PutItem(pit)
	return err
}

func (db *database) Get(name string) (*models.Item, error) {
	kc := map[string]*dynamodb.Condition{
		"Name": {
			AttributeValueList: []*dynamodb.AttributeValue{
				{
					S: aws.String(name),
				},
			},
			ComparisonOperator: aws.String("EQ"),
		},
	}
	qi := &dynamodb.QueryInput{
		TableName:       aws.String(db.tableName),
		ConsistentRead:  aws.Bool(true),
		Select:          aws.String("SPECIFIC_ATTRIBUTES"),
		AttributesToGet: []*string{aws.String("Name"), aws.String("Created"), aws.String("TTL")},
		KeyConditions:   kc,
	}

	qo, err := db.client.Query(qi)
	if err != nil {
		return nil, err
	}

	// Make sure that no or 1 item is returned from DynamoDB
	if qo.Count != nil {
		if *qo.Count == 0 {
			return nil, fmt.Errorf("No item for Name, %s", name)
		} else if *qo.Count > 1 {
			return nil, fmt.Errorf("Expected only 1 item returned from Dynamo, got %d", *qo.Count)
		}
	} else {
		return nil, errors.New("Count not returned")
	}

	if len(qo.Items) < 1 || qo.Items[0] == nil {
		return nil, errors.New("No item returned, count is invalid.")
	}

	n := ""
	c := int64(0)
	t := int64(0)
	for index, element := range qo.Items[0] {
		if index == "Name" {
			n = *element.S
		}
		if index == "Created" {
			c, _ = strconv.ParseInt(*element.N, 10, 0)
		}
		if index == "TTL" {
			t, _ = strconv.ParseInt(*element.N, 10, 0)
		}
	}
	if n == "" || c == 0 || t == 0 {
		return nil, errors.New("The Name, Created and TTL keys were not found in the Dynamo result")
	}
	i := &models.Item{
		Name:    n,
		Created: c,
		TTL:     t,
	}
	return i, nil
}

func (db *database) Delete(name string) error {
	k := map[string]*dynamodb.AttributeValue{
		"Name": {
			S: aws.String(name),
		},
	}
	dii := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.tableName),
		Key:       k,
	}
	_, err := db.client.DeleteItem(dii)
	return err
}
