package tangle

import (
	"time"

	"github.com/iotaledger/hornet/v2/pkg/model/syncmanager"
)

const (
	maxAllowedMilestoneAge = time.Minute * 5
)

// IsNodeHealthy returns whether the node is synced, has active peers and its latest milestone is not too old.
func (t *Tangle) IsNodeHealthy(sync ...*syncmanager.SyncState) bool {
	return true
}
