/*-------------------------------------------------------------------------
 *
 * wizard_test.go
 *    Tests for interactive wizard (integration requires stdin).
 *
 *-------------------------------------------------------------------------
 */

package wizard

import (
	"testing"

	"github.com/neurondb/NeuronAgent/cli/pkg/client"
)

func TestRunWizard_RequiresInteractive(t *testing.T) {
	t.Skip("RunWizard requires interactive stdin")
}

func TestWizardPackage(t *testing.T) {
	_ = client.NewClient("http://localhost:8080", "")
	/* Package loads and client dependency is valid. */
}
