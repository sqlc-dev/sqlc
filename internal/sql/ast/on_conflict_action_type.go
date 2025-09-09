package ast

const (
	OnConflictAction_ON_CONFLICT_ACTION_UNDEFINED OnConflictAction = 0
	OnConflictAction_ONCONFLICT_NONE              OnConflictAction = 1
	OnConflictAction_ONCONFLICT_NOTHING           OnConflictAction = 2
	OnConflictAction_ONCONFLICT_UPDATE            OnConflictAction = 3

	// YQL-specific
	OnConflictAction_INSERT_OR_ABORT   OnConflictAction = 4
	OnConflictAction_INSERT_OR_REVERT  OnConflictAction = 5
	OnConflictAction_INSERT_OR_IGNORE  OnConflictAction = 6
	OnConflictAction_UPSERT            OnConflictAction = 7
	OnConflictAction_REPLACE           OnConflictAction = 8
)
