package nginx

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUpstreamWithDefaultServer(t *testing.T) {
	assert := assert.New(t)
	u := NewUpstreamWithDefaultServer("default")
	assert.Equal("default", u.Name)
	assert.Contains(u.UpstreamServers, UpstreamServer{Address: "127.0.0.1", Port: "8181"})
}

func TestShellOut(t *testing.T) {
	t.Run("is executing the command", func(t *testing.T) {
		assert := assert.New(t)
		err := shellOut("echo 'test'")
		assert.NoError(err)
	})

	t.Run("is returning an error if the command does not exist", func(t *testing.T) {
		assert := assert.New(t)
		err := shellOut("a-non-existing-command")
		assert.Error(err)
	})

	t.Run("is returning an error on unclean exit", func(t *testing.T) {
		assert := assert.New(t)
		err := shellOut("exit 2")
		assert.Error(err)
	})
}

func TestNginxController(t *testing.T) {
	var n *NginxController
	confPath := "/tmp/nginx-controller"

	t.Run("NewNginxController", func(t *testing.T) {
		assert := assert.New(t)
		nginx, err := NewNginxController(confPath, false, false)
		n = nginx
		assert.NoError(err)
	})

	t.Run("AddOrUpdateCertAndKey", func(t *testing.T) {
		assert := assert.New(t)
		expectedFileName := "/tmp/nginx-controller/ssl/my-secret.pem"
		fileName := n.AddOrUpdateCertAndKey("my-secret", "cert", "key")
		assert.Equal(expectedFileName, fileName)
		_, err := os.Stat(expectedFileName)
		assert.NoError(err)
	})

	t.Run("AddOrUpdateDHParam", func(t *testing.T) {
		assert := assert.New(t)
		expectedFileName := "/tmp/nginx-controller/ssl/dhparam.pem"
		fileName, err := n.AddOrUpdateDHParam("dhparam")
		if assert.NoError(err) {
			assert.Equal(expectedFileName, fileName)
			_, err := os.Stat(expectedFileName)
			assert.NoError(err)
		}
	})

	t.Run("AddOrUpdateConfig", func(t *testing.T) {
		assert := assert.New(t)
		expectedFileName := "/tmp/nginx-controller/conf.d/test_test.conf"
		n.AddOrUpdateConfig("test_test", Server{})
		_, err := os.Stat(expectedFileName)
		assert.NoError(err)
	})

	t.Run("AddOrUpdateConfig emptyHost", func(t *testing.T) {
		assert := assert.New(t)
		expectedFileName := "/tmp/nginx-controller/conf.d/default.conf"
		n.AddOrUpdateConfig(emptyHost, Server{})
		_, err := os.Stat(expectedFileName)
		assert.NoError(err)
	})

	t.Run("DeleteConfig", func(t *testing.T) {
		assert := assert.New(t)
		expectedFileName := "/tmp/nginx-controller/conf.d/test_test.conf"
		n.DeleteConfig("test_test")
		_, err := os.Stat(expectedFileName)
		if assert.Error(err) {
			assert.True(os.IsNotExist(err))
		}
	})

	// cleanup
	os.RemoveAll(confPath)
}
