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
	// return random 8 tips from blocks confirmed by the last milestone (17011900)
	return &tipsResponse{Tips: []string{
		"0xae75d93d1ea3e0a8bee2245cd33f394f39f1cb400b0e856ba9f4e29ed8c67b42",
		"0x860c510dc4dc21c612ae617e560c3b997b9fa111713c4952450a69f4b0fe2dde",
		"0x6bb8f6d4b4bc69b385a99844581f3f0d8bcb74d629bcda3b88fca390ddb517ce",
		"0x4b9a57dca736697df5dfc5e4eaf67f8607f7715f8713036b815e0b67ebe8fb07",
		"0x2740938a995a7c87446c66f1b013cb10ee269ce7c8719eafd6ba4e66a686b0ff",
		"0x53be53cf1614d7e3fc05a479134c8925ddaa49031aedf6b8b9d999bccf52d8c6",
		"0xba7a2ccb27f6fe91b72f480ea5730a622f8f85577081d22cb2af78cd52ad4b03",
		"0x77f062528b22f7e29d6a157792534c1ceef2765e4e47d6e9a2775fe71fbac6d9",
	}}, nil
}
