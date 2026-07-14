package service

import "time"

// MaxIdleDuration は放置報酬の上限時間
// これ無いと、放置報酬が無限に増え続ける
const MaxIdleDuration = 12 * time.Hour

// CalcIdleGain は最終claim時刻と生産レートから　現時点で受け取れる放置コイン量を計算する
func CalcIdleGain(lastClaimedAt time.Time, ratePerSecond uint32, now time.Time) int64 {
	elapsed := now.Sub(lastClaimedAt)
	if elapsed <= 0 {
		return 0
	}
	if elapsed > MaxIdleDuration {
		elapsed = MaxIdleDuration
	}
	seconds := int64(elapsed.Seconds())
	return seconds * int64(ratePerSecond)
}
