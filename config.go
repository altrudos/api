package charityhonor

import "github.com/charityhonor/ch-api/pkg/justgiving"

var defaultJG *justgiving.JustGiving

func MustGetDefaultJustGiving() *justgiving.JustGiving {
	if defaultJG == nil {
		mode := GetEnv("JG_MODE", "")
		appId := GetEnv("JG_APPID", "")

		if mode == "" {
			panic("JG_MODE env var not set")
		}

		if appId == "" {
			panic("JG_APPID env var not set")
		}

		defaultJG = &justgiving.JustGiving{
			AppId: appId,
			Mode:  justgiving.Mode(mode),
		}
	}

	return defaultJG
}
