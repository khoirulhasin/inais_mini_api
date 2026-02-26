package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/helpers"
)

// Add the gorm tags to the model definition
func addGormTags(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	var models []string

	for _, model := range b.Models {
		models = append(models, model.Name)
	}

	for _, model := range b.Models {

		for _, field := range model.Fields {
			if field.Name == "id" {
				field.Tag += ` gorm:"column:id;uniqueIndex;primaryKey;autoIcrement"`
			} else if field.Name == "uuid" {
				field.Tag += ` gorm:"column:uuid;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
			} else if field.Name == "createdAt" {
				field.Tag += ` gorm:"column:created_at;type:bigint;autoCreateTime:milli"`
			} else if field.Name == "updatedAt" {
				field.Tag += ` gorm:"column:updated_at;type:bigint;autoUpdateTime:milli"`
			} else if field.Name == "deletedAt" {
				field.Tag += ` gorm:"column:deleted_at;type:bigint;softDelete:milli;default:0"`
			} else if field.Name == "void" {
				field.Tag += ` gorm:"column:void;type:boolean;default:false"`
			} else if field.Name == "owner" {
				field.Tag += ` gorm:"foreignKey:OwnerID;references:ID"`
			} else if field.Name == "imei" {
				field.Tag += ` gorm:"uniqueIndex:idx_` + strings.ToLower(model.Name) + `_imei,WHERE:deleted_at=0;column:imei"`
			} else if field.Name == "code" {
				field.Tag += ` gorm:"uniqueIndex:idx_` + strings.ToLower(model.Name) + `_code,WHERE:deleted_at=0;column:code"`
			} else if field.Name == "username" {
				field.Tag += ` gorm:"uniqueIndex:idx_` + strings.ToLower(model.Name) + `_username,WHERE:deleted_at=0;column:username"`
			} else if field.Name == "email" {
				field.Tag += ` gorm:"uniqueIndex:idx_` + strings.ToLower(model.Name) + `_email,WHERE:deleted_at=0;column:email"`
			} else if field.Name == "name" {
				field.Tag += ` gorm:"index:idx_` + strings.ToLower(model.Name) + `_` + helpers.ToSnakeCase(field.Name) + `;column:name"`
			} else if helpers.Contains(models, strings.ReplaceAll(strings.ReplaceAll(field.Type.String(), "*", ""), helpers.GetModuleRootName()+"/app/models.", "")) == false && field.Name != "menus" {
				field.Tag += ` gorm:"column:` + helpers.ToSnakeCase(field.Name) + `"`
			}
		}
	}

	return b
}

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	// Attaching the mutation function onto modelgen plugin
	p := modelgen.Plugin{
		MutateHook: addGormTags,
	}

	err = api.Generate(cfg,
		api.NoPlugins(),
		api.AddPlugin(&p),
	)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
