# Changelog

## Unreleased

### Added

* Reviewable request for create, update and remove data 
* Ability to set custom rules and actions in permissions
* Filter by effect in get history endpoint
* Fitler by `created_before`, `created_after` timestamps for reviewable requests
* Filter by `all_tasks`, `all_tasks_any_of`, `all_tasks_not_set` for reviewable requests
* Field `wait_for_result` (bool, default value: `true`) to `POST /v3/transactions` request which allows not waiting for transaction result from core (submission finishes faster but invalid tx handling is sender's responsibility). Not guaranteed that transactions submitted using this flag will be applied successfully.
* Vote creation time
* Participants count in `/v3/sales/{id}`
* Operation to update owner of the data entry and reviewable request for update data. New endpoints: `/v3/data_owner_update_requests` to get a list of requests and `/v3/data_owner_update_requests/{id}` to get a request by id

### Changed

* Disable returning errors when using parameters not supported by endpoint
* Filter parameter `filter[request_details.asset]` on endpoint `/v3/create_withdraw_requests` now accepts slice of asset codes

### Fixed

* `/v3/order_book/{id}` now uses renders numbers with trailing-digits-count-based dot
* Match participant in manage offer
* Error on '/v3/sales/{id}/relationships/participation'
* Duplicated sale entries on '/v3/sales?filter[participant]='

## 3.9.1

### Fixed

* panic on `/v3/accounts/{id}/sales`


## 3.9.0

### Added
 
* Traefik
* Endpoint to get operation `/v3/operations/{id}`

### Fixed 

* Overflow quote amount `/v3/order_books/{base}:{quote}:{order_book_id}`
* Incorrect 404 on `/v3/create_issuance_requests` with filter by receiver
* Fee in unlocked effect 
* Invalid conversion price `/v3/accounts/{id}/converted_balances/{asset_code}`

## 3.8.3

### Added

* Ingest version bumped up
* Creator details in `PaymentReviewableRequest`, `ManageOfferReviewableRequest`

### Fixed 

* Ingest of preissuance request details.

## 3.8.0

### Added

* Endpoint for getting account list `v3/accounts`
* Filter by participant for `/v3/sales`
* Filter by asset codes (`v3/assets`)
* Filter by sales ids (`v3/sales`)

### Fixed

* Parsing array query parameters
* slow `/v3/transactions` get requests (same as 3.5.3)
* incorrect `/v3/asset_pairs` filtering
* panic on `/v3/history`
* empty relationships in create asset request

## 3.7.2

### Fixed

* Can see the unlocking amount transaction when an order is matched on the price less than was locked
* 500 on get history (create manage limits request)
* participants statistics for immediate sale
* empty sale participants effects in admin panel
* participant effects for manage offer and payment requests
* 500 on `/v3/create_withdraw_requests`

## 3.7.1

### Added

* Endpoints for redemption reviewable request

## 3.7.0

### Added  

* Operations endpoint (`/v3/operations`)
* Filter by status for asset endpoint

### Fixed 

* ingest of participant effects for sale
* 500 error on endpoints that include deleted assets
* returning 404 when cannot find tx

### Changed

* increased speed of ingest v2

## 3.6.2

### Added

* Cache for `GET` requests

### Fixed

* Batch inserter corner case

## 3.6.1

### Added

* `/v3/license` endpoint to return current license info

### Removed

* Request for account endpoint (`/v3/accounts/{id}/requests/{request_id}`)

### Fixed

* RequestID in create issuance request operation details
* includes for swaps (`/v3/swaps`)
* docs consistency with actual responses

## 3.6.0

### Added

* `/v3/manage_offer_requests` endpoint
* `/v3/create_payment_requests` endpoint
* `CreateManageOfferRequestOp` ingestion
* `CreatePaymentRequestOp` ingestion
* `CreatePaymentRequestOp` ingestion
* `PaymentRequest` ingestion
* `ManageOfferRequest` ingestion
* `RemoveAssetOp` ingestion
* `OpenSwapOp` ingestion
* `CloseSwapOp` ingestion
* `Swap` ingestion
* `/v3/swaps` endpoint to get filtered list of swaps
* `/v3/swaps/{id}` endpoint to get swap by id

### Fixed

* panic on `v3/balances`
* error on reviewable request ingestion

## 3.5.3

###Fixed

* slow `/v3/transactions` get requests

## 3.5.2

### Fixed

* Key for atomic swap ask quote asset relation
* Panic on `ManagePoll` ingestion

# 3.5.1.2

### Fixed 

* KYC recovery request ingestion

### Added 

* Ingest version

## 3.5.1

### Changed

* deprecate `/v3`

### Added

* filter `owner` for balances list
* endpoint `/v3/accounts/{account-id}/requests/{request-id}`
* `base_asset` to `AtomicSwapAsk` relationships
* `/v3/info`
* `request_details.ask_id` filter for `/v3/create_atomic_swap_bid_requests`
* `request_details.ask_owner` filter for `/v3/create_atomic_swap_bid_requests`

### Fixed

* panic on ingesting reviewable requests
* internal error on `v3/accounts`
* get reviewable requests docs
* atomic swap request filters
* checking signature on `/v3/accounts` if kyc_data included
* checking signature on `/v3/create_atomic_swap_bid_requests` for ask owner

## 3.5.0

### Added

* Adding `details` from `ClosePollOp` to `creator_details` of polls after its closing
* filter `asset_owner` for balances list
* Endpoint `/v3/votes/{voter}` which returns all the _votes_ created by specific _voter_ with relationships `Poll` (_poll_ where the _vote_ was created) and `Account` (_voter_ account)

### Fixed

* Sales not allowed to participate in appearing on`/v3/account/{id}/sales`
* Transaction failure response for `/v3/transactions`
* Waiting for transaction ingestion for `/v3/transactions`
* request details for create atomic swap ask and bid requests
* response on `v3/create_atomic_swap_aks_requests` (quote assets include)
* vote id in `v3/votes` response

## 3.4.0

### Removed

* handling `/atomic_swap_bids` request

### Added
* `GetAtomicSwapAskList` handler to handle `v3/atomic_swap_asks` request
* `GetAtomicSwapAsk` handler to handle `v3/atomic_swap_asks/{id}` request
* `InitiateKYCRecovery` operation ingestion
* `CreateKYCRecoveryRequest` operation ingestion
* `KYCRecovery` reviewable request ingestion
* `GetKYCRecoveryRequests` to handle `/v3/kyc_recovery_requests`, `/v3/kyc_recovery_requests/{id}`
* Docs on kyc recovery
* `kyc_data` relationship to `account` resource type
* `kyc_data` include parameter to `/v3/accounts/{id}`
* `RemoveAssetPairOp` operation ingestion
* Tests for operation details handling
* `/v3/accounts/{id}/sales/{sale_id}` endpoint to get sale by id, if account is allowed to participate
* `GetSaleForAccount` to handle `/v3/accounts/{id}/sales/{sale_id}`
* `Asset` to withdraw request relationships
* Filter by asset for `/v3/create_withdraw_requests` and `/v3/create_withdraw_requests/{id}`
* Ability to include asset in responses`/v3/create_withdraw_requests` and `/v3/create_withdraw_requests/{id}`
* Option for asset issuer to access all history/movements through `/v3/history` or `/v3/movements` using filter by asset
* KYC recovery status tracking

### Fixed

* docs key value u32 value
* add `trailing_digits_count` to asset creation request
* Panic on `/v3/history` handling `ManageAccountSpecificRuleOp` details
* Own sales appearing on `/v3/account/{id}/sales`
* Legacy withdraw requests endpoint asset filter

### Changed

* `CurrentIngestVersion` increased
* Operation Details resource handling
* swap `v3/create_atomic_swap_bid_requests` and `v3/create_atomic_swap_ask_requests` names

## 3.3.2

### Fixed

* Syntax error in `sale_access_definition` migration

## 3.3.1

### Added

* Filter by address to `/v3/sales/{id}/relationships/whitelist`
* `access_definition_type` attribute to sale and create sale request resources

## 3.3.1-x.2

### Fixed

* Sale participation response rendering
* Names of multiple resources responses changed (e.g. `AccountsResponse` -> `AccountListResponse`)
* Missing `base_hard_cap` in sale attributes
* Non-unique identifiers for `sale-quote-assets` resources
* Missing soft and hard caps in sale request details

### Added

* `/v3/matches` endpoint to get the history of secondary market matches
* `asset` relationship to `converted-balances-collections` resource type
* `asset` include parameter to `/v3/accounts/{id}/converted_balances{asset_code}`

### Removed

* Sale participation on sale close ingestion

### Changed

* Sale participations are not taken from `participant_effects` table

## 3.3.1-x.1

### Fixed

* Empty response on `/v3/sale/relationships/participation` when sale state is opened

## 3.3.1-x.0

### Added

* `/v3/accounts/{id}/sales` endpoint to get list of sales account is allowed to participate in
* `/v3/sales/{id}/relationships/whitelist` endpoint to get list of accounts explicitly granted 
permission to participate in sale.
* `/v3/sales/{id}/relationships/participation` endpoint to get list of sale participations
* `/v3/account/{id}/converted_balances/{asset_code}` endpoint to get collection of converted balances
* Sale participation on sale close ingestion
* Account specific rules ingestion
* Sale version ingestion
* Bumped up ingest version
* `limits_with_stats` to `account` resource

## 3.3.0

### Added

* Allow forcing reingest on application start up (see ingest2/main.go)
* `/v3/limits` endpoint to get limits list with filters
* Squash matches with the same price in history
* `/v3/movements` endpoint to get participants effect related to balance movements
* `/v3` endpoint to get horizon info
* Tx submission system v2 with retry
* `/v3/transactions` to submit txs
* Cancel poll and update poll end time ingestion

## Fixed
* `/v3/account` to include all limits imposed for account
* Issue with `/v3/assets` for assets with max_uint64 fields
* Issues with balances and limits having fields types for amount int64 instead of uint64. 
* Update doorman with SignerOf check fixed
* Problem with CPU usage on OS X
* Fixed `null` messages on failed transaction submit

## 3.2.0-x.9

## Fixed

* `/v3/transactions/` includes were never populated

## 3.2.0-x.8

## Fixed

* `/v3/sales` was setting the wrong caps if default quote asset was also in the list of quote assets

## 3.2.0-x.7

### Fixed

* `/v3/votes` ignoring page params and always returning all the votes

## 3.2.0-x.6

### Added

* `/v3/transactions/{id}` endpoint to get transaction by hash or ID

## 3.2.0-x.5

### Added

* `/v3/balances/{id}` endpoint
* docs for `/v3/balances/{id}` endpoint

### Fixed

* 500 for `GET /balances`  because of nil pointer on balance relationships
* Default quote asset include for `v3/sales`

## 3.2.0-x.4

### Fixed

* 500 for `GET /accounts?include=balances` because of nil pointer on balance relationships
* Panics on manage vote/manage poll/manage poll creation requests operation details

## 3.2.0-x.0

### Removed

* dependency from regources/v2

### Fixed

* saving equal values to transaction's `result_meta_xdr` and `result_xdr`
* all docs for proper regources generation
* `ManageAssetOpAttributes.PreissuanceSigner` -> `ManageAssetOpAttributes.PreIssuanceSigner`
* json: `preissued_signer` -> `pre_issuance_signer`
* missing poll id in create poll request response
* order book entries sorting by price in alphabetical order

### Added

* generated regources/generated in vendor
* new endpoint `/v3/order_books/{id}`

### Changed

* old endpoint `/v3/order_book/{id}` marked as deprecated

## 3.1.1

### Fixed

* psql error on poll ingesting

## 3.1.0-rc.0

### Added

* `GetTransactionList` handler to handler `v3/transactions` endpoint
* Messages for new error codes
* `GetBalanceList` handler to handle `v3/balances` request
* `GetPublicKeyEntry` handler to handle `v3/public_key_entries/{id}` request
* Receiver filter for create issuance requests
* Polls and votes ingestion
* `GetPollList` to handle `/v3/polls`
* `GetPoll` to handle `/v3/polls/{id}`
* `GetVoteList` to handle `/v3/polls/{id}/relationships/votes`
* `GetVote` to handle `/v3/polls/{id}/relationships/votes{voter}`
* `GetCreatePollRequests` to handle `/v3/create_poll_requests`, `/v3/create_poll_requests/{id}`
*  Docs on polls

### Fixed

* panic on handle set fee operation when account or account role does not exist
* panic on actions' `isAllowed` checking
* error on get reviewable request by reviewer
* order book is now sorted
* create change role ingestion causing 500 on operations
* Now Bad Request is correctly returned in case of invalid signature

### Removed

* charts, endpoint for charts etc.
* influx

## 3.0.1-x.21
## Fixed
* Invalid fees joined for account

## 3.0.1-x.20
## Fixed
* Not valid percent amount for manage offer op

## 3.0.1-x.19
## Added
* Error codes for manageSigner and skiped op
## Fixed
* Information disclosure without permission for reviewable requests v1

## 3.0.1-x.18
## Fixed

* Ingest v2 stop on sale participation delete

## 3.0.1-x.17
* Fixed issue with 0 sale state

## 3.0.1-x.16

* Added license op and stamp op to history

## 3.0.1-x.15

### Fixed
* Ingest v2 fail on cancel sale request

## 3.0.1-x.14

### Added

* Asset type to asset and create asset response
* Ingest of details of LicenseOp and StampOp
* Error for v3 history, when trying to request filtration by non existing account/balance

### Changed

* xdr revision
* fix limits select
* Messages for error codes

### Fixed
* Change role request creator details null
* Fixed state of the reviewable requests
* Added xdr type to reviewable request
* Fees list when account type filter is not requested
* 401 on v1 operations when skipSigCheck is on
* Ingest v2 stop on check_sale_state op


# 3.0.1-x.13

### Fixed

* Returning account and signer rule action using xdr type in account and signer rules responses
* (internal) proper types for reviewable requests amounts
* Removed KYCData from ChangeRoleRequest

# 3.0.1-x.11

### Fixed

* 500 on `/history` when there are update signer role operations
* pagination in docs
* all endpoints in docs has padlocks and 401 errcode where needed
* required & non-required fields

# 3.0.1-x.10

### Added

* Added xdr revision to root response

### Fixed
* Fixed issue with empty filter present in url params handled as requested filter

# 3.0.1-x.9

### Changed

* Rename `Details`/`Reason`/`NewDetails` field names in operation which create reqeusts unified to `creatorDetails`
* Rename isForbid to forbids
* Fixed issues with fee
* Fixed issues with filters by flags or dest account for change role requests in v1

### Fixed

* Error codes
* `request_details` format to satisfy JSON API spec

# 3.0.1-x.7

### Added

* Signers endpoint
* Account/Signer Role/Rules
* Endpoints for reviewable requests
* Docs on reviewable requests and sale

### Changed

* `Details`/`Reason`/`NewDetails` field names in reviewable request types unified to `creatorDetails`
* Corresponding messages in `messages` map
* Calculated fee response

### Removed

* `UpdateSaleEndTimeRequest` type

### Fixed

* Fee bounds
* Fee calculation for account

# 3.0.1-x.6

### Fixed
* 500 on fees

# 3.0.1-x.5

### Added

* Helper "test" for quick transaction envelopes unmarshal

### Fixed

* 500 error on /history endpoint when receiving history of payments
* Fixed issue with balanceID not been set for balance effects
* Fixed Fee 500 if asset does not exists


# 3.0.1-x.4

## Fixed

* Switched endpoint version to only track major version
* Switched config back to default

# 3.0.1-x.3

## Added

Added back proxy to API

## Fixed

* Allow to specify filter for primary market offers via orderBookID = -1
* (internal) Fixed ingest v2. Participant effect has not included asset.
* (internal) Fixed ingest v2. Correctly handle `fulfilled` on review of request

# 3.0.1-x.2

## Added

* added handler for `POST /v3.0/transactions`

## Changed

* (internal) Janus config
* `v2`-prefixed endpoints updated to `v3.0`-prefixed

## Fixed

* Fixed rendering success operation state when corresponding request is rejected

# 3.0.1-x.1

## Fixed

* (internal) Fixed panic on ingest_v2 create account op trying to get referrer accountID which might not exist
* (internal) Fixed nil pointer exception on ingest_v2 on withdrawal participant effect handling

# 3.0.1-x.0

## Changed

* Updated XDR

# 3.0.0-x.2

## Added

* `/sales` endpoint
* `/order_book` endpoint

## Fixed

* changelog format

# 3.0.0-x.1

* `assets` owner is not rendering as `null` anymore
