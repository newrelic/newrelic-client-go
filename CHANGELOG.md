<a name="v2.73.1"></a>
## [v2.73.1] - 2025-11-05
### Bug Fixes
- **dashboards:** add support for `showApplyAction` in variable options ([#1343](https://github.com/newrelic/newrelic-client-go/issues/1343))

<a name="v2.73.0"></a>
## [v2.73.0] - 2025-10-29
### Features
- **alerts:** add support for target entity in NRQL conditions ([#1341](https://github.com/newrelic/newrelic-client-go/issues/1341))
- **azure:** add azure autodiscovery integration ([#1344](https://github.com/newrelic/newrelic-client-go/issues/1344))

<a name="v2.72.0"></a>
## [v2.72.0] - 2025-10-15
### Features
- **alerts:** revert alerts target entity support ([#1337](https://github.com/newrelic/newrelic-client-go/issues/1337)) ([#1340](https://github.com/newrelic/newrelic-client-go/issues/1340))
- **alerts:** revert alerts outlier detection support ([#1338](https://github.com/newrelic/newrelic-client-go/issues/1338)) ([#1339](https://github.com/newrelic/newrelic-client-go/issues/1339))

<a name="v2.71.0"></a>
## [v2.71.0] - 2025-10-14
### Features
- **alerts:** add outlier detection support ([#1338](https://github.com/newrelic/newrelic-client-go/issues/1338))
- **alerts:** add support for target entity in NRQL conditions ([#1337](https://github.com/newrelic/newrelic-client-go/issues/1337))

<a name="v2.70.2"></a>
## [v2.70.2] - 2025-09-26
### Bug Fixes
- **billboard-settings:** add support for billboard settings in dashboards ([#1324](https://github.com/newrelic/newrelic-client-go/issues/1324))

<a name="v2.70.1"></a>
## [v2.70.1] - 2025-09-24
### Bug Fixes
- **cloud:** add support for OCI Logs integrations in API and types `CloudOciLogsIntegration` ([#1334](https://github.com/newrelic/newrelic-client-go/issues/1334))

<a name="v2.70.0"></a>
## [v2.70.0] - 2025-09-24
### Features
- **cloud:** add OCI provider functions for linkAccount and queries ([#1333](https://github.com/newrelic/newrelic-client-go/issues/1333))

<a name="v2.69.0"></a>
## [v2.69.0] - 2025-09-04
### Features
- **nrqldroprules:** deprecate package, add support for Pipeline Cloud Rule ID ([#1331](https://github.com/newrelic/newrelic-client-go/issues/1331))

<a name="v2.68.1"></a>
## [v2.68.1] - 2025-09-02
### Bug Fixes
- **change-tracking:** update depth to 2 and fix test failures ([#1325](https://github.com/newrelic/newrelic-client-go/issues/1325))

<a name="v2.68.0"></a>
## [v2.68.0] - 2025-09-01
### Features
- **pipelinecontrol:** add a new Pipeline Control package to the client ([#1330](https://github.com/newrelic/newrelic-client-go/issues/1330))

<a name="v2.67.1"></a>
## [v2.67.1] - 2025-08-05
### Bug Fixes
- **change-tracking:** remove test print ([#1323](https://github.com/newrelic/newrelic-client-go/issues/1323))

<a name="v2.67.0"></a>
## [v2.67.0] - 2025-07-29
### Features
- **dashboard:** support multiple account id's when creating dashboards ([#1319](https://github.com/newrelic/newrelic-client-go/issues/1319))

<a name="v2.66.0"></a>
## [v2.66.0] - 2025-07-21
### Features
- **tooltip-support:** add tooltip support in dashboard rawConfiguration ([#1320](https://github.com/newrelic/newrelic-client-go/issues/1320))

<a name="v2.65.0"></a>
## [v2.65.0] - 2025-07-14
### Features
- **change-tracking:** included change tracking api ([#1317](https://github.com/newrelic/newrelic-client-go/issues/1317))

<a name="v2.64.0"></a>
## [v2.64.0] - 2025-06-18
### Features
- **nrdb:** add support to the nrdb package to support facet, timeseries result unmarshaling ([#1312](https://github.com/newrelic/newrelic-client-go/issues/1312))

<a name="v2.63.0"></a>
## [v2.63.0] - 2025-06-17
### Features
- add Workflow Automation Destination and Channel types ([#1309](https://github.com/newrelic/newrelic-client-go/issues/1309))

<a name="v2.62.0"></a>
## [v2.62.0] - 2025-06-13
### Features
- **cloud:** added auto-discovery slug to cloud configure integration api ([#1307](https://github.com/newrelic/newrelic-client-go/issues/1307))

<a name="v2.61.2"></a>
## [v2.61.2] - 2025-06-09
### Bug Fixes
- fixed a distributed tracing bug in the browser package ([#1304](https://github.com/newrelic/newrelic-client-go/issues/1304))
- **build:** add PR checks similar to Terraform, fix failing tests, broken code and obsolete deps ([#1295](https://github.com/newrelic/newrelic-client-go/issues/1295))

<a name="v2.61.1"></a>
## [v2.61.1] - 2025-05-26
### Bug Fixes
- fixes users pkg ([#1298](https://github.com/newrelic/newrelic-client-go/issues/1298))

<a name="v2.61.0"></a>
## [v2.61.0] - 2025-05-22
### Features
- **alerts:** Add support for silent alerts via `disableHealthStatusReporting` ([#1296](https://github.com/newrelic/newrelic-client-go/issues/1296))

<a name="v2.60.0"></a>
## [v2.60.0] - 2025-05-12
### Features
- **entityrelationship:** add new package to manage entity relationships

<a name="v2.59.0"></a>
## [v2.59.0] - 2025-05-06
### Features
- add Metrics Client ([#1293](https://github.com/newrelic/newrelic-client-go/issues/1293))

<a name="v2.58.0"></a>
## [v2.58.0] - 2025-04-16
### Features
- **alerts:** adds support for `signalSeasonality` for NRQL baseline conditions ([#1291](https://github.com/newrelic/newrelic-client-go/issues/1291))

<a name="v2.57.0"></a>
## [v2.57.0] - 2025-04-14
### Features
- **alerts:** revert addition of support for nrql condition signal seasonality ([#1290](https://github.com/newrelic/newrelic-client-go/issues/1290))

<a name="v2.56.1"></a>
## [v2.56.1] - 2025-04-14
### Bug Fixes
- **synthetics:** Remove `omitempty` from `ResponseValidationText` to allow unsetting ([#1287](https://github.com/newrelic/newrelic-client-go/issues/1287))

<a name="v2.56.0"></a>
## [v2.56.0] - 2025-04-10
### Features
- **alerts:** add support for nrql condition signal seasonality ([#1280](https://github.com/newrelic/newrelic-client-go/issues/1280))

<a name="v2.55.4"></a>
## [v2.55.4] - 2025-04-03
### Bug Fixes
- **docs:** mock release for a test inclusive of a README update ([#1281](https://github.com/newrelic/newrelic-client-go/issues/1281))

<a name="v2.55.3"></a>
## [v2.55.3] - 2025-03-28
### Bug Fixes
- **build:** add codecov token to be used from secrets

<a name="v2.55.2"></a>
## [v2.55.2] - 2025-03-17
### Bug Fixes
- **release:** updated missing depenedencies in release script ([#1275](https://github.com/newrelic/newrelic-client-go/issues/1275))
- **release:** moved release script to make files to avoid race condition with release tags([#1274](https://github.com/newrelic/newrelic-client-go/issues/1274))

<a name="v2.55.1"></a>
## [v2.55.1] - 2025-03-14
### Bug Fixes
- **release:** updated release script to populate release notes ([#1273](https://github.com/newrelic/newrelic-client-go/issues/1273))

<a name="v2.55.0"></a>
## [v2.55.0] - 2025-03-10
### Features
- **application_settings:** update application settings to latest fields ([#1267](https://github.com/newrelic/newrelic-client-go/issues/1267))

<a name="v2.54.0"></a>
## [v2.54.0] - 2025-02-21
### Features
- **alerts:** Add support for prediction to NRQL alert conditions ([#1271](https://github.com/newrelic/newrelic-client-go/issues/1271))

<a name="v2.53.0"></a>
## [v2.53.0] - 2025-02-03
### Bug Fixes
- **dependabot_alerts:** Updated dependencies to fix dependabot alerts ([#1265](https://github.com/newrelic/newrelic-client-go/issues/1265))
- **integration-tests:** Resolve `newrelic-client-go` failed integration tests ([#1264](https://github.com/newrelic/newrelic-client-go/issues/1264))

### Features
- **workflows:** adds new `notification_trigger` type `INVESTIGATING` ([#1270](https://github.com/newrelic/newrelic-client-go/issues/1270))

<a name="v2.52.0"></a>
## [v2.52.0] - 2025-01-08
### Features
- Add optimized error message to the recipe event ([#1256](https://github.com/newrelic/newrelic-client-go/issues/1256))
- **alert_muting_rule:** Add `action_on_muting_rule_window_ended` attribute in `newrelic_alert_muting_rule` Terraform Resource ([#1259](https://github.com/newrelic/newrelic-client-go/issues/1259))

<a name="v2.51.3"></a>
## [v2.51.3] - 2024-11-11
### Bug Fixes
- **entities:** remove relationships field from entity, entities, entitySearch queries ([#1255](https://github.com/newrelic/newrelic-client-go/issues/1255))

<a name="v2.51.2"></a>
## [v2.51.2] - 2024-11-05
### Bug Fixes
- **goreleaser:** try replacing build with builds for v2 config ([#1254](https://github.com/newrelic/newrelic-client-go/issues/1254))

<a name="v2.51.1"></a>
## [v2.51.1] - 2024-11-05
### Bug Fixes
- **goreleaser:** print goreleaser version ([#1252](https://github.com/newrelic/newrelic-client-go/issues/1252))

<a name="v2.51.0"></a>
## [v2.51.0] - 2024-11-04
### Features
- **codeql:** fixes to the action with path changes ([#1238](https://github.com/newrelic/newrelic-client-go/issues/1238))

<a name="v2.50.1"></a>
## [v2.50.1] - 2024-10-28
### Bug Fixes
- **goreleaser:** update version constraint ([#1245](https://github.com/newrelic/newrelic-client-go/issues/1245))

<a name="v2.50.0"></a>
## [v2.50.0] - 2024-10-28
### Features
- minor change in README.md test commit ([#1244](https://github.com/newrelic/newrelic-client-go/issues/1244))

<a name="v2.49.0"></a>
## [v2.49.0] - 2024-10-28
### Bug Fixes
- **azure_link_account:** add update capability for all fields in azure cloud link account  ([#1242](https://github.com/newrelic/newrelic-client-go/issues/1242))

### Features
- **automation:** execute code generation when new API endpoints are detected
- **dashboards:** Added a new field in dashboards -> variable -> options called <excluded>

### Refactor
- **automation:** don't send slack message if no new endpoints

<a name="v2.48.2"></a>
## [v2.48.2] - 2024-10-04
### Bug Fixes
- Init agent library ([#1230](https://github.com/newrelic/newrelic-client-go/issues/1230))

<a name="v2.48.1"></a>
## [v2.48.1] - 2024-10-03
### Bug Fixes
- Expose agent library ([#1229](https://github.com/newrelic/newrelic-client-go/issues/1229))

<a name="v2.48.0"></a>
## [v2.48.0] - 2024-10-03
### Features
- Current agent release ([#1228](https://github.com/newrelic/newrelic-client-go/issues/1228))

<a name="v2.47.0"></a>
## [v2.47.0] - 2024-09-25
### Features
- **keytransaction:** add `keytransaction` package to support managing key transactions ([#1213](https://github.com/newrelic/newrelic-client-go/issues/1213))

<a name="v2.46.0"></a>
## [v2.46.0] - 2024-09-25
### Features
- **dashboards:** add support for `dataFormat` to `rawConfiguration` in dashboards

<a name="v2.45.0"></a>
## [v2.45.0] - 2024-09-10
### Features
- **synthetics:** adds `browsers` `devices` fields to some synthetics mutations, queries ([#1198](https://github.com/newrelic/newrelic-client-go/issues/1198))

<a name="v2.44.0"></a>
## [v2.44.0] - 2024-09-02
### Features
- **dashboard:** add support for initial sorting and refresh rate ([#1206](https://github.com/newrelic/newrelic-client-go/issues/1206))

<a name="v2.43.2"></a>
## [v2.43.2] - 2024-08-29
### Bug Fixes
- **dashboard:**  changing `DashboardWidgetLegend` -> `Enabled` field to pointer type ([#1210](https://github.com/newrelic/newrelic-client-go/issues/1210))

<a name="v2.43.1"></a>
## [v2.43.1] - 2024-08-14
### Bug Fixes
- **browseragent:** changing cookiesEnabled field to pointer type ([#1205](https://github.com/newrelic/newrelic-client-go/issues/1205))

<a name="v2.43.0"></a>
## [v2.43.0] - 2024-08-12
### Features
- **alerts:** add support for `ignoreOnExpectedTermination` in `expiration` to NRQL alert conditions ([#1180](https://github.com/newrelic/newrelic-client-go/issues/1180))
- **alerts:** add support for `titleTemplate` to NRQL alert conditions ([#1141](https://github.com/newrelic/newrelic-client-go/issues/1141))
- **automation:** report when NerdGraph API endpoints change signatures (i.e. breaking changes)

<a name="v2.42.1"></a>
## [v2.42.1] - 2024-08-07
### Bug Fixes
- **synthetics:** alter []SyntheticsCustomHeaderInput to allow empty custom headers in the request ([#1203](https://github.com/newrelic/newrelic-client-go/issues/1203))

<a name="v2.42.0"></a>
## [v2.42.0] - 2024-08-07
### Features
- **workflows:** added new notification channel type called servicenow_app ([#1200](https://github.com/newrelic/newrelic-client-go/issues/1200))

<a name="v2.41.3"></a>
## [v2.41.3] - 2024-08-05
### Bug Fixes
- change datatypes of to and from in line and table thresholds from float64 to string ([#1199](https://github.com/newrelic/newrelic-client-go/issues/1199))
- **entity:** regenerated 'entity' query max depth of 4 ([#1196](https://github.com/newrelic/newrelic-client-go/issues/1196))

<a name="v2.41.2"></a>
## [v2.41.2] - 2024-07-15
### Bug Fixes
- **dashboard:** update line and table widget threshold to and from fields datatypes to float64 ([#1192](https://github.com/newrelic/newrelic-client-go/issues/1192))

<a name="v2.41.1"></a>
## [v2.41.1] - 2024-07-11
### Bug Fixes
- **release:** testing org-wide branch protection rules in release process
- **release:** this is just a test

<a name="v2.41.0"></a>
## [v2.41.0] - 2024-07-10
### Features
- add `organizationRevokeSharedAccount` and `accountManagementCancelAccount` mutations ([#1187](https://github.com/newrelic/newrelic-client-go/issues/1187))

<a name="v2.40.0"></a>
## [v2.40.0] - 2024-07-03
### Features
- **nrqldroprules:** Added `GetDropRuleByID` function to fetch a droprule by specified ID ([#1183](https://github.com/newrelic/newrelic-client-go/issues/1183))

<a name="v2.39.1"></a>
## [v2.39.1] - 2024-07-01
### Bug Fixes
- **entities:** manual adjustments were needed for some dashboard types (pointers)

<a name="v2.39.0"></a>
## [v2.39.0] - 2024-07-01
### Features
- **nrql_alert_conditions:** Add data_account_id for NRQL alert condition ([#1184](https://github.com/newrelic/newrelic-client-go/issues/1184))

<a name="v2.38.0"></a>
## [v2.38.0] - 2024-06-27
### Features
- **automation:** alert [@hero](https://github.com/hero) when new API endpoints are detected
- **workflows:** add `updateOriginalMessage` to destination configurations ([#1177](https://github.com/newrelic/newrelic-client-go/issues/1177))

<a name="v2.37.1"></a>
## [v2.37.1] - 2024-06-12
### Bug Fixes
- **customeradministration:** change `OrganizationTargetIdInput` to have a pointer ([#1173](https://github.com/newrelic/newrelic-client-go/issues/1173))

<a name="v2.37.0"></a>
## [v2.37.0] - 2024-06-06
### Bug Fixes
- **changetracking:** alter the functioning of milli and nanoseconds in timestamps ([#1151](https://github.com/newrelic/newrelic-client-go/issues/1151))

### Features
- **entities:** add support for tags in EntityOutline types ([#1146](https://github.com/newrelic/newrelic-client-go/issues/1146))

<a name="v2.36.2"></a>
## [v2.36.2] - 2024-05-27
### Bug Fixes
- **infra:** Discard the infra test api call to prevent timeout issues ([#1145](https://github.com/newrelic/newrelic-client-go/issues/1145))

<a name="v2.36.1"></a>
## [v2.36.1] - 2024-05-23
### Bug Fixes
- **customeradministration:** change auth_domain filter types in the package to use pointers ([#1134](https://github.com/newrelic/newrelic-client-go/issues/1134))

<a name="v2.36.0"></a>
## [v2.36.0] - 2024-05-23
### Features
- **dashboards:** add yAxisRight, isLabelVisible and alter the functioning of thresholds ([#1144](https://github.com/newrelic/newrelic-client-go/issues/1144))

<a name="v2.35.0"></a>
## [v2.35.0] - 2024-05-22
### Bug Fixes
- **entity:** refactored entity GUID validation functions to remove padding ([#1143](https://github.com/newrelic/newrelic-client-go/issues/1143))

### Features
- **automation:** report new API features and updates via Tutone

<a name="v2.34.1"></a>
## [v2.34.1] - 2024-05-14
### Bug Fixes
- **entity:** Entity guid validation ([#1138](https://github.com/newrelic/newrelic-client-go/issues/1138))

<a name="v2.34.0"></a>
## [v2.34.0] - 2024-05-07
### Features
- **cloud:** addition of msElasticCache and AiPlatformIntegration integrations to AWS, GCP ([#1133](https://github.com/newrelic/newrelic-client-go/issues/1133))

### Refactor
- **nrdb:** clean up package, refactor structure, api endpoints and more ([#1132](https://github.com/newrelic/newrelic-client-go/issues/1132))

<a name="v2.33.0"></a>
## [v2.33.0] - 2024-04-30
### Features
- **build:** upgrade to Go v1.21 ✨ ([#1128](https://github.com/newrelic/newrelic-client-go/issues/1128))

<a name="v2.32.0"></a>
## [v2.32.0] - 2024-04-30
### Features
- **nrdb:** fix NRDB query to eliminate embeddedCharts ([#1131](https://github.com/newrelic/newrelic-client-go/issues/1131))

<a name="v2.31.0"></a>
## [v2.31.0] - 2024-04-19
### Features
- **destinations:** add custom header auth and secureUrl to destinations ([#1122](https://github.com/newrelic/newrelic-client-go/issues/1122))

<a name="v2.30.0"></a>
## [v2.30.0] - 2024-04-18
### Features
- addition of client code corresponding to packages customeradministration +2 ([#1126](https://github.com/newrelic/newrelic-client-go/issues/1126))

<a name="v2.29.0"></a>
## [v2.29.0] - 2024-04-09
### Features
- **entities:** add KeyTransactionEntity types (manual update)

<a name="v2.28.1"></a>
## [v2.28.1] - 2024-04-08
### Bug Fixes
- **entities:** include 'reporting' if part of entity search NRQL query

<a name="v2.28.0"></a>
## [v2.28.0] - 2024-04-02
### Bug Fixes
- **testhelpers:** add insights API mock URL to testconfig ([#1077](https://github.com/newrelic/newrelic-client-go/issues/1077))

### Features
- **entities:** add entity search query helpers

<a name="v2.27.0"></a>
## [v2.27.0] - 2024-03-20
### Bug Fixes
- Logger string interpolation ([#1104](https://github.com/newrelic/newrelic-client-go/issues/1104))
- **nrdb:** update to nrdb query to include more options ([#1110](https://github.com/newrelic/newrelic-client-go/issues/1110))

### Features
- synthetic monitors `MUTED` status EOL in newrelic-client-go ([#1102](https://github.com/newrelic/newrelic-client-go/issues/1102))
- **dashboards:** Added ignore time picker as options field to the dashboard variables ([#1106](https://github.com/newrelic/newrelic-client-go/issues/1106))
- **synthetics:** addition of support for next-gen runtime to step, cert check, broken links monitors ([#1114](https://github.com/newrelic/newrelic-client-go/issues/1114))

<a name="v2.26.1"></a>
## [v2.26.1] - 2024-03-01
### Bug Fixes
- Logger string interpolation ([#1104](https://github.com/newrelic/newrelic-client-go/issues/1104)) ([#1105](https://github.com/newrelic/newrelic-client-go/issues/1105))

<a name="v2.26.0"></a>
## [v2.26.0] - 2024-02-28
### Bug Fixes
- **group_management:** alter group management queries and add tests ([#1103](https://github.com/newrelic/newrelic-client-go/issues/1103))

### Features
- **destinations:** expose destination guid ([#1096](https://github.com/newrelic/newrelic-client-go/issues/1096))

<a name="v2.25.0"></a>
## [v2.25.0] - 2024-02-14
### Features
- **user_management:** addition of client code to manage users ([#1090](https://github.com/newrelic/newrelic-client-go/issues/1090))

<a name="v2.24.0"></a>
## [v2.24.0] - 2024-01-31
### Features
- **authentication_domains:** addition of client code to fetch auth domains ([#1086](https://github.com/newrelic/newrelic-client-go/issues/1086))

<a name="v2.23.0"></a>
## [v2.23.0] - 2023-12-19
### Features
- **synthetics:** addition of mutations for the monitor downtime feature ([#1070](https://github.com/newrelic/newrelic-client-go/issues/1070))

<a name="v2.22.2"></a>
## [v2.22.2] - 2023-11-14
### Bug Fixes
- integration tests ([#1064](https://github.com/newrelic/newrelic-client-go/issues/1064))

<a name="v2.22.1"></a>
## [v2.22.1] - 2023-11-07
### Bug Fixes
- **changetracking:** remove customAttributes from response fields to avoid error
- **synthetics_automated_testing:** update to schema based on NG updates ([#1062](https://github.com/newrelic/newrelic-client-go/issues/1062))

<a name="v2.22.0"></a>
## [v2.22.0] - 2023-10-25
### Features
- **changetracking:** Support custom attributes JSON ([#1060](https://github.com/newrelic/newrelic-client-go/issues/1060))

<a name="v2.21.2"></a>
## [v2.21.2] - 2023-10-16
### Bug Fixes
- **BrowserApplicationEntity:** addition of browserProperties fields ([#1035](https://github.com/newrelic/newrelic-client-go/issues/1035))

<a name="v2.21.1"></a>
## [v2.21.1] - 2023-09-05
### Bug Fixes
- Generate 1 more level depth for Synthetics simple browser ([#1057](https://github.com/newrelic/newrelic-client-go/issues/1057))
- **GenericEntityOutline:** addition of tags to attributes returned ([#1051](https://github.com/newrelic/newrelic-client-go/issues/1051))

<a name="v2.21.0"></a>
## [v2.21.0] - 2023-08-28
### Features
- Add capability id header for all NerdGraph requests ([#1049](https://github.com/newrelic/newrelic-client-go/issues/1049))
- **synthetics:** Update device emulation test for Simple Browser monitor synthetics ([#1054](https://github.com/newrelic/newrelic-client-go/issues/1054))

<a name="v2.20.0"></a>
## [v2.20.0] - 2023-08-01
### Features
- **changetracking:** add deployment custom attributes ([#1047](https://github.com/newrelic/newrelic-client-go/issues/1047))

<a name="v2.19.6"></a>
## [v2.19.6] - 2023-07-20
### Bug Fixes
- **apiaccess:** converted CreatedAt to int64

<a name="v2.19.5"></a>
## [v2.19.5] - 2023-07-20
### Bug Fixes
- **apiaccess:** Returning createdAt on keysearch

<a name="v2.19.4"></a>
## [v2.19.4] - 2023-07-11
### Bug Fixes
- **data_partition:** fix issue with enable attribute ([#1034](https://github.com/newrelic/newrelic-client-go/issues/1034))

<a name="v2.19.3"></a>
## [v2.19.3] - 2023-06-19
### Bug Fixes
- **dashboards:** remove omitempty on widgets to allow creating pages with no widgets ([#1033](https://github.com/newrelic/newrelic-client-go/issues/1033))

<a name="v2.19.2"></a>
## [v2.19.2] - 2023-06-14
### Bug Fixes
- **notifications:** add name to filter ([#1014](https://github.com/newrelic/newrelic-client-go/issues/1014))

<a name="v2.19.1"></a>
## [v2.19.1] - 2023-05-15
### Bug Fixes
- **alerts:** Removed omitempty from RunbookURL in NrqlConditionUpdateBase ([#1030](https://github.com/newrelic/newrelic-client-go/issues/1030))

<a name="v2.19.0"></a>
## [v2.19.0] - 2023-05-02
### Features
- **cloud:** update to queries, types to support azure monitor integration from TF ([#1029](https://github.com/newrelic/newrelic-client-go/issues/1029))

<a name="v2.18.1"></a>
## [v2.18.1] - 2023-04-18
### Bug Fixes
- **dashboards:** addition of the attribute zero to DashboardWidgetYAxisLeft

<a name="v2.18.0"></a>
## [v2.18.0] - 2023-04-18
### Bug Fixes
- **dashboard:** altering the 'Min' field in DashboardWidgetYAxisLeft ([#1020](https://github.com/newrelic/newrelic-client-go/issues/1020))

### Features
- **entity:** Add method to support the key transaction entities in client go. ([#1019](https://github.com/newrelic/newrelic-client-go/issues/1019))
- **synthetics:** add device emulation functionality ([#1017](https://github.com/newrelic/newrelic-client-go/issues/1017))

<a name="v2.17.1"></a>
## [v2.17.1] - 2023-04-07
### Bug Fixes
- **apiaccess:** fix to the error thrown when API Access Key creation fails ([#1015](https://github.com/newrelic/newrelic-client-go/issues/1015))

<a name="v2.17.0"></a>
## [v2.17.0] - 2023-03-22
### Features
- Update servicelevel schema ([#1012](https://github.com/newrelic/newrelic-client-go/issues/1012))

<a name="v2.16.0"></a>
## [v2.16.0] - 2023-03-15
### Features
- **account_management:** added account management api and test cases

<a name="v2.15.1"></a>
## [v2.15.1] - 2023-03-08
### Bug Fixes
- **one_dashboard:** added raw configuration properties to one dashboard resource ([#1001](https://github.com/newrelic/newrelic-client-go/issues/1001))

<a name="v2.15.0"></a>
## [v2.15.0] - 2023-03-07
### Bug Fixes
- **EpochTime:** fix unmarshal/marshal of empty EpochTime ([#997](https://github.com/newrelic/newrelic-client-go/issues/997))

### Features
- **agentapplication:** add ability to manage browser applications ([#991](https://github.com/newrelic/newrelic-client-go/issues/991))

<a name="v2.14.0"></a>
## [v2.14.0] - 2023-03-02
### Features
- **data_partition:** new attribute nrql and updated tests

<a name="v2.13.0"></a>
## [v2.13.0] - 2023-02-17
### Features
- **workflows:** expose workflow guid

<a name="v2.12.0"></a>
## [v2.12.0] - 2023-02-08
### Features
- facilitate additional service names for NewRelic-Requesting-Services request header via env var

<a name="v2.11.2"></a>
## [v2.11.2] - 2023-02-07
### Bug Fixes
- **synthetics:** resolved error targeting legacy runtimes

<a name="v2.11.1"></a>
## [v2.11.1] - 2023-01-25
### Bug Fixes
- **tutone:** release build

<a name="v2.11.0"></a>
## [v2.11.0] - 2023-01-18
### Features
- Add evaluation delay to nrql conditions

<a name="v2.10.0"></a>
## [v2.10.0] - 2023-01-06
### Features
- Remove value_function from NRQL Alert Conditions
- **parsingrule:** generated code and tests

<a name="v2.9.0"></a>
## [v2.9.0] - 2022-12-16
### Bug Fixes
- **workflows:** fix issue that prevents creation of disabled workflows

### Features
- **testgrok:** generated code and added tests

<a name="v2.8.0"></a>
## [v2.8.0] - 2022-12-14
### Bug Fixes
- **dashboards:** make some fields on variables nullable

### Features
- **workflows:** support a new flag to optionally disable channel deletion on workflow updates/deletes

<a name="v2.7.0"></a>
## [v2.7.0] - 2022-12-06
### Bug Fixes
- **synthetics_secure_cred:** handle null timestamp

<a name="v2.6.1"></a>
## [v2.6.1] - 2022-12-06
### Bug Fixes
- **drop_rules:** generated errors from schema :bug:

<a name="v2.6.0"></a>
## [v2.6.0] - 2022-11-30
### Features
- **dashboard:** update dashboard entity query with variables

<a name="v2.5.0"></a>
## [v2.5.0] - 2022-11-30
### Features
- **Workflow:** Add notification triggers
- **cloud:** add newly supported cloud service integrations
- **dashboards:** add variables
- **data_pratition_rule:** Added Data partition and tests

<a name="v2.4.0"></a>
## [v2.4.0] - 2022-11-16
### Features
- **obfuscation_rule:** Added Obfuscation rule and tests

<a name="v2.3.0"></a>
## [v2.3.0] - 2022-11-10
### Features
- **obfuscation_expression:** added obfuscation expression

<a name="v2.2.2"></a>
## [v2.2.2] - 2022-11-04
### Bug Fixes
- **workflows:** allow to remove enrichments from workflows

<a name="v2.2.1"></a>
## [v2.2.1] - 2022-11-03
### Bug Fixes
- **nrdb:** Use actor.queryHistory NG endpoint

<a name="v2.2.0"></a>
## [v2.2.0] - 2022-10-28
### Features
- add ChangeTracking to NewRelic client

<a name="v2.1.0"></a>
## [v2.1.0] - 2022-10-25
### Features
- **change_tracking:** add change tracking create deployment endpoint

<a name="v2.0.3"></a>
## [v2.0.3] - 2022-10-24
### Bug Fixes
- remove slug field which is causing timeout

<a name="v2.0.2"></a>
## [v2.0.2] - 2022-10-19
### Bug Fixes
- update module path to have v2
- **alert_conditions:** added missing userdef values

<a name="v2.0.1"></a>
## [v2.0.1] - 2022-10-18
### Bug Fixes
- update module path to have v2

<a name="v2.0.0"></a>
## [v2.0.0] - 2022-10-17
<a name="v1.1.0"></a>
## [v1.1.0] - 2022-10-17
### Features
- get workload collection

<a name="v1.0.0"></a>
## [v1.0.0] - 2022-09-26
### Bug Fixes
- add servicelevel select types as nullables

### Features
- update servicelevel model

<a name="v0.91.3"></a>
## [v0.91.3] - 2022-09-26
### Bug Fixes
- make some workloads types nullable

<a name="v0.91.2"></a>
## [v0.91.2] - 2022-09-22
### Bug Fixes
- omit filter ID when empty

<a name="v0.91.1"></a>
## [v0.91.1] - 2022-09-07
### Bug Fixes
- **notifications:** add missing destination type via tutone

<a name="v0.91.0"></a>
## [v0.91.0] - 2022-08-17
### Features
- **dashboards:** added RawConfiguration structure

<a name="v0.90.0"></a>
## [v0.90.0] - 2022-08-15
### Bug Fixes
- **notifications:** fix tests

### Features
- **workflows:** fix lint
- **workflows:** fix intgration tests
- **workflows:** fix unit tests + add readme
- **workflows:** add workflows API - fix tests
- **workflows:** add workflows API

<a name="v0.89.1"></a>
## [v0.89.1] - 2022-08-15
### Bug Fixes
- **notifications:** add fileds to error interface

<a name="v0.89.0"></a>
## [v0.89.0] - 2022-08-01
### Features
- **synthetics:** generate code for queries synthetics.script and synthetics.steps

<a name="v0.88.1"></a>
## [v0.88.1] - 2022-07-24
### Bug Fixes
- **destinations:** change credentials type to pointer

<a name="v0.88.0"></a>
## [v0.88.0] - 2022-07-15
### Bug Fixes
- **channels:** add small fix
- **channels:** add unit tests
- **channels:** add integration tests and destinations missing data

### Features
- **channels:** remove duplicate declartion
- **channels:** remove duplicate code and fix tests
- **channels:** add notifications channels API using tutone tool

<a name="v0.87.1"></a>
## [v0.87.1] - 2022-07-14
### Bug Fixes
- **muting_rules:** client not setting err.NotFound

<a name="v0.87.0"></a>
## [v0.87.0] - 2022-07-13
### Bug Fixes
- **destinations:** add integration tests and unit tests
- **destinations:** add integration tests + small fix for union type
- **destinations:** use tutone generator
- **destinations:** fix lint

### Features
- **destinations:** fix tests
- **destinations:** fix tests
- **notifications:** add notifications destinations api calls

<a name="v0.86.5"></a>
## [v0.86.5] - 2022-07-11
### Bug Fixes
- change private location GUID to string from int

<a name="v0.86.4"></a>
## [v0.86.4] - 2022-07-08
### Bug Fixes
- remove deprecated field from service level query
- remove deprecated field from service level query
- **synthetics:** use *bool type to avoid removing false values

<a name="v0.86.3"></a>
## [v0.86.3] - 2022-06-27
### Bug Fixes
- **synthetics:** use *bool type to avoid removing false values

<a name="v0.86.2"></a>
## [v0.86.2] - 2022-06-23
### Bug Fixes
- remove deprecated field from service level query

<a name="v0.86.1"></a>
## [v0.86.1] - 2022-06-15
### Bug Fixes
- remove deprecated field from service level query

<a name="v0.86.0"></a>
## [v0.86.0] - 2022-06-06
### Features
- **http:** Add retry condition for graphql TOO_MANY_REQUESTS error response on json

<a name="v0.85.0"></a>
## [v0.85.0] - 2022-05-24
### Features
- added testing scripts to test synthetic monitors
- added synthetics monitoring

<a name="v0.84.0"></a>
## [v0.84.0] - 2022-05-23
### Features
- **build:** upgrade to Go 1.18

<a name="v0.83.0"></a>
## [v0.83.0] - 2022-05-23
### Features
- Generate the client code for synthetics private locations

<a name="v0.82.0"></a>
## [v0.82.0] - 2022-05-23
### Features
- add entitySearch with query parameter

<a name="v0.81.0"></a>
## [v0.81.0] - 2022-05-23
### Documentation Updates
- add upgrade instructions and update example usage steps

### Features
- **dashboards:** remove deprecated and disabled legacy dashboards REST API methods

<a name="v0.80.0"></a>
## [v0.80.0] - 2022-05-12
### Features
- **alerts:** Adds 3 term threshold operators for NRQL conditions

<a name="v0.79.0"></a>
## [v0.79.0] - 2022-05-09
### Features
- add synthetics secure credentials GraphQL API

<a name="v0.78.0"></a>
## [v0.78.0] - 2022-04-28
### Features
- Expose EntityGUID on NRQL Conditions when using NerdGraph.

<a name="v0.77.0"></a>
## [v0.77.0] - 2022-04-28
### Documentation Updates
- update minimum Go version requirement in development section
- Update example in readme to compile and run with v0.73.0

### Features
- Expose EntityGUID on NRQL Conditions.

<a name="v0.76.0"></a>
## [v0.76.0] - 2022-04-26
### Features
- **build:** compile on Go 1.17.x

<a name="v0.75.0"></a>
## [v0.75.0] - 2022-04-13
### Features
- **errors:** handle 402 payment required HTTP response scenario

<a name="v0.74.2"></a>
## [v0.74.2] - 2022-03-23
### Bug Fixes
- use correct input type for cloud disable integrations mutation

<a name="v0.74.1"></a>
## [v0.74.1] - 2022-03-04
### Bug Fixes
- remove integrations from getLinkedAccounts query

<a name="v0.74.0"></a>
## [v0.74.0] - 2022-03-03
### Features
- **auth:** Add X-Account-ID header if value exists in request context

<a name="v0.73.0"></a>
## [v0.73.0] - 2022-02-09
### Features
- **entities:** add new entity types

<a name="v0.72.0"></a>
## [v0.72.0] - 2022-02-01
### Features
- **nrql_conditions:** add optional SlideBy field to signal

<a name="v0.71.0"></a>
## [v0.71.0] - 2022-01-25
### Features
- **cloud:** add query to get a single linked account

<a name="v0.70.0"></a>
## [v0.70.0] - 2022-01-19
### Features
- **installevents:** add recipe event metadata field, update mutation via tutone

<a name="v0.69.0"></a>
## [v0.69.0] - 2021-12-28
### Features
- **events:** Add license key authorization for the Event API

<a name="v0.68.3"></a>
## [v0.68.3] - 2021-12-03
### Bug Fixes
- **dashboards:** make billboard widget thresholds optional, add test cases around them

<a name="v0.68.2"></a>
## [v0.68.2] - 2021-12-03
### Bug Fixes
- **entities:** handle deprecated field errors in tests
- **http:** check if the NerdGraph error is a deprecation warning, and still pass on the response (with error)

<a name="v0.68.1"></a>
## [v0.68.1] - 2021-11-29
### Bug Fixes
- **release:** use our changelog for release notes

<a name="v0.68.0"></a>
## [v0.68.0] - 2021-10-22
### Features
- use improved error handling for muting rules

<a name="v0.67.0"></a>
## [v0.67.0] - 2021-10-21
### Features
- **tags:** added method to get only mutable tags

<a name="v0.66.2"></a>
## [v0.66.2] - 2021-10-21
### Bug Fixes
- update NRQL query for alert condition tests

<a name="v0.66.1"></a>
## [v0.66.1] - 2021-10-07
### Bug Fixes
- use pointer for EvaluationOffset

<a name="v0.66.0"></a>
## [v0.66.0] - 2021-10-06
### Features
- Provide additional context in GraphQL errors for Alerts operations

<a name="v0.65.0"></a>
## [v0.65.0] - 2021-10-05
### Features
- **alerts:** streaming triggers for nrql alerts

<a name="v0.64.1"></a>
## [v0.64.1] - 2021-09-28
### Bug Fixes
- add spell check for auto-generated CHANGELOG.md
- let goreleaser generate the release notes. git-chglog for CHANGELOG

<a name="v0.64.0"></a>
## [v0.64.0] - 2021-09-28
### Features
- **install:** add DETECTED status via code gen

<a name="v0.63.5"></a>
## [v0.63.5] - 2021-09-27
### Bug Fixes
- **build:** more error checking in the release script

<a name="v0.63.4"></a>
## [v0.63.4] - 2021-09-24
### Bug Fixes
- release test

<a name="v0.63.3"></a>
## [v0.63.3] - 2021-09-24
### Bug Fixes
- release test
- release test

<a name="v0.63.2"></a>
## [v0.63.2] - 2021-09-23
### Bug Fixes
- release test

<a name="v0.63.1"></a>
## [v0.63.1] - 2021-09-23
### Bug Fixes
- add additional output to verify release tag info
- update to correct current version in version.go
- use all branches for tag-mode  to get current and next tag with svu
- release test
- release test
- release test
- release
- **servicelevel:** Update code gen strategy

<a name="v0.63.0"></a>
## [v0.63.0] - 2021-09-21
### Bug Fixes
- **servicelevel:** avoid import cycle
- **servicelevel:** Initialize service level API with config

### Features
- **servicelevel:** update code gen strategy
- **servicelevel:** generate servicelevel API

### Refactor
- move EntityGUID to a common package

<a name="v0.62.1"></a>
## [v0.62.1] - 2021-08-04
### Bug Fixes
- update error handling to reflect schema changes

<a name="v0.62.0"></a>
## [v0.62.0] - 2021-08-03
### Bug Fixes
- override ID type as string

### Features
- add installstatus schema for install-events-service

### Refactor
- delete installationeventresult if statement

<a name="v0.61.4"></a>
## [v0.61.4] - 2021-07-28
### Bug Fixes
- update error handling code for alert policies

<a name="v0.61.3"></a>
## [v0.61.3] - 2021-07-28
### Bug Fixes
- retire usages of deprecated error schema

<a name="v0.61.2"></a>
## [v0.61.2] - 2021-07-22
### Bug Fixes
- **logging:** export LogrusLogger for use in other projects

<a name="v0.61.1"></a>
## [v0.61.1] - 2021-07-20
### Bug Fixes
- **dashboard:** skip DashboardBillboardWidgetThresholdInput not DashboardBillboardWidgetConfigurationInput

### Refactor
- **tutone:** Add error wrapping to mutation results

<a name="v0.61.0"></a>
## [v0.61.0] - 2021-07-13
### Bug Fixes
- type for validation duration
- **dashboard:** Linked entities must be the page GUID, update the test

### Features
- **Error:** Add InvalidInput error
- **installevents:** start package to track install-events-service

<a name="v0.60.2"></a>
## [v0.60.2] - 2021-06-29
<a name="v0.60.1"></a>
## [v0.60.1] - 2021-06-28
### Features
- **apiaccess:** add context-aware methods for insights keys
- **entity:** Add EntityInterface.GetTags()
- **events:** add context-aware method to event creation method

### Refactor
- **workloads:** Generate workload code, deprecate old functions

<a name="v0.60.0"></a>
## [v0.60.0] - 2021-06-11
### Bug Fixes
- **client:** remove over-strict cast

### Features
- add context-aware methods

<a name="v0.59.4"></a>
## [v0.59.4] - 2021-06-10
### Bug Fixes
- **dashboards:** remove goldenTags from dashboard query

<a name="v0.59.3"></a>
## [v0.59.3] - 2021-06-10
### Bug Fixes
- **dashboards:** remove goldenMetrics from dashboard query

<a name="v0.59.2"></a>
## [v0.59.2] - 2021-06-10
### Bug Fixes
- **apm:** remove applicationsREST unused funcs
- **http:** Look inside response body for downstream NotFound errors

### Features
- **apm:** allow passing context to apm methods
- **apm:** allow passing context to applicationsREST funcs

<a name="v0.59.1"></a>
## [v0.59.1] - 2021-05-24
### Bug Fixes
- **region:** fix insights key management api url

<a name="v0.59.0"></a>
## [v0.59.0] - 2021-05-13
### Features
- **apiaccess:** add methods for managing insights insert keys
- **serialization:** Add Unix() command to EpochTime

<a name="v0.58.5"></a>
## [v0.58.5] - 2021-04-27
### Bug Fixes
- **graphql:** retry on server errors

<a name="v0.58.4"></a>
## [v0.58.4] - 2021-04-15
### Bug Fixes
- **cloud:** regenerate types
- **synthetics:** add paging to monitors resource

<a name="v0.58.3"></a>
## [v0.58.3] - 2021-02-19
### Bug Fixes
- **dashboards:** return an error.NotFound instead of nil

<a name="v0.58.2"></a>
## [v0.58.2] - 2021-02-18
### Bug Fixes
- **dashboards:** Prevent nil dereference on GetDashboardEntity

<a name="v0.58.1"></a>
## [v0.58.1] - 2021-02-17
### Bug Fixes
- **dashboards:** Return rawConfiguration on get, needed for all viz types
- **nrqldroprules:** Actually return Nrqldroprules client
- **region_constants:** corrected insightsBaseURL for EU

### Features
- **alerts:** adding id to alertsMutingRulesQuery

<a name="v0.58.0"></a>
## [v0.58.0] - 2021-02-12
### Bug Fixes
- **typegen:** do not attempt to unmarshal null data

### Features
- **nrqldroprules:** Implement NrqlDropRules

### Refactor
- Update all code-gen unmarshals with new typegen template

<a name="v0.57.2"></a>
## [v0.57.2] - 2021-02-01
### Refactor
- **alerts:** remove omitEmpty from MutingRuleScheduleUpdateInput

<a name="v0.57.1"></a>
## [v0.57.1] - 2021-01-29
### Refactor
- Tutone auto-naming conflict with schema
- EpochTime as a pointer to allow for null value in JSON unmarshaling

<a name="v0.57.0"></a>
## [v0.57.0] - 2021-01-27
### Bug Fixes
- **dashboards:** MANUAL CHANGE: remove queries until it is out of the schema

### Features
- **users:** Add users package, and replace references in existing packages
- **users:** Add users package

<a name="v0.56.2"></a>
## [v0.56.2] - 2021-01-22
### Bug Fixes
- **dashboard:** Fetch permalink for dashboards

<a name="v0.56.1"></a>
## [v0.56.1] - 2021-01-22
### Bug Fixes
- **http:** fix panics when resp is nil

### Refactor
- **testhelpers:** Remove hard-coded TestAccountID

<a name="v0.56.0"></a>
## [v0.56.0] - 2021-01-22
### Bug Fixes
- **http:** display underlying errors on max retries

### Features
- **alerts:** Add muting rule schedule fields
- **dashboards:** add linkedEntities to getDashboardEntityQuery

<a name="v0.55.8"></a>
## [v0.55.8] - 2021-01-15
### Refactor
- **dashboards:** Use nrqlQueries in place of queries

<a name="v0.55.7"></a>
## [v0.55.7] - 2021-01-15
<a name="v0.55.6"></a>
## [v0.55.6] - 2021-01-15
<a name="v0.55.5"></a>
## [v0.55.5] - 2021-01-14
### Bug Fixes
- **entities:** unmarshal Minutes as an int

<a name="v0.55.4"></a>
## [v0.55.4] - 2021-01-13
### Bug Fixes
- **http:** move logger initialization to NewClient()

<a name="v0.55.3"></a>
## [v0.55.3] - 2021-01-12
### Bug Fixes
- **dashboards:** remove manual changes so code generation works again

<a name="v0.55.2"></a>
## [v0.55.2] - 2021-01-11
### Bug Fixes
- **http:** slightly better error message for 401 status code

<a name="v0.55.1"></a>
## [v0.55.1] - 2021-01-11
### Bug Fixes
- **dashboards:** DashboardWidgetConfigurationInput needs to be nullable

<a name="v0.55.0"></a>
## [v0.55.0] - 2021-01-05
### Features
- **entities:** Add more methods to Entity(Outline)Interface

<a name="v0.54.1"></a>
## [v0.54.1] - 2021-01-05
<a name="v0.54.0"></a>
## [v0.54.0] - 2021-01-04
### Features
- **entities:** Generate Getter helpers for EntityInterfaces

### Refactor
- **entities:** change the get functions to not be on ptrs

<a name="v0.53.0"></a>
## [v0.53.0] - 2020-12-28
### Bug Fixes
- **entities:** DashboardWidgetRawConfiguration custom unmarshal as []byte (raw JSON)
- **entities:** Force ID fields to be a string
- **nerdgraphclient:** Template fixes to prevent nil pointers
- **typegen:** Avoid nil pointer on custom UnmarshalJSON

### Features
- **dashboards:** Add GetDashboardEntity()
- **dashboards:** Auto-generate GraphQL code for dashboards (early access)
- **entities:** Code-gen tag mutations
- **entities:** Mostly code generated Entities queries (entities/entity/entitySearch)

### Refactor
- Add omitempty to nullable and input objects for all packages
- Selective generation on type in Entities, Infrastructure, and NerdStorage
- DRY up some of the time based items into nrtime
- **cloud:** Cloud auto-generating via tutone
- **entities:** Cleanup unused structs in entities
- **entities:** DRY up Nrdb from Entities
- **typegen:** If we override a type to be in another package, properly generate the Unmarshal func call

<a name="v0.52.0"></a>
## [v0.52.0] - 2020-12-08
### Bug Fixes
- **typegen:** Avoid nil pointer on unmarshall

### Features
- **accounts:** Add AccountReference
- **nrql_conditions:** add violation_time_limit_seconds

<a name="v0.51.0"></a>
## [v0.51.0] - 2020-12-01
### Features
- **alerts:** allow passing context to alerts methods
- **nerdgraph:** allow passing context to underlying client
- **nrdb:** allow passing context to nrdb query methods

<a name="v0.50.0"></a>
## [v0.50.0] - 2020-11-20
### Features
- **nerdgraph:** allow custom unmarshal structs for queries

<a name="v0.49.0"></a>
## [v0.49.0] - 2020-11-13
### Bug Fixes
- **graphql:** include downstream error retry condition

### Features
- **config:** add a local region

<a name="v0.48.1"></a>
## [v0.48.1] - 2020-11-10
### Bug Fixes
- **http:** include INTERNAL_SERVER_ERROR as a retry reason
- **muting_rules:** ensure updates to disable rule are respected

<a name="v0.48.0"></a>
## [v0.48.0] - 2020-11-04
### Bug Fixes
- **cloud:** manually update generated code to fix cloud account methods

### Documentation Updates
- **cloud:** add cloud account resource example

### Features
- **cloud:** add cloud domain to client API

<a name="v0.47.3"></a>
## [v0.47.3] - 2020-10-28
### Bug Fixes
- **apm:** add extra comments
- **apm:** fix linting error
- **apm:** add backwards compatible fix and test

<a name="v0.47.2"></a>
## [v0.47.2] - 2020-10-27
### Bug Fixes
- **alerts:** don't omitempty for muting rule's enabled field

<a name="v0.47.1"></a>
## [v0.47.1] - 2020-10-23
### Bug Fixes
- **alerts:** remove pagination from ListMultiLocationSyntheticsConditions
- **http:** include 500 errors in reasons to retry requests

<a name="v0.47.0"></a>
## [v0.47.0] - 2020-10-16
### Features
- **http:** retry on nerdgraph server timeout

<a name="v0.46.0"></a>
## [v0.46.0] - 2020-10-15
### Bug Fixes
- **alerts:** make error handling more resilient for alert policies
- **build:** update changelog action for improved standards
- **build:** use DTK token for auto-PR process
- **edge:** trace observer schema updates

### Documentation Updates
- update changelog

### Features
- **cloud:** include initial cloud client support

<a name="v0.45.0"></a>
## [v0.45.0] - 2020-10-05
### Documentation Updates
- update changelog

### Features
- **nrql alert condition:** add signal.aggregation_window

<a name="v0.44.0"></a>
## [v0.44.0] - 2020-10-02
### Documentation Updates
- update changelog

### Features
- remove admin API key as an authentication mechanism
- **application_instances:** add an application instance resource

<a name="v0.43.0"></a>
## [v0.43.0] - 2020-10-01
### Documentation Updates
- update changelog

### Features
- **synthetics:** change resources to use personal api keys

<a name="v0.42.1"></a>
## [v0.42.1] - 2020-09-30
### Bug Fixes
- **events:** dereference the data pointer

### Documentation Updates
- update supported Go information and test config

<a name="v0.42.0"></a>
## [v0.42.0] - 2020-09-23
### Features
- **alerts:** enable personal api key auth for infra conditions

<a name="v0.41.2"></a>
## [v0.41.2] - 2020-09-16
### Refactor
- **alerts:** remove widespread change, limit scope to only nrql condition error resp handling

<a name="v0.41.1"></a>
## [v0.41.1] - 2020-09-15
### Bug Fixes
- **http:** handle 'not found' downstream response

<a name="v0.41.0"></a>
## [v0.41.0] - 2020-09-11
### Bug Fixes
- **entities:** filter out read-only tag values

### Features
- **logs:** support insert key

<a name="v0.40.0"></a>
## [v0.40.0] - 2020-09-04
### Features
- **alerts:** add new fields 'expiration' and 'signal' to nrql_conditions

<a name="v0.39.0"></a>
## [v0.39.0] - 2020-08-27
### Features
- **logs:** implement log batch mode
- **logs:** implement log batch mode

<a name="v0.38.0"></a>
## [v0.38.0] - 2020-08-25
### Bug Fixes
- **changelog:** drop reviewers and assignees

### Documentation Updates
- update changelog

### Features
- **logs:** Implement Log API

<a name="v0.37.0"></a>
## [v0.37.0] - 2020-08-20
### Features
- **apiaccess:** add search api access keys method

<a name="v0.36.0"></a>
## [v0.36.0] - 2020-08-20
### Features
- **apiaccesskeys:** add new api access keys package

<a name="v0.35.1"></a>
## [v0.35.1] - 2020-08-03
### Bug Fixes
- **newrelic:** Allow just an insert key for the newrelic package

<a name="v0.35.0"></a>
## [v0.35.0] - 2020-08-03
### Features
- **events:** Batch event insertion

<a name="v0.34.0"></a>
## [v0.34.0] - 2020-07-30
### Bug Fixes
- **alerts:** infra condition threshold value should be pointer to support zero-value thresholds
- **alerts:** always send infra condition description even if its an empty string
- **region:** make ConfigRegion case-insensitive

### Documentation Updates
- update README.md

### Features
- **graphql:** capture validation error output from response

<a name="v0.33.2"></a>
## [v0.33.2] - 2020-07-27
### Bug Fixes
- **alerts:** infra condition threshold value should be pointer to support zero-value thresholds

<a name="v0.33.1"></a>
## [v0.33.1] - 2020-07-24
### Bug Fixes
- **alerts:** always send infra condition description even if its an empty string

<a name="v0.33.0"></a>
## [v0.33.0] - 2020-07-23
### Bug Fixes
- fix http client compression
- **build:** run generate-tutone once per make command

### Features
- add a custom event resource
- **alerts:** add description field to infra alert condition

<a name="v0.32.1"></a>
## [v0.32.1] - 2020-07-17
### Bug Fixes
- **accounts:** add accounts API to client
- **nrql_conditions:** proper zero-value marshalling for threshold

<a name="v0.32.0"></a>
## [v0.32.0] - 2020-07-10
### Features
- **accounts:** add an accounts resource
- **synthetics:** add a monitor locations resource

<a name="v0.31.3"></a>
## [v0.31.3] - 2020-07-02
### Bug Fixes
- **alerts:** allow a value of 0 for NRQL condition thresholds

<a name="v0.31.2"></a>
## [v0.31.2] - 2020-07-01
### Bug Fixes
- **alerts:** better nerdgraph operator usage

<a name="v0.31.1"></a>
## [v0.31.1] - 2020-06-26
### Bug Fixes
- **alerts:** drop incorrect MonitorID flag from MultiLocationSyntheticsCondition

<a name="v0.31.0"></a>
## [v0.31.0] - 2020-06-18
### Bug Fixes
- add goreleaser back
- remove goreleaser

### Features
- **entities:** surface underlying application IDs for mobile application entities
- **eventstometrics:** add EventsToMetrics

<a name="v0.30.2"></a>
## [v0.30.2] - 2020-06-15
### Bug Fixes
- **nrdb:** Add String() to Epoch* types
- **serialization:** set tz UTC on MarshalJSON
- **serialization:** Fix nanosecond handling, set tz UTC, add EpochTime.String()

<a name="v0.30.1"></a>
## [v0.30.1] - 2020-06-12
### Bug Fixes
- **deps:** revert goreleaser v0.138.0 (causing import issues)

<a name="v0.30.0"></a>
## [v0.30.0] - 2020-06-12
### Features
- **nrdb:** Fetch nrql query history

<a name="v0.29.1"></a>
## [v0.29.1] - 2020-06-10
### Bug Fixes
- **nrdb:** Correctly unwrap the graphql context for Query, better testing

<a name="v0.29.0"></a>
## [v0.29.0] - 2020-06-10
### Bug Fixes
- **nerdgraph:** Add omitempty for yaml output
- **typegen:** Remove some overrides so types are generated without modifications

### Features
- **nrdb:** Add nrdb.Query(acct, NRQL) via NerdGraph
- **typegen:** Add imports to typegen

<a name="v0.28.1"></a>
## [v0.28.1] - 2020-06-03
### Bug Fixes
- **alerts:** add missing Outlier type to condition NrqlConditionTypes enum

<a name="v0.28.0"></a>
## [v0.28.0] - 2020-06-03
### Bug Fixes
- **alerts:** rename NrqlConditionTerms to singular for accuracy

### Features
- **alerts:** add DeleteNrqlConditionMutation as a proxy to DeleteConditionMutation
- **alerts:** add UpdateNrqlConditionOutlierMutation method for updating outlier NRQL alert conditions
- **alerts:** add CreateNrqlConditionOutlierMutation method for creating outlier NRQL alert conditions

### Refactor
- **nrql_conditions:** conditionID should be a string for consistency in ng implementation

<a name="v0.27.1"></a>
## [v0.27.1] - 2020-05-29
### Bug Fixes
- **serialization:** EpochTime handling of Unix timestamp with milliseconds

### Refactor
- **alerts:** update NG-based condition IDs to string type

<a name="v0.27.0"></a>
## [v0.27.0] - 2020-05-28
### Refactor
- **alerts:** update NG-based policy and condition IDs to string types

<a name="v0.26.0"></a>
## [v0.26.0] - 2020-05-27
### Bug Fixes
- **http:** Use default User-Agent header if none provided
- **http:** Refactor config.Compression and use it in http

### Features
- **http:** Enable compression handling for requests, consolidate POST methods
- **region:** Add Insights insert API URLs

<a name="v0.25.1"></a>
## [v0.25.1] - 2020-05-20
### Bug Fixes
- **alerts:** fix json tag for NrqlConditionInput.ValueFunction

<a name="v0.25.0"></a>
## [v0.25.0] - 2020-05-20
### Features
- enable APIKS auth for alerts and plugins packages
- **alerts:** add nerdgraph-based alert condition deletion

<a name="v0.24.1"></a>
## [v0.24.1] - 2020-05-19
### Bug Fixes
- **apm:** don't return empty zero values for floats for MetricTimesliceValues which can be misleading
- **apm:** adjust MetricDataParams json tags to support query params as arrays

<a name="v0.24.0"></a>
## [v0.24.0] - 2020-05-15
### Features
- **edge:** add trace observer resource

<a name="v0.23.4"></a>
## [v0.23.4] - 2020-05-13
### Bug Fixes
- **alerts:** allow blank runbook URL to be sent

<a name="v0.23.3"></a>
## [v0.23.3] - 2020-05-12
### Bug Fixes
- **workloads:** fix query type for entity guid

<a name="v0.23.2"></a>
## [v0.23.2] - 2020-05-11
### Bug Fixes
- **alerts:** Updating incorrect AlertEvents params
- **region:** gracefully fall back to default region

<a name="v0.23.1"></a>
## [v0.23.1] - 2020-05-04
### Bug Fixes
- **alerts:** Alerts paging was incorrectly generating URLs

<a name="v0.23.0"></a>
## [v0.23.0] - 2020-05-01
### Bug Fixes
- **build:** Github actions `make lint` for commit messages need full history
- **workloads:** Removing deprecated field `name` from `entitySearchQuery`

### Documentation Updates
- **typegen:** Add a README for typegen

### Features
- **events:** add alert events

### Refactor
- **build:** Split up github actions a bit
- **typegen:** Split/refactor much of typegen internals

<a name="v0.22.0"></a>
## [v0.22.0] - 2020-04-23
### Bug Fixes
- **alerts:** include "equal" operator for NRQL condition terms

### Features
- **dashboards:** add grid_column_count field

### Refactor
- **workloads:** query with GUID per upstream API change

<a name="v0.21.1"></a>
## [v0.21.1] - 2020-04-15
### Bug Fixes
- **alerts:** return a NotFound error when policies are not found in NerdGraph

### Refactor
- **workloads:** remove deprecated fields

<a name="v0.21.0"></a>
## [v0.21.0] - 2020-04-06
### Bug Fixes
- **build:** goreleaser now supports libraries, remove hack in config, skip build

### Features
- **typegen:** handle scalar types

### Refactor
- **alerts:** Update go:generate for types signature
- **newrelic:** Continue to fail, do not log invalid region
- **region:** Better error types/more tests
- **region:** Split parsing / fetching of region data so Parse can be reused
- **typegen:** Filter internal context off of generated descriptions if present
- **typegen:** Break out some structs
- **typegen:** Move all schema generation stuff into typegen

<a name="v0.20.1"></a>
## [v0.20.1] - 2020-04-01
### Bug Fixes
- **alerts:** use pointers for the result struct

### Refactor
- **nerdstorage:** check scope ID for zero values during nerdstorage operations

<a name="v0.20.0"></a>
## [v0.20.0] - 2020-03-31
### Bug Fixes
- **alerts:** include missing types to generate
- **nrql_conditions:** fix spelling errors, fix types
- **region:** Better URL building
- **typegen:** ensure we also generate non-input fields
- **typegen:** ensure handling of slice/LIST types
- **typegen:** default to resolving nested types

### Documentation Updates
- **README:** Update example in README, closes [#225](https://github.com/newrelic/newrelic-client-go/issues/225)
- **nerdstorage:** add examples

### Features
- **alerts:** implement NerdGraph policy search
- **internal:** add additional error context to graphQLError
- **nerdstorage:** add a nerdstorage package and resource
- **region:** Add a region package
- **typegen:** handle types of Kind OBJECT

### Refactor
- **alerts:** fix alerts tests
- **alerts:** Move FQDN/URL creation into package, out of http client for REST
- **apm:** Move FQDN/URL creation into package, out of http client for REST
- **apm:** Update apm tests
- **config:** Remove unused config fields
- **config:** Add new func for config
- **config:** Migrate Region to pkg/region
- **dashboards:** Move FQDN/URL creation into package, out of http client for REST
- **dashboards:** Update dashboard tests
- **entities:** Update entities tests
- **http:** introduce a request-scoped API for NerdGraph queries
- **http:** Move HTTP client to use new region format
- **http:** Remove assumption that we are talking to a REST endpoint
- **nerdgraph:** Update nerdgraph tests
- **plugins:** Move FQDN/URL creation into package, out of http client for REST
- **plugins:** Update plugin tests
- **region:** Change access to config.Region to ensure it exists
- **synthetics:** Update synthetics tests
- **synthetics:** Move FQDN/URL creation into package, out of http client for REST
- **typegen:** Convert to using go generate to run typegen, `make generate` to test
- **workloads:** Update workloads tests

<a name="v0.19.0"></a>
## [v0.19.0] - 2020-03-25
### Bug Fixes
- **alerts:** policy update response test
- **workloads:** remove nullable struct fields unless necessary

### Features
- **alerts:** add search method for NRQL conditions
- **alerts:** add get method for query NRQL conditions
- **alerts:** add update methods for baseline and static NRQL conditions
- **alerts:** add create methods for baseline and static NRQL conditions
- **nerdgraph:** begin generating structs from schema

### Refactor
- **alerts:** consolidate Nrql condition structs for better reusability

<a name="v0.18.0"></a>
## [v0.18.0] - 2020-03-20
### Bug Fixes
- **workloads:** fix some bugs in the workloads implementation
- **workloads:** export the workloads API via the newrelic package

### Features
- **alerts:** implement muting rules

<a name="v0.17.1"></a>
## [v0.17.1] - 2020-03-18
### Bug Fixes
- **alerts:** add custom unmarshaling for ConditionTerm
- **workloads:** use epoch time for EntitySearchQuery.CreatedAt

<a name="v0.17.0"></a>
## [v0.17.0] - 2020-03-17
### Bug Fixes
- **workloads:** map non-nullable fields to structs correctly

### Documentation Updates
- **alerts:** add package-level documentation and examples
- **apm:** add package-level documentation and examples
- **client:** add synopses for all packages
- **config:** add package-level documentation
- **dashboards:** add package-level documentation and examples
- **entities:** add package-level documentation and examples
- **errors:** update package-level documentation
- **infrastructure:** add package-level documentation
- **nerdgraph:** add package-level documentation and examples
- **newrelic:** use single-letter vars for receivers
- **newrelic:** add package-level documentation and examples
- **plugins:** add package-level documentation and examples
- **synthetics:** add package-level documentation and examples

### Features
- **alerts:** implement graphql policy methods
- **workloads:** add update operation, rework integration test scenario
- **workloads:** add delete and duplicate mutations
- **workloads:** add a workload create operation
- **workloads:** add a workloads resource, list and get methods

### Refactor
- **alerts:** Fix lint issue
- **alerts:** add types for fields with known values
- **apm:** Move Application REST implementation, use interface
- **dashboards:** add types for fields with known values
- **http:** Move NewRequest, have it follow New* func format
- **http:** Move graphql code out to file
- **http:** Make all fields private, add some setters/getters, more tests
- **http:** Consolidate GraphQL client, rename to http.Client
- **http:** Move GraphQL into http.NewRelicClient as Query()

<a name="v0.16.0"></a>
## [v0.16.0] - 2020-03-11
### Bug Fixes
- **build:** Force pull tags after each checkout

### Documentation Updates
- update community support information

### Features
- **entities:** Add some more details from BrowserApplicationEntity
- **entities:** Return more data on ApmApplicationEntity, and be consistent in what we return between fetch and search

### Refactor
- **entities:** Change Entity.Type type... Add more to the ENUMs

<a name="v0.15.0"></a>
## [v0.15.0] - 2020-03-09
### Bug Fixes
- **apm:** remove unused field
- **build:** Remove working dir config for CircleCI
- **http:** allow overriding of service name

### Refactor
- **alert_conditions:** remove transient PolicyID from struct for consistency with API response
- **alerts:** use consistent types for incident timestamp fields
- **build:** Make the build system consistent with other projects
- **nrql_conditions:** remove transient PolicyID from struct for consistency with API response
- **plugins_conditions:** remove transient PolicyID from struct for consistency with API response

<a name="v0.14.0"></a>
## [v0.14.0] - 2020-03-05
### Features
- **newrelic:** add types for fields with well known values

<a name="v0.13.0"></a>
## [v0.13.0] - 2020-03-03
### Bug Fixes
- **entities:** include applicationId for ApmApplicationEntity results
- **entities:** Make ApplicationID optional in results, omit if not returned by the API
- **http:** create a new errorValue for every request

### Refactor
- **apm:** refactor deployments resource to use new auth strategy
- **http:** refactor client to a request-scoped config context

<a name="v0.12.0"></a>
## [v0.12.0] - 2020-02-28
### Bug Fixes
- **docs:** Fix the release badge

### Features
- **nerdgraph:** implement ability to make raw graphql query

### Refactor
- **alerts:** Move structs into implementing files
- **apm:** Move structs into implementing files
- **config:** BREAKING CHANGE: Change environment vars and rename APIKey to AdminApiKey
- **dashboards:** Move structs into implementing files
- **entities:** Move structs into implementing files
- **synthetics:** Move structs into implementing files

<a name="v0.11.0"></a>
## [v0.11.0] - 2020-02-27
### Features
- **http:** allow personal API keys to be used for alerts and APM resources

### Refactor
- **http:** refactor authentication out of http client

<a name="v0.10.1"></a>
## [v0.10.1] - 2020-02-20
### Bug Fixes
- **entities:** tags filter needs to use type TagValue in graphql query
- **newrelic:** Add option to set ServiceName in Config

<a name="v0.10.0"></a>
## [v0.10.0] - 2020-02-19
### Features
- **ci:** add release make target
- **ci:** the beginnings of some release automation
- **synthetics:** add secure credentials resource
- **synthetics:** implement label monitor support

<a name="v0.9.0"></a>
## [v0.9.0] - 2020-02-05
### Bug Fixes
- allow string representations of JSON for alert channel webhook and payload
- **http:** Clear client responses between pages

### Features
- **alerts:** Implement multi-location synthetics conditions
- **http:** add trace logging with additional request info

<a name="v0.8.0"></a>
## [v0.8.0] - 2020-01-29
### Bug Fixes
- **alerts:** ensure multiple channels can be added via /alerts_policy_channel.json endpoint ([#114](https://github.com/newrelic/newrelic-client-go/issues/114))

### Features
- **apm:** Add support application metric names and data

<a name="v0.7.1"></a>
## [v0.7.1] - 2020-01-24
### Bug Fixes
- **alerts:** handle more complex JSON structures in headers and/or payload
- **logging:** use global methods for the default logger rather than a logrus instance

### Refactor
- **entities:** rename SearchEntities params struct per convention
- **newrelic:** remove reference to pointer for http transport config

<a name="v0.7.0"></a>
## [v0.7.0] - 2020-01-23
### Features
- **newrelic:** add ConfigOptions for logging
- **newrelic:** add the ability to configure base URLs per API

### Refactor
- **newrelic:** incorporate code review feedback

<a name="v0.6.0"></a>
## [v0.6.0] - 2020-01-22
### Features
- **alerts:** add GetSyntheticsCondition method ([#105](https://github.com/newrelic/newrelic-client-go/issues/105))

<a name="v0.5.1"></a>
## [v0.5.1] - 2020-01-21
### Bug Fixes
- **alerts:** custom unmarshal of channel configuration Headers and Payload fields ([#102](https://github.com/newrelic/newrelic-client-go/issues/102))

<a name="v0.5.0"></a>
## [v0.5.0] - 2020-01-16
### Documentation Updates
- **newrelic:** update API key configuration documentation

### Refactor
- **newrelic:** validate that at least one API key is provided

<a name="v0.4.0"></a>
## [v0.4.0] - 2020-01-15
### Bug Fixes
- retry HTTP requests on 429 status codes

### Features
- **entities:** add entities search and entity tagging

### Refactor
- update test helpers to use new mock server, consistent patterns in tests

<a name="v0.3.0"></a>
## [v0.3.0] - 2020-01-13
### Bug Fixes
- make use of ErrorNotFound type for Get methods that are based on List methods
- add policy ID to alert condition

### Documentation Updates
- update example
- **build:** Update README for commit message format
- **changelog:** Add auto-generation of CHANGELOG from git comments via `make changelog`

### Features
- add top-level logging package for convenience
- add option for JSON logging and fail gracefully when log level cannot be parsed
- introduce logging
- update monitor scripts with return design pattern, update tests

### Refactor
- update monitors to use return design pattern where applicable, update tests
- consistent use of pointers for &reqBody structs
- move logging config code into logging package
- use centralized test helpers and remove old ones
- update ListDashboards to return array of pointers, update Dashboard test to use assert
- update ListApplications to return array of pointers, update tests to use assert
- refactor to package-based types files
- remove config pointer references
- remove unnecessary else
- rescope vars for integration tests to avoid variable name conflicts
- update alert channels to return array of pointers, update tests to use assert
- update synthetics conditions to return array of pointers
- create a logger instance per package
- remove redundant 'alert' from file names
- remove redundant 'Alert' from naming convention
- incorporate code review feedback
- update alerts incidents to follow return design pattern, parallelize and use assert lib in alert incidents tests
- use require lib for dashboards integration tests
- **alerts:** Update mock server format, continue to have pkg helper
- **alerts:** Spike example of changes to the mock setup
- **config:** Change Region to a string, then parse with region package
- **newrelic:** Extract config setting to opts ... format
- **region:** Move region out of config into package, add Parse(string)

<a name="v0.2.0"></a>
## [v0.2.0] - 2020-01-08
### Documentation Updates
- update readme example

<a name="v0.1.0"></a>
## v0.1.0 - 2020-01-07
### Bug Fixes
- rename variables to fix redeclared error
- update unit tests to use new method sigs
- fix monitor ID type and GetMonitor URL
- http client needs to handle other 'success' response status codes such as 201
- add godoc as a dep, and a warning about GOPATH and godoc
- fix paging bug for v2 API
- **lint:** formatting fixes for linter

### Documentation Updates
- add alerts package docs
- temporarily checking in broken import paths in generated markdown docs
- add inline documentation
- add badges to README
- fill in missing inline documentation
- document some methods

### Features
- add DeletePluginCondition
- add CreatePluginCondition
- add UpdatePluginCondition
- add GetPluginCondition
- add ListPluginsConditions
- encode monitor script text
- add ability to use 'detailed' query param in ListPlugins method
- add GetPlugin
- add ListPlugins
- publicly expose error types
- finish components endpoints
- add Components
- add internal utils package, move IntArrayToString() util to new home
- add integration tests for key transactions
- add query param filters for ListKeyTransactions
- add GetKeyTransaction
- add ListKeyTransactions
- add DeleteLabel
- add CreateLabel
- add ListLabels, add GetLabel
- add DeleteDeployment
- add CreateDeployment
- add ListDeployments
- centralize apm test helpers
- add DeleteNrqlAlertCondition
- add UpdateNrqlAlertCondition
- add CreateNrqlAlertCondition
- add GetNrqlAlertCondition
- add ListNrqlAlertConditions
- add UpdateAlertPolicy
- add DeleteAlertCondition
- add CreateAlertCondition
- add GetAlertCondition
- add ListAlertConditions
- get infra condition integration tests passing
- add InfrastructureConditions
- add MonitorScripts
- add MonitorScript
- add DeleteAlertPolicyChannel, update unit tests, add integration test (might need to remove this)
- add alert policy channels
- add synthetics alert conditions
- add synthetics alert conditions
- add GetAlertChannel method
- add CreateAlertChannel, ListAlertChannels, DeleteAlertChannel
- add DeleteMonitor
- add UpdateMonitor
- add CreateMonitor
- add dashboards
- add DeleteAlertPolicy method
- add UpdateAlertPolicy method
- add CreateAlertPolicy method
- add GetAlertPolicy method
- add ListAlertPolicies method
- alerts package
- create remaining CRUD methods for application resource
- add new dependency-free client implementation
- add version.go per auto-versioning docs
- add ListAlertConditions for infrastructure
- add infra namespace
- add catchall newrelic package
- add New Relic environment enum
- maximize page size for ListMonitors
- add ListMonitors method for Synthetics monitors
- add application filtering for ListApplications
- get TestListApplications passing

### Refactor
- updates per code review
- use proper noun Plugins in naming convention
- update key txns to use new query string parsing mechanism
- simplify integration test scenarios for components
- move components to the plugins package
- move query string parsing to an external package
- represent query params as a struct rather than a map
- return slices of pointers instead of slices of structs
- simplify parameter handling logic
- optimize IntArrayToString() per review, add test cases
- add integration tests, update unit tests, links should be a pointer for omission
- Makefile cleanup
- optimize pushing to array of pointers
- refactor synthetics conditions to established patterns
- refactor alerts package to established patterns
- refactor synthetics package to established patterns
- update local var names for consistency
- update Epoch to EpochTime
- remove redundant 'Alert' from naming convention
- remove pointer from AlertChannelConfiguration
- utilize testify assert library, other minor refactors
- refactor unit tests to use testify assertions
- add concrete types for field with known possible values
- use Epoch type for date types instead of int64
- consolidate request body structs into one alertPolicyRequestBody
- no pointers for param fields
- integrate new http client
- simplify HTTP method signatures
- add the remaining HTTP methods
- rename the new client types
- remove the old resty-based client
- put new client in place for all resources
- make ListApplications use the new client
- move version into its own internal package for now
- incorporate linter suggestions
- clean up the configuration API for NewRelicClient
- restructuring project files
- extract cross cutting concern for apm resources
- extract paging implementation
- rename packages for clarity, promote Config to the public package

[Unreleased]: https://github.com/newrelic/newrelic-client-go/compare/v2.73.1...HEAD
[v2.73.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.73.0...v2.73.1
[v2.73.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.72.0...v2.73.0
[v2.72.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.71.0...v2.72.0
[v2.71.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.70.2...v2.71.0
[v2.70.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.70.1...v2.70.2
[v2.70.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.70.0...v2.70.1
[v2.70.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.69.0...v2.70.0
[v2.69.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.68.1...v2.69.0
[v2.68.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.68.0...v2.68.1
[v2.68.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.67.1...v2.68.0
[v2.67.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.67.0...v2.67.1
[v2.67.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.66.0...v2.67.0
[v2.66.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.65.0...v2.66.0
[v2.65.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.64.0...v2.65.0
[v2.64.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.63.0...v2.64.0
[v2.63.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.62.0...v2.63.0
[v2.62.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.61.2...v2.62.0
[v2.61.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.61.1...v2.61.2
[v2.61.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.61.0...v2.61.1
[v2.61.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.60.0...v2.61.0
[v2.60.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.59.0...v2.60.0
[v2.59.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.58.0...v2.59.0
[v2.58.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.57.0...v2.58.0
[v2.57.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.56.1...v2.57.0
[v2.56.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.56.0...v2.56.1
[v2.56.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.55.4...v2.56.0
[v2.55.4]: https://github.com/newrelic/newrelic-client-go/compare/v2.55.3...v2.55.4
[v2.55.3]: https://github.com/newrelic/newrelic-client-go/compare/v2.55.2...v2.55.3
[v2.55.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.55.1...v2.55.2
[v2.55.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.55.0...v2.55.1
[v2.55.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.54.0...v2.55.0
[v2.54.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.53.0...v2.54.0
[v2.53.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.52.0...v2.53.0
[v2.52.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.51.3...v2.52.0
[v2.51.3]: https://github.com/newrelic/newrelic-client-go/compare/v2.51.2...v2.51.3
[v2.51.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.51.1...v2.51.2
[v2.51.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.51.0...v2.51.1
[v2.51.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.50.1...v2.51.0
[v2.50.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.50.0...v2.50.1
[v2.50.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.49.0...v2.50.0
[v2.49.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.48.2...v2.49.0
[v2.48.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.48.1...v2.48.2
[v2.48.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.48.0...v2.48.1
[v2.48.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.47.0...v2.48.0
[v2.47.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.46.0...v2.47.0
[v2.46.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.45.0...v2.46.0
[v2.45.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.44.0...v2.45.0
[v2.44.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.43.2...v2.44.0
[v2.43.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.43.1...v2.43.2
[v2.43.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.43.0...v2.43.1
[v2.43.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.42.1...v2.43.0
[v2.42.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.42.0...v2.42.1
[v2.42.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.41.3...v2.42.0
[v2.41.3]: https://github.com/newrelic/newrelic-client-go/compare/v2.41.2...v2.41.3
[v2.41.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.41.1...v2.41.2
[v2.41.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.41.0...v2.41.1
[v2.41.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.40.0...v2.41.0
[v2.40.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.39.1...v2.40.0
[v2.39.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.39.0...v2.39.1
[v2.39.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.38.0...v2.39.0
[v2.38.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.37.1...v2.38.0
[v2.37.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.37.0...v2.37.1
[v2.37.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.36.2...v2.37.0
[v2.36.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.36.1...v2.36.2
[v2.36.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.36.0...v2.36.1
[v2.36.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.35.0...v2.36.0
[v2.35.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.34.1...v2.35.0
[v2.34.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.34.0...v2.34.1
[v2.34.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.33.0...v2.34.0
[v2.33.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.32.0...v2.33.0
[v2.32.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.31.0...v2.32.0
[v2.31.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.30.0...v2.31.0
[v2.30.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.29.0...v2.30.0
[v2.29.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.28.1...v2.29.0
[v2.28.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.28.0...v2.28.1
[v2.28.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.27.0...v2.28.0
[v2.27.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.26.1...v2.27.0
[v2.26.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.26.0...v2.26.1
[v2.26.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.25.0...v2.26.0
[v2.25.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.24.0...v2.25.0
[v2.24.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.23.0...v2.24.0
[v2.23.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.22.2...v2.23.0
[v2.22.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.22.1...v2.22.2
[v2.22.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.22.0...v2.22.1
[v2.22.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.21.2...v2.22.0
[v2.21.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.21.1...v2.21.2
[v2.21.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.21.0...v2.21.1
[v2.21.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.20.0...v2.21.0
[v2.20.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.19.6...v2.20.0
[v2.19.6]: https://github.com/newrelic/newrelic-client-go/compare/v2.19.5...v2.19.6
[v2.19.5]: https://github.com/newrelic/newrelic-client-go/compare/v2.19.4...v2.19.5
[v2.19.4]: https://github.com/newrelic/newrelic-client-go/compare/v2.19.3...v2.19.4
[v2.19.3]: https://github.com/newrelic/newrelic-client-go/compare/v2.19.2...v2.19.3
[v2.19.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.19.1...v2.19.2
[v2.19.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.19.0...v2.19.1
[v2.19.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.18.1...v2.19.0
[v2.18.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.18.0...v2.18.1
[v2.18.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.17.1...v2.18.0
[v2.17.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.17.0...v2.17.1
[v2.17.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.16.0...v2.17.0
[v2.16.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.15.1...v2.16.0
[v2.15.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.15.0...v2.15.1
[v2.15.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.14.0...v2.15.0
[v2.14.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.13.0...v2.14.0
[v2.13.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.12.0...v2.13.0
[v2.12.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.11.2...v2.12.0
[v2.11.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.11.1...v2.11.2
[v2.11.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.11.0...v2.11.1
[v2.11.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.10.0...v2.11.0
[v2.10.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.9.0...v2.10.0
[v2.9.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.8.0...v2.9.0
[v2.8.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.7.0...v2.8.0
[v2.7.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.6.1...v2.7.0
[v2.6.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.6.0...v2.6.1
[v2.6.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.5.0...v2.6.0
[v2.5.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.4.0...v2.5.0
[v2.4.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.3.0...v2.4.0
[v2.3.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.2.2...v2.3.0
[v2.2.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.2.1...v2.2.2
[v2.2.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.2.0...v2.2.1
[v2.2.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.1.0...v2.2.0
[v2.1.0]: https://github.com/newrelic/newrelic-client-go/compare/v2.0.3...v2.1.0
[v2.0.3]: https://github.com/newrelic/newrelic-client-go/compare/v2.0.2...v2.0.3
[v2.0.2]: https://github.com/newrelic/newrelic-client-go/compare/v2.0.1...v2.0.2
[v2.0.1]: https://github.com/newrelic/newrelic-client-go/compare/v2.0.0...v2.0.1
[v2.0.0]: https://github.com/newrelic/newrelic-client-go/compare/v1.1.0...v2.0.0
[v1.1.0]: https://github.com/newrelic/newrelic-client-go/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.91.3...v1.0.0
[v0.91.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.91.2...v0.91.3
[v0.91.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.91.1...v0.91.2
[v0.91.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.91.0...v0.91.1
[v0.91.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.90.0...v0.91.0
[v0.90.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.89.1...v0.90.0
[v0.89.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.89.0...v0.89.1
[v0.89.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.88.1...v0.89.0
[v0.88.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.88.0...v0.88.1
[v0.88.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.87.1...v0.88.0
[v0.87.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.87.0...v0.87.1
[v0.87.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.86.5...v0.87.0
[v0.86.5]: https://github.com/newrelic/newrelic-client-go/compare/v0.86.4...v0.86.5
[v0.86.4]: https://github.com/newrelic/newrelic-client-go/compare/v0.86.3...v0.86.4
[v0.86.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.86.2...v0.86.3
[v0.86.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.86.1...v0.86.2
[v0.86.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.86.0...v0.86.1
[v0.86.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.85.0...v0.86.0
[v0.85.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.84.0...v0.85.0
[v0.84.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.83.0...v0.84.0
[v0.83.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.82.0...v0.83.0
[v0.82.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.81.0...v0.82.0
[v0.81.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.80.0...v0.81.0
[v0.80.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.79.0...v0.80.0
[v0.79.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.78.0...v0.79.0
[v0.78.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.77.0...v0.78.0
[v0.77.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.76.0...v0.77.0
[v0.76.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.75.0...v0.76.0
[v0.75.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.74.2...v0.75.0
[v0.74.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.74.1...v0.74.2
[v0.74.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.74.0...v0.74.1
[v0.74.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.73.0...v0.74.0
[v0.73.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.72.0...v0.73.0
[v0.72.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.71.0...v0.72.0
[v0.71.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.70.0...v0.71.0
[v0.70.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.69.0...v0.70.0
[v0.69.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.68.3...v0.69.0
[v0.68.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.68.2...v0.68.3
[v0.68.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.68.1...v0.68.2
[v0.68.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.68.0...v0.68.1
[v0.68.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.67.0...v0.68.0
[v0.67.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.66.2...v0.67.0
[v0.66.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.66.1...v0.66.2
[v0.66.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.66.0...v0.66.1
[v0.66.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.65.0...v0.66.0
[v0.65.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.64.1...v0.65.0
[v0.64.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.64.0...v0.64.1
[v0.64.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.63.5...v0.64.0
[v0.63.5]: https://github.com/newrelic/newrelic-client-go/compare/v0.63.4...v0.63.5
[v0.63.4]: https://github.com/newrelic/newrelic-client-go/compare/v0.63.3...v0.63.4
[v0.63.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.63.2...v0.63.3
[v0.63.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.63.1...v0.63.2
[v0.63.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.63.0...v0.63.1
[v0.63.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.62.1...v0.63.0
[v0.62.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.62.0...v0.62.1
[v0.62.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.61.4...v0.62.0
[v0.61.4]: https://github.com/newrelic/newrelic-client-go/compare/v0.61.3...v0.61.4
[v0.61.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.61.2...v0.61.3
[v0.61.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.61.1...v0.61.2
[v0.61.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.61.0...v0.61.1
[v0.61.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.60.2...v0.61.0
[v0.60.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.60.1...v0.60.2
[v0.60.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.60.0...v0.60.1
[v0.60.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.59.4...v0.60.0
[v0.59.4]: https://github.com/newrelic/newrelic-client-go/compare/v0.59.3...v0.59.4
[v0.59.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.59.2...v0.59.3
[v0.59.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.59.1...v0.59.2
[v0.59.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.59.0...v0.59.1
[v0.59.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.58.5...v0.59.0
[v0.58.5]: https://github.com/newrelic/newrelic-client-go/compare/v0.58.4...v0.58.5
[v0.58.4]: https://github.com/newrelic/newrelic-client-go/compare/v0.58.3...v0.58.4
[v0.58.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.58.2...v0.58.3
[v0.58.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.58.1...v0.58.2
[v0.58.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.58.0...v0.58.1
[v0.58.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.57.2...v0.58.0
[v0.57.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.57.1...v0.57.2
[v0.57.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.57.0...v0.57.1
[v0.57.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.56.2...v0.57.0
[v0.56.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.56.1...v0.56.2
[v0.56.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.56.0...v0.56.1
[v0.56.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.55.8...v0.56.0
[v0.55.8]: https://github.com/newrelic/newrelic-client-go/compare/v0.55.7...v0.55.8
[v0.55.7]: https://github.com/newrelic/newrelic-client-go/compare/v0.55.6...v0.55.7
[v0.55.6]: https://github.com/newrelic/newrelic-client-go/compare/v0.55.5...v0.55.6
[v0.55.5]: https://github.com/newrelic/newrelic-client-go/compare/v0.55.4...v0.55.5
[v0.55.4]: https://github.com/newrelic/newrelic-client-go/compare/v0.55.3...v0.55.4
[v0.55.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.55.2...v0.55.3
[v0.55.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.55.1...v0.55.2
[v0.55.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.55.0...v0.55.1
[v0.55.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.54.1...v0.55.0
[v0.54.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.54.0...v0.54.1
[v0.54.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.53.0...v0.54.0
[v0.53.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.52.0...v0.53.0
[v0.52.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.51.0...v0.52.0
[v0.51.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.50.0...v0.51.0
[v0.50.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.49.0...v0.50.0
[v0.49.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.48.1...v0.49.0
[v0.48.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.48.0...v0.48.1
[v0.48.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.47.3...v0.48.0
[v0.47.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.47.2...v0.47.3
[v0.47.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.47.1...v0.47.2
[v0.47.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.47.0...v0.47.1
[v0.47.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.46.0...v0.47.0
[v0.46.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.45.0...v0.46.0
[v0.45.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.44.0...v0.45.0
[v0.44.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.43.0...v0.44.0
[v0.43.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.42.1...v0.43.0
[v0.42.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.42.0...v0.42.1
[v0.42.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.41.2...v0.42.0
[v0.41.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.41.1...v0.41.2
[v0.41.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.41.0...v0.41.1
[v0.41.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.40.0...v0.41.0
[v0.40.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.39.0...v0.40.0
[v0.39.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.38.0...v0.39.0
[v0.38.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.37.0...v0.38.0
[v0.37.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.36.0...v0.37.0
[v0.36.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.35.1...v0.36.0
[v0.35.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.35.0...v0.35.1
[v0.35.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.34.0...v0.35.0
[v0.34.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.33.2...v0.34.0
[v0.33.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.33.1...v0.33.2
[v0.33.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.33.0...v0.33.1
[v0.33.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.32.1...v0.33.0
[v0.32.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.32.0...v0.32.1
[v0.32.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.31.3...v0.32.0
[v0.31.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.31.2...v0.31.3
[v0.31.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.31.1...v0.31.2
[v0.31.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.31.0...v0.31.1
[v0.31.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.30.2...v0.31.0
[v0.30.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.30.1...v0.30.2
[v0.30.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.30.0...v0.30.1
[v0.30.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.29.1...v0.30.0
[v0.29.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.29.0...v0.29.1
[v0.29.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.28.1...v0.29.0
[v0.28.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.28.0...v0.28.1
[v0.28.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.27.1...v0.28.0
[v0.27.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.27.0...v0.27.1
[v0.27.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.26.0...v0.27.0
[v0.26.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.25.1...v0.26.0
[v0.25.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.25.0...v0.25.1
[v0.25.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.24.1...v0.25.0
[v0.24.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.24.0...v0.24.1
[v0.24.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.23.4...v0.24.0
[v0.23.4]: https://github.com/newrelic/newrelic-client-go/compare/v0.23.3...v0.23.4
[v0.23.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.23.2...v0.23.3
[v0.23.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.23.1...v0.23.2
[v0.23.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.23.0...v0.23.1
[v0.23.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.22.0...v0.23.0
[v0.22.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.21.1...v0.22.0
[v0.21.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.21.0...v0.21.1
[v0.21.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.20.1...v0.21.0
[v0.20.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.20.0...v0.20.1
[v0.20.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.19.0...v0.20.0
[v0.19.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.0...v0.19.0
[v0.18.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.17.1...v0.18.0
[v0.17.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.17.0...v0.17.1
[v0.17.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.16.0...v0.17.0
[v0.16.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.15.0...v0.16.0
[v0.15.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.14.0...v0.15.0
[v0.14.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.13.0...v0.14.0
[v0.13.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.12.0...v0.13.0
[v0.12.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.11.0...v0.12.0
[v0.11.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.10.1...v0.11.0
[v0.10.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.10.0...v0.10.1
[v0.10.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.9.0...v0.10.0
[v0.9.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.8.0...v0.9.0
[v0.8.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.7.1...v0.8.0
[v0.7.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.7.0...v0.7.1
[v0.7.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.6.0...v0.7.0
[v0.6.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.5.1...v0.6.0
[v0.5.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.1.0...v0.2.0
