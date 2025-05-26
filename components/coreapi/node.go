package coreapi

import (
	"time"

	"github.com/labstack/echo/v4"

	iotago "github.com/iotaledger/iota.go/v3"
)

//nolint:unparam // even if the error is never used, the structure of all routes should be the same
func info() (*infoResponse, error) {

	var blocksPerSecond, referencedBlocksPerSecond, referencedRate float64
	lastConfirmedMilestoneMetric := deps.Tangle.LastConfirmedMilestoneMetric()
	if lastConfirmedMilestoneMetric != nil {
		blocksPerSecond = lastConfirmedMilestoneMetric.BPS
		referencedBlocksPerSecond = lastConfirmedMilestoneMetric.RBPS
		referencedRate = lastConfirmedMilestoneMetric.ReferencedRate
	}

	syncState := deps.SyncManager.SyncState()

	// latest milestone
	var latestMilestoneIndex = syncState.LatestMilestoneIndex
	//var latestMilestoneTimestamp uint32
	var latestMilestoneIDHex string
	cachedMilestoneLatest := deps.Storage.CachedMilestoneByIndexOrNil(latestMilestoneIndex) // milestone +1
	if cachedMilestoneLatest != nil {
		//latestMilestoneTimestamp = cachedMilestoneLatest.Milestone().TimestampUnix()
		latestMilestoneIDHex = cachedMilestoneLatest.Milestone().MilestoneIDHex()
		cachedMilestoneLatest.Release(true) // milestone -1
	}

	// confirmed milestone index
	var confirmedMilestoneIndex = syncState.ConfirmedMilestoneIndex
	//var confirmedMilestoneTimestamp uint32
	var confirmedMilestoneIDHex string
	cachedMilestoneConfirmed := deps.Storage.CachedMilestoneByIndexOrNil(confirmedMilestoneIndex) // milestone +1
	if cachedMilestoneConfirmed != nil {
		//confirmedMilestoneTimestamp = cachedMilestoneConfirmed.Milestone().TimestampUnix()
		confirmedMilestoneIDHex = cachedMilestoneConfirmed.Milestone().MilestoneIDHex()
		cachedMilestoneConfirmed.Release(true) // milestone -1
	}

	// pruning index
	var pruningIndex iotago.MilestoneIndex
	snapshotInfo := deps.Storage.SnapshotInfo()
	if snapshotInfo != nil {
		pruningIndex = snapshotInfo.PruningIndex()
	}

	return &infoResponse{
		Name:    deps.AppInfo.Name,
		Version: deps.AppInfo.Version,
		Status: nodeStatus{
			IsHealthy: deps.Tangle.IsNodeHealthy(syncState),
			LatestMilestone: milestoneInfoResponse{
				Index:       latestMilestoneIndex,
				Timestamp:   uint32(time.Now().Unix()), //latestMilestoneTimestamp,
				MilestoneID: latestMilestoneIDHex,
			},
			ConfirmedMilestone: milestoneInfoResponse{
				Index:       confirmedMilestoneIndex,
				Timestamp:   uint32(time.Now().Unix()), //confirmedMilestoneTimestamp,
				MilestoneID: confirmedMilestoneIDHex,
			},
			PruningIndex: pruningIndex,
		},
		SupportedProtocolVersions: deps.ProtocolManager.SupportedVersions(),
		ProtocolParameters:        deps.ProtocolManager.Current(),
		PendingProtocolParameters: deps.ProtocolManager.Pending(),
		BaseToken:                 deps.BaseToken,
		Metrics: nodeMetrics{
			BlocksPerSecond:           blocksPerSecond,
			ReferencedBlocksPerSecond: referencedBlocksPerSecond,
			ReferencedRate:            referencedRate,
		},
		Features: features,
	}, nil
}

func tips(c echo.Context) (*tipsResponse, error) {
	return &tipsResponse{Tips: []string{"0x0000000000000000000000000000000000000000000000000000000000000000"}}, nil
}
