package apihandler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"sort"
	"strings"
)

func indexChecker(input []int32) error {
	// check Index validation
	sort.Slice(input , func(i, j int) bool { return input[i] < input[j] })
	for i, i2 := range input {
		if int32(i) != i2 {
			return status.Errorf(codes.InvalidArgument, "Index %d missing. Please check indexes.", i)
		}
	}
	return nil
}

func pageIdGenerator() string {
	bytesArray, _ := generateRandomBytes(32)
	hasher := md5.New()
	hasher.Write(bytesArray)
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func isMongoKeyExists(tileColl *mongo.Collection, key string, ctx context.Context) bool {
	if count, _ := tileColl.CountDocuments(ctx, bson.M{key : bson.M{ "$exists" : true}}); count <= 0 {
		return false
	}
	return true
}

func trimSpaces(values []string ) []string {
	var result []string
	for _, value := range values {
		result = append(result, strings.TrimSpace(value))
	}
	return result
}























