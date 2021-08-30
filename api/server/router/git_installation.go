package router

import (
	"github.com/go-chi/chi"
	"github.com/porter-dev/porter/api/server/handlers/gitinstallation"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/types"
)

func NewGitInstallationScopedRegisterer(children ...*Registerer) *Registerer {
	return &Registerer{
		GetRoutes: GetGitInstallationScopedRoutes,
		Children:  children,
	}
}

func GetGitInstallationScopedRoutes(
	r chi.Router,
	config *config.Config,
	basePath *types.Path,
	factory shared.APIEndpointFactory,
	children ...*Registerer,
) []*Route {
	routes, projPath := getGitInstallationRoutes(r, config, basePath, factory)

	if len(children) > 0 {
		r.Route(projPath.RelativePath, func(r chi.Router) {
			for _, child := range children {
				childRoutes := child.GetRoutes(r, config, basePath, factory, child.Children...)

				routes = append(routes, childRoutes...)
			}
		})
	}

	return routes
}

func getGitInstallationRoutes(
	r chi.Router,
	config *config.Config,
	basePath *types.Path,
	factory shared.APIEndpointFactory,
) ([]*Route, *types.Path) {
	relPath := "/gitrepos/{git_installation_id}"

	newPath := &types.Path{
		Parent:       basePath,
		RelativePath: relPath,
	}

	routes := make([]*Route, 0)

	// GET /api/projects/{project_id}/gitrepos/{git_installation_id} -> gitinstallation.NewGitInstallationGetHandler
	getEndpoint := factory.NewAPIEndpoint(
		&types.APIRequestMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &types.Path{
				Parent:       basePath,
				RelativePath: relPath,
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.ProjectScope,
				types.GitInstallationScope,
			},
		},
	)

	getHandler := gitinstallation.NewGitInstallationGetHandler(
		config,
		factory.GetResultWriter(),
	)

	routes = append(routes, &Route{
		Endpoint: getEndpoint,
		Handler:  getHandler,
		Router:   r,
	})

	// GET /api/projects/{project_id}/gitrepos/{git_installation_id}/repos ->
	// gitinstallation.GithubListReposHandler
	listReposEndpoint := factory.NewAPIEndpoint(
		&types.APIRequestMetadata{
			Verb:   types.APIVerbList,
			Method: types.HTTPVerbGet,
			Path: &types.Path{
				Parent:       basePath,
				RelativePath: relPath + "/repos",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.ProjectScope,
				types.GitInstallationScope,
			},
		},
	)

	listReposHandler := gitinstallation.NewGithubListReposHandler(
		config,
		factory.GetResultWriter(),
	)

	routes = append(routes, &Route{
		Endpoint: listReposEndpoint,
		Handler:  listReposHandler,
		Router:   r,
	})

	return routes, newPath
}
