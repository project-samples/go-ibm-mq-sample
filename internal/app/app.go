package app

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/core-go/health"
	hm "github.com/core-go/health/mongo"
	"github.com/core-go/ibmmq"
	w "github.com/core-go/mongo/writer"
	"github.com/core-go/mq"
	v "github.com/core-go/mq/validator"
	"github.com/core-go/mq/zap"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	Subscribe     func(ctx context.Context, handle func(context.Context, []byte))
	Handle        func(context.Context, []byte)
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	log.Initialize(cfg.Log)
	client, er1 := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.Uri))
	if er1 != nil {
		log.Error(ctx, "Cannot connect to MongoDB: Error: "+er1.Error())
		return nil, er1
	}
	db := client.Database(cfg.Mongo.Database)

	logError := log.ErrorMsg
	var logInfo func(context.Context, string)
	if log.IsInfoEnable() {
		logInfo = log.InfoMsg
	}

	subscriber, er2 := ibmmq.NewSubscriberByConfig(cfg.IBMMQ.SubscriberConfig, cfg.IBMMQ.MQAuth, log.ErrorMsg)
	if er2 != nil {
		log.Error(ctx, "Cannot create a new subscriber. Error: "+er2.Error())
	}
	validator, err := v.NewValidator[*User]()
	if err != nil {
		return nil, err
	}
	errorHandler := mq.NewErrorHandler[*User](logError)
	writer := w.NewWriter[*User](db, "user")
	handler := mq.NewHandlerByConfig[User](cfg.Handler, writer.Write, validator.Validate, errorHandler.Reject, errorHandler.HandleError, logError, logInfo)
	mongoChecker := hm.NewHealthChecker(client)
	subscriberChecker := ibmmq.NewHealthCheckerByConfig(&cfg.IBMMQ.QueueConfig, &cfg.IBMMQ.MQAuth, cfg.IBMMQ.SubscriberConfig.Topic)
	healthHandler := health.NewHandler(mongoChecker, subscriberChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		Subscribe:     subscriber.Subscribe,
		Handle:        handler.Handle,
	}, nil
}
