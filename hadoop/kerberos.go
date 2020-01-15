package hadoop

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/xloger/g"
	krb "gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/credentials"
)

func getKrbosCachePath() (string, error) {
	var v, p string
	v = os.Getenv("KRB5CCNAME")
	// KRB5CCNAME=FILE:/opt/userdata/krb5cache/1006/krb5cc_1006

	if v == "" {
		u, err := user.Current()
		if err != nil {
			return "", err
		}

		p = fmt.Sprintf("/tmp/krb5cc_%s", u.Uid)
		return p, nil
	}

	if !strings.HasPrefix(v, "FILE:") {
		errMsg := fmt.Sprintf("unusable cache_path: %s", v)
		return "", errors.New(errMsg)
	}

	pn := strings.SplitN(v, ":", 2)
	if len(pn) != 2 {
		errMsg := fmt.Sprintf("unsplitN %s", v)
		return "", errors.New(errMsg)
	}

	return pn[1], nil
}

func NewKerberosClient() (*krb.Client, error) {
	krbConf, err := g.GetKrb5Conf()
	if err != nil {
		return nil, err
	}

	cfg, err := config.Load(krbConf)
	if err != nil {
		return nil, err
	}

	cachePath, err := getKrbosCachePath()
	if err != nil {
		return nil, err
	}

	ccache, err := credentials.LoadCCache(cachePath)
	if err != nil {
		return nil, err
	}

	client, err := krb.NewClientFromCCache(ccache, cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}
