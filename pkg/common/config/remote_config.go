package config

import (
	"HugeSpaceship/pkg/common/config/model/remote"
	"HugeSpaceship/pkg/common/db"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	lbpApiSection = "lbp_api"
)

func GetLBPAPIConfig() remote.LBPAPIConfig {
	store, err := db.GetConnection().GetConfigSection(lbpApiSection)
	if err.Error() == "no rows in result set" {
		err := SaveLBPAPIConfig(remote.LBPAPIConfig{
			PrimaryDigest: "",
		})
		if err != nil {

		}
	}
	for k, v := range store {
		fmt.Println(k, v)
	}
	return remote.LBPAPIConfig{PrimaryDigest: "", AlternateDigest: ""}
}

func SaveLBPAPIConfig(config remote.LBPAPIConfig) error {
	store := pgtype.Hstore{}
	err := store.Scan(&config)
	if err != nil {
		return err
	}

	err = db.GetConnection().StoreConfig(lbpApiSection, store)
	return err
}
