// +build ee

package billing

import (
	"net/http"

	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/config"

	"github.com/porter-dev/porter/ee/api/server/handlers/billing"
)

var NewBillingGetTokenHandler func(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
	writer shared.ResultWriter,
) http.Handler

var NewBillingWebhookHandler func(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
) http.Handler

var NewBillingAddProjectHandler func(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
) http.Handler

func init() {
	NewBillingGetTokenHandler = billing.NewBillingGetTokenHandler
	NewBillingWebhookHandler = billing.NewBillingWebhookHandler
	NewBillingAddProjectHandler = billing.NewBillingAddProjectHandler
}
