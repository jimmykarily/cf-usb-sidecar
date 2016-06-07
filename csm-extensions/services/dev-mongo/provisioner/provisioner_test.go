package provisioner

import (
	"log"
	"os"
	"testing"

	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mongo/config"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("mongo-provisioner-test")

var mongoConConfig = struct {
	User            string
	Pass            string
	Host            string
	Port            string
	TestProvisioner MongoProvisionerInterface
}{}

func init() {
	mongoConConfig.User = os.Getenv("MONGO_USER")
	mongoConConfig.Pass = os.Getenv("MONGO_PASS")
	mongoConConfig.Host = os.Getenv("MONGO_HOST")
	mongoConConfig.Port = os.Getenv("MONGO_PORT")

	mongo := config.MongoDriverConfig{
		Host: mongoConConfig.Host,
		Port: mongoConConfig.Port,
		Pass: mongoConConfig.Pass,
		User: mongoConConfig.User,
	}

	mongoConConfig.TestProvisioner = New(mongo, logger)
}

func TestCreateDb(t *testing.T) {
	assert := assert.New(t)

	dbName := "test_createdb"
	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	log.Println("Creating test database")
	err := mongoConConfig.TestProvisioner.CreateDatabase(dbName)

	assert.Nil(err)
}

func TestCreateDbExists(t *testing.T) {
	assert := assert.New(t)
	dbName := "test_createdb"

	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	log.Println("Testing if database exists")
	created, err := mongoConConfig.TestProvisioner.IsDatabaseCreated(dbName)
	assert.Nil(err)
	assert.True(created)
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	dbName := "test_createdb"

	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	log.Println("Creating test user")
	err := mongoConConfig.TestProvisioner.CreateUser(dbName, "mytestUser", "mytestPass")
	assert.Nil(err)
}

func TestCreateUserExists(t *testing.T) {
	assert := assert.New(t)

	dbName := "test_createdb"

	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	log.Println("Testing if user exists")
	created, err := mongoConConfig.TestProvisioner.IsUserCreated(dbName, "mytestUser")
	assert.Nil(err)
	assert.True(created)
}

func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)

	dbName := "test_createdb"

	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	log.Println("Removing test user")
	err := mongoConConfig.TestProvisioner.DeleteUser(dbName, "mytestUser")
	assert.Nil(err)
}

func TestDeleteTheDatabase(t *testing.T) {
	assert := assert.New(t)

	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	dbName := "test_createdb"
	log.Println("Removing test database")

	err := mongoConConfig.TestProvisioner.DeleteDatabase(dbName)
	assert.Nil(err)
}