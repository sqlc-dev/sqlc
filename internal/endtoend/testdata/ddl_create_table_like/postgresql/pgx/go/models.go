// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package querytest

import ()

type Change struct {
	Ranked int32
}

type ChangesRanked struct {
	Ranked                 int32
	RankByEffectSize       int32
	RankByAbsPercentChange int32
}
