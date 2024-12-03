package shell

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/kloudlite/kl/domain/apiclient"
	"github.com/kloudlite/kl/domain/fileclient"
	fn "github.com/kloudlite/kl/pkg/functions"
)

func ParseEnvVarsAndMounts(cfg *fileclient.KLFileType) (envVars []string, mounts map[string][]byte, err error) {
	// TODO: validate whether Team with name: cfg.TeamName exists
	cookie, err := fileclient.GetCookieString(fn.MakeOption("teamName", cfg.TeamName))
	if err != nil {
		return nil, nil, err
	}

	currMreses := cfg.EnvVars.GetMreses()
	currSecs := cfg.EnvVars.GetSecrets()
	currConfs := cfg.EnvVars.GetConfigs()

	currMounts := cfg.Mounts.GetMounts()

	m := map[string]any{
		"envName": cfg.DefaultEnv,
		"configQueries": func() []any {
			var queries []any
			for _, v := range currConfs {
				for _, vv := range v.Env {
					queries = append(queries, map[string]any{
						"configName": v.Name,
						"key":        vv.RefKey,
					})
				}
			}

			for _, fe := range currMounts {
				if fe.Type == fileclient.ConfigType {
					queries = append(queries, map[string]any{
						"configName": fe.Name,
						"key":        fe.Key,
					})
				}
			}

			return queries
		}(),

		"mresQueries": func() []any {
			var queries []any
			for _, rt := range currMreses {
				for _, v := range rt.Env {
					queries = append(queries, map[string]any{
						"secretName": rt.Name,
						"key":        v.RefKey,
					})
				}
			}

			return queries
		}(),

		"secretQueries": func() []any {
			var queries []any
			for _, v := range currSecs {
				for _, vv := range v.Env {
					queries = append(queries, map[string]any{
						"secretName": v.Name,
						"key":        vv.RefKey,
					})
				}
			}

			for _, fe := range currMounts {
				if fe.Type == fileclient.SecretType {
					queries = append(queries, map[string]any{
						"secretName": fe.Name,
						"key":        fe.Key,
					})
				}
			}
			return queries
		}(),
	}

	respData, err := klFetch("cli_getConfigSecretMap", m, &cookie)
	if err != nil {
		return nil, nil, fn.NewE(err)
	}

	var resp apiclient.EnvRsp
	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, nil, err
	}

	rm := make(map[string]string)

	for _, v := range resp.Configs {
		rm[fmt.Sprintf("config/%s/%s", v.ConfigName, v.Key)] = v.Value
	}

	for _, v := range resp.Secrets {
		rm[fmt.Sprintf("secret/%s/%s", v.SecretName, v.Key)] = v.Value
	}

	for _, v := range resp.Mreses {
		rm[fmt.Sprintf("mres/%s/%s", v.SecretName, v.Key)] = v.Value
	}

	result := make([]string, 0, len(cfg.EnvVars))

	for _, v := range cfg.EnvVars {
		switch {
		case v.ConfigRef != nil:
			result = append(result, fmt.Sprintf("%s=%s", v.Key, rm["config/"+*v.ConfigRef]))
		case v.SecretRef != nil:
			result = append(result, fmt.Sprintf("%s=%s", v.Key, rm["secret/"+*v.SecretRef]))
		case v.MresRef != nil:
			result = append(result, fmt.Sprintf("%s=%s", v.Key, rm["mres/"+*v.MresRef]))
		}
	}

	mounts = make(map[string][]byte, len(cfg.Mounts))

	for _, v := range cfg.Mounts {
		switch {
		case v.ConfigRef != nil:
			mounts[v.Path] = []byte(rm["config/"+*v.ConfigRef])
		case v.SecretRef != nil:
			mounts[v.Path] = []byte(rm["secret/"+*v.SecretRef])
		}
	}

	slog.Debug("API (cli_getConfigSecretMap), got", "resp", respData)
	return result, mounts, nil
}
