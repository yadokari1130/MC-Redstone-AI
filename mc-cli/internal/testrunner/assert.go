package testrunner

import "mc-cli/internal/model"

// checkAssertions は指定されたアサーションリストを取得済みブロックマップと照合します。
// blockMap のキーは "x,y,z" 形式の文字列です。
func checkAssertions(assertions []Assertion, blockMap map[string]model.BlockData) []AssertionFailure {
	var failures []AssertionFailure

	for _, a := range assertions {
		key := blockKey(a.X, a.Y, a.Z)
		actual, found := blockMap[key]

		if !found {
			// 座標にブロックが存在しない（空気ブロックはAPIから除外されるため）
			// blockがairを期待していて、かつstateチェックも不要な場合はOK
			if isAirExpected(a) {
				continue
			}
			failures = append(failures, AssertionFailure{
				X:        a.X,
				Y:        a.Y,
				Z:        a.Z,
				Expected: a,
				Actual:   nil,
				Reason:   "座標にブロックが存在しません（空気または範囲外）",
			})
			continue
		}

		// ブロックIDのチェック（省略された場合はスキップ）
		if a.Block != "" && a.Block != actual.Block {
			failures = append(failures, AssertionFailure{
				X:        a.X,
				Y:        a.Y,
				Z:        a.Z,
				Expected: a,
				Actual:   &actual,
				Reason:   "ブロックIDが一致しません",
			})
			continue
		}

		// stateの部分一致チェック
		if len(a.State) > 0 {
			for expectKey, expectVal := range a.State {
				actualVal, exists := actual.State[expectKey]
				if !exists || actualVal != expectVal {
					failures = append(failures, AssertionFailure{
						X:        a.X,
						Y:        a.Y,
						Z:        a.Z,
						Expected: a,
						Actual:   &actual,
						Reason:   "ブロック状態が一致しません",
					})
					break
				}
			}
		}
	}

	return failures
}

// isAirExpected はアサーションが空気ブロックを期待しているかを判定します。
func isAirExpected(a Assertion) bool {
	if a.Block == "minecraft:air" || a.Block == "air" {
		return true
	}
	// blockが未指定でstateも未指定の場合はチェックなし（常に合格）
	if a.Block == "" && len(a.State) == 0 {
		return true
	}
	return false
}
