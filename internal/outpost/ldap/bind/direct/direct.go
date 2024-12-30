package direct

import (
	"context"
	"errors"
	"strings"

	goldap "github.com/go-ldap/ldap/v3"
	"go.uber.org/zap"
	"goauthentik.io/internal/config"
	"goauthentik.io/internal/outpost/flow"
	"goauthentik.io/internal/outpost/ldap/server"
	"goauthentik.io/internal/outpost/ldap/utils"
)

const ContextUserKey = "ak_user"

type DirectBinder struct {
	si  server.LDAPServerInstance
	log *zap.Logger
}

func NewDirectBinder(si server.LDAPServerInstance) *DirectBinder {
	db := &DirectBinder{
		si:  si,
		log: config.Get().Logger().Named("authentik.outpost.ldap.binder.direct"),
	}
	db.log.Info("initialised direct binder")
	return db
}

func (db *DirectBinder) GetUsername(dn string) (string, error) {
	if !utils.HasSuffixNoCase(dn, db.si.GetBaseDN()) {
		return "", errors.New("invalid base DN")
	}
	dns, err := goldap.ParseDN(dn)
	if err != nil {
		return "", err
	}
	for _, part := range dns.RDNs {
		for _, attribute := range part.Attributes {
			if strings.ToLower(attribute.Type) == "cn" {
				return attribute.Value, nil
			}
		}
	}
	return "", errors.New("failed to find cn")
}

func (db *DirectBinder) TimerFlowCacheExpiry(ctx context.Context) {
	fe := flow.NewFlowExecutor(ctx, db.si.GetAuthenticationFlowSlug(), db.si.GetAPIClient().GetConfig(), []zap.Field{})
	fe.Params.Add("goauthentik.io/outpost/ldap", "true")
	fe.Params.Add("goauthentik.io/outpost/ldap-warmup", "true")

	err := fe.WarmUp()
	if err != nil {
		db.log.Warn("failed to warm up flow cache", zap.Error(err))
	}
}
