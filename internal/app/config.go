package app

import (
	"github.com/core-go/health/server"
	"github.com/core-go/ibmmq"
	"github.com/core-go/mq"
	"github.com/core-go/mq/zap"
)

type Config struct {
	Server  server.ServerConf `mapstructure:"server"`
	Log     log.Config        `mapstructure:"log"`
	Mongo   MongoConfig       `mapstructure:"mongo"`
	IBMMQ   IBMMQConfig       `mapstructure:"ibmmq"`
	Handler mq.HandlerConfig  `mapstructure:"handler"`
}

type MongoConfig struct {
	Uri      string `yaml:"uri" mapstructure:"uri" json:"uri,omitempty" gorm:"column:uri" bson:"uri,omitempty" dynamodbav:"uri,omitempty" firestore:"uri,omitempty"`
	Database string `yaml:"database" mapstructure:"database" json:"database,omitempty" gorm:"column:database" bson:"database,omitempty" dynamodbav:"database,omitempty" firestore:"database,omitempty"`
}

type IBMMQConfig struct {
	QueueConfig      ibmmq.QueueConfig      `mapstructure:"queue_config"`
	SubscriberConfig ibmmq.SubscriberConfig `mapstructure:"subscriber_config"`
	MQAuth           ibmmq.MQAuth           `mapstructure:"mq_auth"`
}
