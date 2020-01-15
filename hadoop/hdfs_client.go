package hadoop

import (
	"errors"

	"github.com/colinmarc/hdfs"
	"github.com/colinmarc/hdfs/hadoopconf"
	"github.com/xloger/g"
)

func HdfsClient(namenodes []string) (*hdfs.Client, error) {

	dir, err := g.GetHaoodpConfDir()
	if err != nil {
		return nil, err
	}

	cfg, err := hadoopconf.Load(dir)
	if err != nil || cfg == nil {
		return &hdfs.Client{}, err
	}

	options := hdfs.ClientOptionsFromConf(cfg)
	if len(namenodes) == 0 {
		return &hdfs.Client{}, errors.New("No namenodes server")
	}
	options.Addresses = namenodes

	if options.KerberosClient != nil {
		options.KerberosClient, err = NewKerberosClient()
		if err != nil {
			return &hdfs.Client{}, err
		}
	}
	options.KerberosClient.Config.LibDefaults.DNSLookupKDC = false
	options.KerberosClient.Config.LibDefaults.DNSLookupRealm = false

	client, err := hdfs.NewClient(options)

	if err != nil {
		return &hdfs.Client{}, err
	}
	return client, nil
}
