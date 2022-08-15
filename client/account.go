package client

import (
	"bytes"
	"encoding/json"
	"net/url"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/spf13/cobra"
)

var getAccountCmd = &cobra.Command{
	Use:   "account:get",
	Short: "gets account information",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client/account", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		buf, err := http.HandleItemResponse(res, cfg)
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}

var getPermissionsCmd = &cobra.Command{
	Use:     "account:perms",
	Aliases: []string{"account:p"},
	Short:   "gets system permissions",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client/permissions", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}
		if res == nil {
			return
		}

		buf, err := http.HandleItemResponse(res, cfg)
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}

var getTwoFactorCodeCmd = &cobra.Command{
	Use:     "account:2fa:get",
	Aliases: []string{"2fa:get"},
	Short:   "gets account two-factor code",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		log.Warn("BUG: image_url_data '&' is escaped")
		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client/account/two-factor", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		var model struct {
			Data struct {
				ImageURLData string `json:"image_url_data"`
				Secret       string `json:"secret"`
			} `json:"data"`
		}
		if err = json.Unmarshal(res, &model); err != nil {
			log.WithError(err)
			return
		}

		model.Data.ImageURLData, err = url.PathUnescape(model.Data.ImageURLData)
		if err != nil {
			log.Error("failed to parse url:").WithError(err)
			return
		}

		var buf []byte

		if cfg.Http.ParseBody {
			if cfg.Http.ParseIndent {
				buf, err = json.MarshalIndent(model.Data, "", "  ")
			} else {
				buf, err = json.Marshal(model.Data)
			}
		} else {
			if cfg.Http.ParseIndent {
				buf, err = json.MarshalIndent(model, "", "  ")
			} else {
				buf, err = json.Marshal(model)
			}
		}
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}

var enableTwoFactorCmd = &cobra.Command{
	Use:     "account:2fa:enable",
	Aliases: []string{"2fa:enable"},
	Short:   "enables two-factor on the account",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		if len(args) == 0 {
			log.Error("no two-factor code specified")
			return
		} else if len(args) == 1 {
			log.Error("no account password specified")
			return
		} else if len(args) > 2 {
			log.Error("more than 2 arguments specified (expected: token, password)")
			return
		}

		data, _ := json.Marshal(map[string]string{"code": args[0], "password": args[1]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/account/two-factor", &body)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		buf, err := http.HandleDataResponse(res, cfg)
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}

var disableTwoFactorCmd = &cobra.Command{
	Use:     "account:2fa:disable",
	Aliases: []string{"2fa:disable"},
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		if len(args) == 0 {
			log.Error("no password specified")
			return
		} else if len(args) > 1 {
			log.Error("more than 1 argument specified (expected: password)")
			return
		}

		data, _ := json.Marshal(map[string]string{"password": args[0]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("DELETE", "/api/client/account/two-factor", &body)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
			return
		}
	},
}

var getServersCmd = &cobra.Command{
	Use:   "servers:get",
	Short: "gets account servers",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		buf, err := http.HandleDataResponse(res, cfg)
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}