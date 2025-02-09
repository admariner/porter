package namespace

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/porter-dev/porter/api/server/authz"
	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/server/shared/requestutils"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/kubernetes"
	"github.com/porter-dev/porter/internal/models"
)

type GetIngressHandler struct {
	handlers.PorterHandlerReadWriter
	authz.KubernetesAgentGetter
}

func NewGetIngressHandler(
	config *config.Config,
	resultWriter shared.ResultWriter,
) *GetIngressHandler {
	return &GetIngressHandler{
		PorterHandlerReadWriter: handlers.NewDefaultPorterHandler(config, nil, resultWriter),
		KubernetesAgentGetter:   authz.NewOutOfClusterAgentGetter(config),
	}
}

func (c *GetIngressHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cluster, _ := r.Context().Value(types.ClusterScope).(*models.Cluster)
	agent, err := c.GetAgent(r, cluster, "")
	name, _ := requestutils.GetURLParamString(r, types.URLParamIngressName)
	namespace, _ := requestutils.GetURLParamString(r, types.URLParamNamespace)

	if err != nil {
		c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	ingress, err := agent.GetIngress(namespace, name)

	if targetErr := kubernetes.IsNotFoundError; errors.Is(err, targetErr) {
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
			fmt.Errorf("ingress %s/%s was not found", namespace, name),
			http.StatusNotFound,
		))

		return
	} else if err != nil {
		c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	c.WriteResult(w, r, ingress)
}
