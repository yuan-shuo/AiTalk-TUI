package cmd

import (
	"aitalk/config"
	"aitalk/utils/json"
)

func Run(c *config.Config) error {
	req := json.NewReqStruct(c)
	err := loopScan(req, c)
	if err != nil {
		return err
	}
	return nil
}
