package teamvault

import (
	"github.com/foomo/htpasswd"
	"github.com/golang/glog"
)

type Htpasswd struct {
	Connector Connector
}

func (c *Htpasswd) Generate(key Key) ([]byte, error) {
	pass, err := c.Connector.Password(key)
	if err != nil {
		glog.V(2).Infof("get password from teamvault for key %v failed: %v", key, err)
		return nil, err
	}
	user, err := c.Connector.User(key)
	if err != nil {
		glog.V(2).Infof("get user from teamvault for key %v failed: %v", key, err)
		return nil, err
	}
	pws := make(htpasswd.HashedPasswords)
	err = pws.SetPassword(string(user), string(pass), htpasswd.HashBCrypt)
	if err != nil {
		glog.V(2).Infof("set password failed for key %v failed: %v", key, err)
		return nil, err
	}
	return pws.Bytes(), nil
}
