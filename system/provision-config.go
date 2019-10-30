package system

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	provision "github.com/cortezaproject/corteza-server/provision/system"
	"github.com/cortezaproject/corteza-server/system/importer"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

func provisionConfig(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
	c.Log.Debug("running configuration provision")
	c.InitServices(ctx, c)

	var provisioned bool

	// Make sure we have all full access for provisioning
	ctx = auth.SetSuperUserContext(ctx)

	// if system is already provisioned, we do partial provisioning:
	// missing settings only
	if provisioned, err = isProvisioned(ctx); err != nil {
		return err
	} else if provisioned {
		c.Log.Debug("configuration already provisioned")
	}

	readers, err := impAux.ReadStatic(provision.Asset)
	if err != nil {
		return err
	}

	if provisioned {
		return partialImportSettings(ctx, service.DefaultSettings, readers...)
	}

	return errors.Wrap(
		importer.Import(ctx, readers...),
		"could not provision configuration for system service",
	)
}

// Provision ONLY when there are no rules for role admins
func isProvisioned(ctx context.Context) (bool, error) {
	return len(service.DefaultPermissions.FindRulesByRoleID(permissions.AdminsRoleID)) > 0, nil
}

func makeDefaultApplications(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	db, err := factory.Database.Get(system)
	if err != nil {
		return err
	}

	repo := repository.Application(ctx, db)

	aa, _, err := repo.Find(types.ApplicationFilter{})
	if err != nil {
		return err
	}

	// List of apps to create.
	//
	// We use Unify.Url field for matching,
	// so make sure it's always present!
	defApps := types.ApplicationSet{
		&types.Application{
			Name:    "CRM",
			Enabled: true,
			Unify: &types.ApplicationUnify{
				Name:   "CRM",
				Listed: true,
				Icon:   "/applications/crust_favicon.png",
				Logo:   "/applications/crust.jpg",
				Url:    "/compose/ns/crm/pages",
			},
		},
	}

	return defApps.Walk(func(defApp *types.Application) error {
		for _, a := range aa {
			if a.Unify != nil && a.Unify.Url == defApp.Unify.Url {
				// App already added.
				return nil
			}
		}

		defApp, err = repo.Create(defApp)
		c.Log.Info(
			"creating default application",
			zap.String("name", defApp.Name),
			zap.Uint64("name", defApp.ID),
			zap.Error(err),
		)
		return err

		return nil
	})
}

// Partial import of settings from provision files
func partialImportSettings(ctx context.Context, ss service.SettingsService, ff ...io.Reader) (err error) {
	var (
		// decoded content from YAML files
		aux interface{}

		si = settings.NewImporter()

		// importer w/o permissions & roles
		// we need only settings
		imp = importer.NewImporter(nil, si, nil)

		// current value
		current settings.ValueSet

		// unexisting values
		unex settings.ValueSet
	)

	for _, f := range ff {
		if err = yaml.NewDecoder(f).Decode(&aux); err != nil {
			return
		}

		err = imp.Cast(aux)
		if err != nil {
			return
		}
	}

	ss = ss.With(ctx)

	// Get all "current" settings storage
	current, err = ss.FindByPrefix("")
	if err != nil {
		return
	}

	// Compare current settings with imported, get all that do not exist yet
	if unex = si.GetValues(); len(unex) > 0 {
		// Store non existing
		err = ss.BulkSet(current.New(unex))
		if err != nil {
			return
		}
	}

	return nil
}