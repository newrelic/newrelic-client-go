package synthetics

import (
	"context"
)

func (a *Synthetics) SyntheticsEditOneTimeMonitorDowntimeWithContext(
	ctx context.Context,
	gUID EntityGUID,
	monitorGUIDs []EntityGUID,
	name string,
	once SyntheticsMonitorDowntimeOnceConfig,
) (*SyntheticsMonitorDowntimeMutationResult, error) {

	resp := SyntheticsEditMonitorDowntimeQueryResponse{}
	vars := map[string]interface{}{
		"guid":         gUID,
		"monitorGuids": monitorGUIDs,
		"name":         name,
		"once":         once,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, SyntheticsEditOneTimeMonitorDowntimeMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.SyntheticsMonitorDowntimeMutationResult, nil
}

func (a *Synthetics) SyntheticsEditDailyMonitorDowntimeWithContext(
	ctx context.Context,
	gUID EntityGUID,
	monitorGUIDs []EntityGUID,
	name string,
	daily SyntheticsMonitorDowntimeDailyConfig,
) (*SyntheticsMonitorDowntimeMutationResult, error) {

	resp := SyntheticsEditMonitorDowntimeQueryResponse{}
	vars := map[string]interface{}{
		"guid":         gUID,
		"monitorGuids": monitorGUIDs,
		"name":         name,
		"daily":        daily,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, SyntheticsEditDailyMonitorDowntimeMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.SyntheticsMonitorDowntimeMutationResult, nil
}

func (a *Synthetics) SyntheticsEditWeeklyMonitorDowntimeWithContext(
	ctx context.Context,
	gUID EntityGUID,
	monitorGUIDs []EntityGUID,
	name string,
	weekly SyntheticsMonitorDowntimeWeeklyConfig,
) (*SyntheticsMonitorDowntimeMutationResult, error) {

	resp := SyntheticsEditMonitorDowntimeQueryResponse{}
	vars := map[string]interface{}{
		"guid":         gUID,
		"monitorGuids": monitorGUIDs,
		"name":         name,
		"weekly":       weekly,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, SyntheticsEditWeeklyMonitorDowntimeMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.SyntheticsMonitorDowntimeMutationResult, nil
}

func (a *Synthetics) SyntheticsEditMonthlyMonitorDowntimeWithContext(
	ctx context.Context,
	gUID EntityGUID,
	monitorGUIDs []EntityGUID,
	name string,
	monthly SyntheticsMonitorDowntimeMonthlyConfig,
) (*SyntheticsMonitorDowntimeMutationResult, error) {

	resp := SyntheticsEditMonitorDowntimeQueryResponse{}
	vars := map[string]interface{}{
		"guid":         gUID,
		"monitorGuids": monitorGUIDs,
		"name":         name,
		"monthly":      monthly,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, SyntheticsEditMonthlyMonitorDowntimeMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.SyntheticsMonitorDowntimeMutationResult, nil
}

const SyntheticsEditOneTimeMonitorDowntimeMutation = `mutation(
	$guid: EntityGuid!,
	$monitorGuids: [EntityGuid],
	$name: String,
	$once: SyntheticsMonitorDowntimeOnceConfig,
) { syntheticsEditMonitorDowntime(
	guid: $guid,
	monitorGuids: $monitorGuids,
	name: $name,
	once: $once,
) {
	accountId
	endTime
	guid
	monitorGuids
	name
	startTime
	timezone
} }`

const SyntheticsEditDailyMonitorDowntimeMutation = `mutation(
	$guid: EntityGuid!,
	$monitorGuids: [EntityGuid],
	$name: String,
	$daily: SyntheticsMonitorDowntimeDailyConfig,
) { syntheticsEditMonitorDowntime(
	guid: $guid,
	monitorGuids: $monitorGuids,
	name: $name,
	daily: $daily,
) {
	accountId
	endRepeat {
		onDate
		onRepeat
	}
	endTime
	guid
	monitorGuids
	name
	startTime
	timezone
} }`

const SyntheticsEditWeeklyMonitorDowntimeMutation = `mutation(
	$guid: EntityGuid!,
	$monitorGuids: [EntityGuid],
	$name: String,
	$weekly: SyntheticsMonitorDowntimeWeeklyConfig,
) { syntheticsEditMonitorDowntime(
	guid: $guid,
	monitorGuids: $monitorGuids,
	name: $name,
	weekly: $weekly,
) {
	accountId
	endRepeat {
		onDate
		onRepeat
	}
	endTime
	guid
	maintenanceDays
	monitorGuids
	name
	startTime
	timezone
} }`

const SyntheticsEditMonthlyMonitorDowntimeMutation = `mutation(
	$guid: EntityGuid!,
	$monitorGuids: [EntityGuid],
	$name: String,
	$monthly: SyntheticsMonitorDowntimeMonthlyConfig,
) { syntheticsEditMonitorDowntime(
	guid: $guid,
	monitorGuids: $monitorGuids,
	name: $name,
	monthly: $monthly,
) {
	accountId
	endRepeat {
		onDate
		onRepeat
	}
	endTime
	frequency {
		daysOfMonth
		daysOfWeek {
			ordinalDayOfMonth
			weekDay
		}
	}
	guid
	monitorGuids
	name
	startTime
	timezone
} }`
