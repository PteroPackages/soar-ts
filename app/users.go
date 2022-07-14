package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/pteropackages/soar/logger"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var log = logger.New()

type user struct {
	ID         int    `json:"id"`
	ExternalID string `json:"external_id"`
	UUID       string `json:"uuid"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Language   string `json:"language"`
	RootAdmin  bool   `json:"root_admin"`
	TwoFactor  bool   `json:"2fa"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type attrModel struct {
	O string `json:"object"`
	A *user  `json:"attributes"`
}

type dataModel struct {
	O string      `json:"object"`
	D []attrModel `json:"data"`
}

var getUsersCmd = &cobra.Command{
	Use: "users:get",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		single, query, err := parseUserQuery(cmd)
		if err != nil {
			log.Error("command error:").WithError(err)
			return
		}

		ctx := http.New(cfg, &cfg.Application, log)
		if single {
			req := ctx.Request("GET", "/api/application/users"+query, nil)
			buf, err := ctx.Execute(req)
			if err != nil {
				log.WithError(err)
				return
			}
			if buf == nil {
				return
			}

			var model attrModel
			if err = json.Unmarshal(buf, &model); err != nil {
				log.Error("failed to parse json:").WithError(err)
				return
			}

			var str []byte
			if cfg.Http.ParseBody {
				str, err = json.MarshalIndent(model.A, "", "  ")
			} else {
				str, err = json.MarshalIndent(model, "", "  ")
			}
			if err != nil {
				log.Error("failed to parse response:").WithError(err)
				return
			}

			log.LineB(str)
			return
		}

		req := ctx.Request("GET", "/api/application/users"+query, nil)
		buf, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}
		if buf == nil {
			return
		}

		var model dataModel
		if err = json.Unmarshal(buf, &model); err != nil {
			log.Error("failed to parse json:").WithError(err)
			return
		}

		var str []byte
		if cfg.Http.ParseBody {
			var inner []user
			for _, u := range model.D {
				inner = append(inner, *u.A)
			}
			str, err = json.MarshalIndent(inner, "", "  ")
		} else {
			str, err = json.MarshalIndent(model, "", "  ")
		}
		if err != nil {
			log.Error("failed to parse response:").WithError(err)
			return
		}

		log.LineB(str)
	},
}

func parseUserQuery(cmd *cobra.Command) (bool, string, error) {
	var query strings.Builder
	single := false
	flags := cmd.Flags()

	if id, _ := flags.GetInt("id"); id != 0 {
		single = true
		query.WriteString(fmt.Sprintf("/%d", id))
	}

	if ext, _ := flags.GetString("external"); ext != "" {
		if query.Len() != 0 {
			return false, "", errors.New("id an external flags specified; pick one")
		}

		query.WriteString("/external/" + ext)
	}

	var params []string
	if val, _ := flags.GetString("username"); val != "" {
		params = append(params, "filter[username]="+val)
	}
	if val, _ := flags.GetString("email"); val != "" {
		params = append(params, "filter[email]="+val)
	}
	if val, _ := flags.GetString("uuid"); val != "" {
		params = append(params, "filter[uuid]="+val)
	}

	if len(params) != 0 {
		query.WriteString("?" + params[0])
		for _, p := range params[1:] {
			query.WriteString("&" + p)
		}
	}

	return single, query.String(), nil
}

var createUserCmd = &cobra.Command{
	Use: "users:create",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		src, _ := cmd.Flags().GetString("src")
		if src == "" {
			log.Error("a source file must be provided")
			return
		}

		buf, err := util.SafeReadFile(src)
		if err != nil {
			log.WithError(err)
			return
		}

		var schema struct {
			Username   string `json:"username"`
			Email      string `json:"email"`
			FirstName  string `json:"first_name"`
			LastName   string `json:"last_name"`
			Password   string `json:"password,omitempty"`
			RootAdmin  bool   `json:"root_admin,omitempty"`
			ExternalID string `json:"external_id,omitempty"`
		}
		if err = json.Unmarshal(buf, &schema); err != nil {
			log.Error("failed to parse json:").WithError(err)
			return
		}

		data, _ := json.Marshal(schema)
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("POST", "/api/application/users", &body)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}
		if buf == nil {
			return
		}

		var model attrModel
		if err = json.Unmarshal(res, &model); err != nil {
			log.Error("failed to parse json:").WithError(err)
			return
		}

		var str []byte
		if cfg.Http.ParseBody {
			str, err = json.MarshalIndent(model.A, "", "  ")
		} else {
			str, err = json.MarshalIndent(model, "", "  ")
		}
		if err != nil {
			log.Error("failed to parse response:").WithError(err)
			return
		}

		log.LineB(str)
	},
}
