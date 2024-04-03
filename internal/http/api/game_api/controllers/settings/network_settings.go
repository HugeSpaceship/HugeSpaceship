package settings

import (
	"net/http"
)

const networkSettingsTemplate = `
ProbabilityOfPacketDelay        0.0
MinPacketDelayFrames    0
MaxPacketDelayFrames    3

ProbabilityOfPacketDrop         0.0

EnableFakeConditionsForLoopback true

NumberOfFramesPredictionAllowedForNonLocalPlayer 1000

EnablePrediction                        true
MinPredictedFrames                      0
MaxPredictedFrames                      10

AllowGameRendCameraSplit        true
FramesBeforeAgressiveCatchup    30


PredictionPadSides                      200
PredictionPadTop                        200
PredictionPadBottom                     200

ShowErrorNumbers true
AllowModeratedLevels true
AllowModeratedPoppetItems true

ShowLevelBoos true

TIMEOUT_WAIT_FOR_JOIN_RESPONSE_FROM_PREV_PARTY_HOST             50.0
TIMEOUT_WAIT_FOR_CHANGE_LEVEL_PARTY_HOST                                30.0
TIMEOUT_WAIT_FOR_CHANGE_LEVEL_PARTY_MEMBER                              45.0
TIMEOUT_WAIT_FOR_REQUEST_JOIN_FRIEND                                    15.0
TIMEOUT_WAIT_FOR_CONNECTION_FROM_HOST                                   30.0
TIMEOUT_WAIT_FOR_ROOM_ID_TO_JOIN                                                60.0
TIMEOUT_WAIT_FOR_GET_NUM_PLAYERS_ONLINE                                 60.0
TIMEOUT_WAIT_FOR_SIGNALLING_CONNECTIONS                                 120.0
TIMEOUT_WAIT_FOR_PARTY_DATA                                                             60.0
TIME_TO_WAIT_FOR_LEAVE_MESSAGE_TO_COME_BACK                             20.0
TIME_TO_WAIT_FOR_FOLLOWING_REQUESTS_TO_ARRIVE                   30.0
TIMEOUT_WAIT_FOR_FINISHED_MIGRATING_HOST                                30.0
TIMEOUT_WAIT_FOR_PARTY_LEADER_FINISH_JOINING                    45.0
TIMEOUT_WAIT_FOR_QUICKPLAY_LEVEL                                                60.0
TIMEOUT_WAIT_FOR_PLAYERS_TO_JOIN                                                30.0
TIMEOUT_WAIT_FOR_DIVE_IN_PLAYERS                                                120.0
TIMEOUT_WAIT_FOR_FIND_BEST_ROOM                                                 30.0
TIMEOUT_DIVE_IN_TOTAL                                                                   120.0

TIMEOUT_WAIT_FOR_SOCKET_CONNECTION                                              120.0
TIMEOUT_WAIT_FOR_REQUEST_RESOURCE_MESSAGE                               120.0
TIMEOUT_WAIT_FOR_LOCAL_CLIENT_TO_GET_RESOURCE_LIST              120.0
TIMEOUT_WAIT_FOR_CLIENT_TO_LOAD_RESOURCES                               120.0
TIMEOUT_WAIT_FOR_LOCAL_CLIENT_TO_SAVE_GAME_STATE                30.0
TIMEOUT_WAIT_FOR_ADD_PLAYERS_TO_TAKE                                    30.0
TIMEOUT_WAIT_FOR_UPDATE_FROM_CLIENT                                             90.0

TIMEOUT_WAIT_FOR_HOST_TO_GET_RESOURCE_LIST                              60.0
TIMEOUT_WAIT_FOR_HOST_TO_SAVE_GAME_STATE                                60.0
TIMEOUT_WAIT_FOR_HOST_TO_ADD_US                                                 30.0
TIMEOUT_WAIT_FOR_UPDATE                                                                 60.0

TIMEOUT_WAIT_FOR_REQUEST_JOIN                                                   50.0

TIMEOUT_WAIT_FOR_AUTOJOIN_PRESENCE                                              60.0
TIMEOUT_WAIT_FOR_AUTOJOIN_CONNECTION                                    120.0

SECONDS_BETWEEN_PINS_AWARDED_UPLOADS                                    300.0

EnableKeepAlive         true
AllowVoIPRecordingPlayback      true

CDNHostName "dev.lbp.valtek.local/api/LBP_XML"
TelemetryServer ""

OverheatingThresholdDisallowMidgameJoin 0.95

MaxCatchupFrames 3
MaxLagBeforeShowLoading 23
MinLagBeforeHideLoading 30

LagImprovementInflectionPoint -1.0
FlickerThreshold 2.0

ClosedDemo2014Version           1
ClosedDemo2014Expired           false

EnablePlayedFilter              true
EnableCommunityDecorations      true

GameStateUpdateRate               10.0
GameStateUpdateRateWithConsumers  1.0

DisableDLCPublishCheck          false
AllowOnlineCreate true
EnableDiveIn                    true
EnableHackChecks                false
`

func NetSettingsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(networkSettingsTemplate))
	}
}
