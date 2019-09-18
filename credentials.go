package docker_credentials

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Credentials struct {
	Auths map[string]struct {
		Auth string
	}
	CredsStore  string
	CredHelpers map[string]string
}

func Of(registry string) (string, string, error) {
	creds, err := Get()
	if err != nil {
		return "", "", err
	}
	return creds.Of(registry)
}

func Get() (Credentials, error) {
	var home string
	if runtime.GOOS == `windows` {
		home = `USERPROFILE`
	} else {
		home = `HOME`
	}
	return File(filepath.Join(os.Getenv(home), `.docker`, `config.json`))
}

func File(path string) (Credentials, error) {
	var creds Credentials
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return creds, nil
	} else if err != nil {
		return creds, err
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return creds, err
	}
	return New(content)
}

func New(content []byte) (Credentials, error) {
	var creds Credentials
	if err := json.Unmarshal(content, &creds); err != nil {
		return creds, err
	}
	return creds, nil
}

func (c Credentials) Of(registry string) (string, string, error) {
	auth, ok := c.AuthOf(registry)
	if !ok {
		return "", "", nil
	}
	if auth != "" {
		return decodeAuth(auth)
	}
	if store := c.StoreOf(registry); store != "" {
		return getAuthFromStore(registry, store)
	}
	return "", "", nil
}

func (c Credentials) AuthOf(registry string) (string, bool) {
	if len(c.Auths) == 0 || registry == "" {
		return "", false
	}
	if v, ok := c.Auths[registry]; ok {
		return v.Auth, true
	}
	for k, v := range c.Auths {
		if u, err := url.Parse(k); err == nil && u.Hostname() == registry {
			return v.Auth, true
		}
	}
	return "", false
}

func (c Credentials) StoreOf(registry string) string {
	if c.CredHelpers != nil && registry != "" {
		if helper, ok := c.CredHelpers[registry]; ok {
			return helper
		}
	}
	return c.CredsStore
}

func decodeAuth(encoded string) (string, string, error) {
	plain, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}
	arr := strings.SplitN(string(plain), ":", 2)
	if len(arr) != 2 {
		return "", "", errors.New("Invalid auth configuration file")
	}
	return arr[0], strings.Trim(arr[1], "\x00"), nil
}

func getAuthFromStore(registry, store string) (string, string, error) {
	cmd := exec.Command("docker-credential-"+store, "get")
	cmd.Stdin = strings.NewReader(registry)
	var output struct {
		Username string
		Secret   string
	}
	if b, err := cmd.Output(); err != nil {
		return "", "", err
	} else if err := json.Unmarshal(b, &output); err != nil {
		return "", "", err
	}
	return output.Username, output.Secret, nil
}
