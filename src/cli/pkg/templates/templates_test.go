/*-------------------------------------------------------------------------
 *
 * templates_test.go
 *    Tests for template listing and loading.
 *
 *-------------------------------------------------------------------------
 */

package templates

import (
	"testing"
)

func TestListTemplates(t *testing.T) {
	list, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates: %v", err)
	}
	/* list may be nil when templates dir does not exist */
	_ = list
}

func TestSearchTemplates(t *testing.T) {
	list, err := SearchTemplates("test")
	if err != nil {
		t.Fatalf("SearchTemplates: %v", err)
	}
	/* list may be nil or empty when no templates exist */
	_ = list
}
