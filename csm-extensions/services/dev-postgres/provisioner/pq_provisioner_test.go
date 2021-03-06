package provisioner

import (
	"os"
	"strings"
	"testing"

	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-postgres/config"
	_ "github.com/lib/pq"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("postgres-provisioner-test")

var testPostgresProv = struct {
	postgresProvisioner PostgresProvisionerInterface
	postgresConfig      config.PostgresConfig
}{}

func init() {
	testPostgresProv.postgresConfig = config.PostgresConfig{
		User:     os.Getenv("SERVICE_POSTGRES_USER"),
		Password: os.Getenv("SERVICE_POSTGRES_PASSWORD"),
		Host:     os.Getenv("SERVICE_POSTGRES_HOST"),
		Port:     os.Getenv("SERVICE_POSTGRES_PORT"),
		Dbname:   os.Getenv("SERVICE_POSTGRES_DBNAME"),
		Sslmode:  os.Getenv("SERVICE_POSTGRES_SSLMODE")}

	testPostgresProv.postgresProvisioner = NewPqProvisioner(logger, testPostgresProv.postgresConfig)
}

func TestPqProvisioner(t *testing.T) {
	assert := assert.New(t)

	newDbName := "testcreatedb"

	if !envVarsOk() {
		t.Skip("Skipping test, not all env variables are set:'SERVICE_POSTGRES_USER','SERVICE_POSTGRES_PASSWORD','SERVICE_POSTGRES_HOST','SERVICE_POSTGRES_PORT','SERVICE_POSTGRES_DBNAME','SERVICE_POSTGRES_SSLMODE'")
	}

	// Create database

	err := testPostgresProv.postgresProvisioner.CreateDatabase(newDbName)
	assert.NoError(err)

	// Check database exists

	exist, err := testPostgresProv.postgresProvisioner.DatabaseExists(newDbName)
	assert.NoError(err)
	assert.True(exist)

	newUser := "testuser"

	// Create User

	err = testPostgresProv.postgresProvisioner.CreateUser(newDbName, newUser, "aPassw0rd")
	assert.NoError(err)

	exist, err = testPostgresProv.postgresProvisioner.UserExists(newUser)
	assert.NoError(err)
	assert.True(exist)

	// Delete user

	err = testPostgresProv.postgresProvisioner.DeleteUser(newDbName, newUser)
	assert.NoError(err)

	// Check user was deleted

	exist, err = testPostgresProv.postgresProvisioner.UserExists(newUser)
	assert.NoError(err)
	assert.False(exist)

	// Delete database

	err = testPostgresProv.postgresProvisioner.DeleteDatabase(newDbName)
	assert.NoError(err)

	// Check database was deleted

	exist, err = testPostgresProv.postgresProvisioner.DatabaseExists(newDbName)
	assert.NoError(err)
	assert.False(exist)
}

func TestParametrizeQuery(t *testing.T) {
	_, err := parametrizeQuery("SELECT COUNT(*) FROM pg_roles WHERE rolname = {{.User}}", map[string]string{"Username": "username"})

	if !strings.Contains(err.Error(), "Invalid parameter passed to query") {
		t.Errorf("Error parametrizing query: %v", err)
	}
}

func envVarsOk() bool {
	return testPostgresProv.postgresConfig.User != "" && testPostgresProv.postgresConfig.Password != "" && testPostgresProv.postgresConfig.Host != "" &&
		testPostgresProv.postgresConfig.Port != "" && testPostgresProv.postgresConfig.Dbname != "" && testPostgresProv.postgresConfig.Sslmode != ""
}
