package jsonapi

// ErrorCode for err code
type ErrorCode int

// validation
const (
	ErrValidation ErrorCode = iota + 1000
	ErrInvalidRequestData
)

// database
const (
	ErrDatabase ErrorCode = iota + 2000
	ErrInsertDatabase
	ErrUpdateDatabase
	ErrDeleteDatabase
	ErrSelectDatabase
	ErrNoResultFromDatabase
	ErrStartTransaction
	ErrCommitTransaction
	ErrRollbackTransaction
)

// third party library or api or something
const (
	ErrThirdParty ErrorCode = iota + 3000
	ErrOperationAPI
)

// caching memcache, redis, koding cache, etc . . .
const (
	ErrCaching ErrorCode = iota + 4000
	ErrRedisGetData
	ErrRedisSetData
	ErrMemcacheGetData
	ErrMemcacheSetData
)

// grpc
const (
	ErrGRPC ErrorCode = iota + 5000
	ErrFailedRequestGRPC
	ErrFailedResponseGRPC
)

// elastic search
const (
	ErrElasticsearch ErrorCode = iota + 6000
	ErrFailedUpdateIndexES
	ErrFailedGetFromES
)

// message queue
const (
	ErrMessageQueue ErrorCode = iota + 7000
	ErrRabbitMQFailedPublishMessage
	ErrRabbitMQFailedConsumeMessage
	ErrGooglePubSubFailedPublishMessage
	ErrGooglePubSubFailedConsumeMessage
)

// others...
const (
	ErrNotCategorized ErrorCode = 157
)
