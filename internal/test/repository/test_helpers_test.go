package repository

import (
	"testing"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func newMongoTest(t *testing.T) *mtest.T {
	return mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
}
