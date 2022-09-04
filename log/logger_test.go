package log

import (
	"context"
	"log"
	"testing"

	"github.com/c2fo/testify/assert"
)

type obj struct {
	Name    string `json:"name"`
	Count   int    `json:"count"`
	Enabled bool   `json:"enabled"`
}

var (
	sampleObjects = []*obj{
		{"a", 1, true},
		{"b", 2, false},
		{"c", 3, true},
		{"d", 4, false},
		{"e", 5, true},
		{"f", 6, false},
		{"g", 7, true},
		{"h", 8, false},
		{"i", 9, true},
		{"j", 0, false},
	}
	sampleArray  = make([]int, 10000)
	sampleString = "some string with a somewhat realistic length"
)

var (
	sampleContext context.Context
	logger        SLogger
)

func init() {
	logger = NewSLogger(sampleString)
	sampleContext = logger.BuildContextDataAndSetValue("Japan", "11")
}

func BenchmarkSLogNative_InfoSimple(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		log.Println(sampleContext, sampleString)
	}
}

func BenchmarkSLog_InfoSimple(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info(sampleContext, sampleString)
	}
}

func BenchmarkSLogNative_LogLargeArray(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		log.Println(sampleContext, sampleArray)
	}
}

func BenchmarkSLog_LogLargeArray(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info(sampleContext, sampleArray)
	}
}

func BenchmarkSLogNative_InfoWithComplexArgs(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		log.Println(sampleContext, sampleObjects)
	}
}

func BenchmarkSLog_InfoWithComplexArgs(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info(sampleContext, sampleObjects)
	}
}

func BenchmarkSLogNative_InfofWithComplexArgs(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		log.Println(sampleContext, sampleString)
	}
}

func BenchmarkSLog_InfofWithComplexArgs(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Infof(sampleContext, sampleString)
	}
}

func TestNoCollisionWhenBuildContextData(t *testing.T) {
	type thisKeyType string
	var (
		thisKey      thisKeyType = thisKeyType(ContextDataMapKey) // to ensure the value equal
		thisKeyValue             = "fakhri"

		randomCountry = "ID"
		randomID      = "randomID"
	)

	ctx := logger.BuildContextDataAndSetValue(randomCountry, randomID)
	newCtx := context.WithValue(ctx, thisKey, thisKeyValue)
	contextDataFromLogger := newCtx.Value(ContextDataMapKey).(map[string]string)

	// Ensure there is no collision
	assert.Equal(t, thisKeyValue, newCtx.Value(thisKey))
	assert.Equal(t, randomCountry, contextDataFromLogger[ContextCountryKey])
	assert.Equal(t, randomID, contextDataFromLogger[ContextIdKey])
}
