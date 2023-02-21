package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SafeRE-IT/horizon/db2"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"time"

	"github.com/spf13/cast"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/strkey"
	"gitlab.com/tokend/go/xdr"
	"github.com/SafeRE-IT/horizon/db2/core"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/utils"
)

const (
	// ParamCursor is a query string param name
	ParamCursor = "cursor"
	// ParamOrder is a query string param name
	ParamOrder = "order"
	// ParamLimit is a query string param name
	ParamLimit = "limit"
	// ParamLimit is a query string param name
	ParamPage = "page"
)

// GetString retrieves a string from either the URLParams, form or query string.
// This method uses the priority (URLParams, Form, Query).
func (base *Base) GetString(name string) string {
	if base.Err != nil {
		return ""
	}

	fromURL, ok := base.GojiCtx.URLParams[name]

	if ok {
		return fromURL
	}

	if base.isJson {
		fromJson := base.JsonValue(name)
		if fromJson != "" {
			return fromJson
		}
	} else {
		fromForm := base.R.FormValue(name)

		if fromForm != "" {
			return fromForm
		}
	}

	return base.R.URL.Query().Get(name)
}

func (base *Base) GetIntArray(name string) []int {
	if base.Err != nil {
		return nil
	}

	stringArray := base.R.URL.Query()[name]

	res, err := getIntArrayFromStringArray(stringArray)
	if err != nil {
		base.SetInvalidField(name, err)
		return nil
	}

	return res
}

func getIntArrayFromStringArray(input []string) (result []int, err error) {
	for _, str := range input {
		value, err := strconv.Atoi(str)
		if err != nil {
			return nil, errors.New("failed to convert to int from string")
		}

		result = append(result, value)
	}

	return
}

// GetNonEmptyString retrieves an string from the action parameter of the given name.
// Populates err if the value is an empty string
func (base *Base) GetNonEmptyString(name string) string {
	if base.Err != nil {
		return ""
	}

	value := base.GetString(name)
	if value == "" {
		base.SetInvalidField(name, errors.New("Must not be empty."))
	}

	return value
}

func (base *Base) TryGetTime(name string) *time.Time {
	if base.Err != nil {
		return nil
	}

	value := base.GetString(name)
	if value != "" {
		ts, err := time.Parse(time.RFC3339, value)
		if err != nil {
			base.SetInvalidField(name, err)
		}
		return &ts
	}
	return nil
}

func (base *Base) GetTime(name string) *time.Time {
	if base.Err != nil {
		return nil
	}
	ts := base.TryGetTime(name)
	if ts == nil && base.Err == nil {
		base.SetInvalidField(name, errors.New("Must not be empty."))
	}
	return ts
}

func (base *Base) GetBalanceIDAsString(name string) string {
	if base.Err != nil {
		return ""
	}

	rawValue := base.GetNonEmptyString(name)
	if base.Err != nil {
		return ""
	}

	_, err := strkey.Decode(strkey.VersionByteBalanceID, rawValue)
	if err != nil {
		base.SetInvalidField(name, err)
	}

	return rawValue
}

func (base *Base) GetRestrictedString(name string, minLength, maxLength int) string {
	rawValue := base.GetNonEmptyString(name)
	if base.Err != nil {
		return ""
	}
	length := len(rawValue)
	validLength := (length >= minLength) && (length <= maxLength)
	if !validLength {
		base.SetInvalidField(name, errors.New(fmt.Sprintf(" is not %d-%d characters", minLength, maxLength)))
		return ""
	}

	return rawValue
}

// GetNonEmptyString retrieves an string from the action parameter of the given name.
// Populates err if the value is an empty string
func (base *Base) GetStringWithFlag(name string, nonEmpty bool) string {
	if nonEmpty {
		return base.GetNonEmptyString(name)
	}

	return base.GetString(name)
}

// GetMobileNumber retrieves an string from the action parameter
// of the given name and split whitespaces
// Populates err if the value is an empty string
func (base *Base) GetStingWithoutWhitespaces(name string) string {
	str := base.GetNonEmptyString(name)
	if base.Err != nil {
		return ""
	}

	return strings.Replace(str, " ", "", -1)
}

// GetInt64 retrieves an int64 from the action parameter of the given name.
// Populates err if the value is not a valid int64
func (base *Base) GetInt64(name string) int64 {
	result := base.GetOptionalInt64(name)
	if result == nil {
		return 0
	}

	return *result
}

func (base *Base) GetPositiveInt64(name string) int64 {
	result := base.GetOptionalInt64(name)
	if result == nil {
		base.SetInvalidField(name, errors.New("must not be empty"))
		return 0
	}

	if *result <= 0 {
		base.SetInvalidField(name, errors.New("must be positive"))
	}

	return *result
}

func (base *Base) GetOptionalInt64(name string) *int64 {
	if base.Err != nil {
		return nil
	}

	asStr := base.GetString(name)

	if asStr == "" {
		return nil
	}

	asI64, err := strconv.ParseInt(asStr, 10, 64)

	if err != nil {
		base.SetInvalidField(name, err)
		return nil
	}

	return &asI64
}

func (base *Base) GetOptionalUint64(name string) *uint64 {
	if base.Err != nil {
		return nil
	}

	asStr := base.GetString(name)

	if asStr == "" {
		return nil
	}

	asU64, err := strconv.ParseUint(asStr, 10, 64)

	if err != nil {
		base.SetInvalidField(name, err)
		return nil
	}

	return &asU64
}

// GetInt32 retrieves an int32 from the action parameter of the given name.
// Populates err if the value is not a valid int32
func (base *Base) GetInt32(name string) int32 {
	if base.Err != nil {
		return 0
	}

	asStr := base.GetString(name)

	if asStr == "" {
		return 0
	}

	asI64, err := strconv.ParseInt(asStr, 10, 32)

	if err != nil {
		base.SetInvalidField(name, err)
		return 0
	}

	return int32(asI64)
}

func (base *Base) GetOptionalBool(name string) *bool {
	if base.Err != nil {
		return nil
	}

	asStr := base.GetString(name)
	if asStr == "" {
		return nil
	}

	result, err := strconv.ParseBool(asStr)
	if err != nil {
		base.SetInvalidField(name, err)
		return nil
	}

	return &result
}

// GetInt32 retrieves an int32 from the action parameter of the given name.
// Populates err if the value is not a valid int32
func (base *Base) GetBool(name string) bool {
	return base.GetBoolOrDefault(name, false)
}

// GetBoolOrDefault - returns boolean values passed, if parameter is not available or value is empty string - returns default value
func (base *Base) GetBoolOrDefault(name string, defaultValue bool) bool {
	if base.Err != nil {
		return false
	}

	asStr := base.GetString(name)
	if asStr == "" {
		return defaultValue
	}

	return cast.ToBool(asStr)
}

// GetUInt64 retrieves an uint64 from the action parameter of the given name.
// Populates err if the value is not a valid uint64
func (base *Base) GetUInt64(name string) uint64 {
	if base.Err != nil {
		return 0
	}

	asStr := base.GetString(name)

	if asStr == "" {
		return 0
	}

	asUI64, err := strconv.ParseUint(asStr, 10, 64)

	if err != nil {
		base.SetInvalidField(name, err)
		return 0
	}

	return asUI64
}

// GetPagingParams returns the cursor/order/limit triplet that is the
// standard way of communicating paging data to a horizon endpoint.
func (base *Base) GetPagingParams() (cursor string, order string, limit uint64) {
	if base.Err != nil {
		return
	}

	cursor = base.GetString(ParamCursor)
	order = base.GetString(ParamOrder)
	limit = base.GetUInt64(ParamLimit)

	if lei := base.R.Header.Get("Last-Event-ID"); lei != "" {
		cursor = lei
	}

	return
}

// GetPageQuery is a helper that returns a new db.PageQuery struct initialized
// using the results from a call to GetPagingParams()
func (base *Base) GetPageQuery() db2.PageQuery {
	if base.Err != nil {
		return db2.PageQuery{}
	}

	r, err := db2.NewPageQuery(base.GetPagingParams())

	if err != nil {
		base.Err = err
	}

	return r
}

func (base *Base) GetPagingParamsV2() (page uint64, limit uint64) {
	if base.Err != nil {
		return
	}

	page = base.GetUInt64(ParamPage)
	limit = base.GetUInt64(ParamLimit)

	return
}

// GetAddress retrieves a stellar address.  It confirms the value loaded is a
// valid stellar address, setting an invalid field error if it is not.
func (base *Base) GetAddress(name string) (result string) {
	if base.Err != nil {
		return
	}

	result = base.GetString(name)

	_, err := strkey.Decode(strkey.VersionByteAccountID, result)

	if err != nil {
		base.SetInvalidField(name, err)
	}

	return result
}

// MustAccountID retireves an xdr.AccountID by attempting to decode a stellar
// address at the provided name.
func (base *Base) GetAccountID(name string) (result xdr.AccountId) {
	raw, err := strkey.Decode(strkey.VersionByteAccountID, base.GetString(name))

	if base.Err != nil {
		return
	}

	if err != nil {
		base.SetInvalidField(name, err)
		return
	}

	var key xdr.Uint256
	copy(key[:], raw)

	result, err = xdr.NewAccountId(xdr.CryptoKeyTypeKeyTypeEd25519, key)
	if err != nil {
		base.SetInvalidField(name, err)
		return
	}

	return
}

// GetAmount returns a native amount (i.e. 64-bit integer) by parsing
// the string at the provided name in accordance with the stellar client
// conventions.
func (base *Base) GetAmount(name string) int64 {
	if base.Err != nil {
		return 0
	}

	var err error
	result, err := amount.Parse(base.GetString(name))

	if err != nil {
		base.SetInvalidField(name, err)
		return 0
	}

	return result
}

// GetOptionalAmount returns a pointer of native amount (i.e. 64-bit integer) by parsing
// the string at the provided name in accordance with the stellar client
// conventions.
func (base *Base) GetOptionalAmount(name string) *int64 {
	if base.Err != nil {
		return nil
	}

	str := base.GetString(name)
	if str == "" {
		return nil
	}

	result, err := amount.Parse(str)
	if err != nil {
		base.SetInvalidField(name, err)
		return nil
	}
	return &result
}

// SetInvalidField establishes an error response triggered by an invalid
// input field from the user.
func (base *Base) SetInvalidField(name string, reason error) {
	br := problem.BadRequest

	br.Extras = map[string]interface{}{}
	br.Extras["invalid_field"] = name
	br.Extras["reason"] = reason.Error()

	base.Err = &br
}

// Path returns the current action's path, as determined by the http.Request of
// this action
func (base *Base) Path() string {
	return base.R.URL.Path
}

// ValidateBodyType sets an error on the action if the requests Content-Type
//  is not `application/x-www-form-urlencoded`
func (base *Base) ValidateBodyType() {
	c := base.R.Header.Get("Content-Type")

	if c == "" {
		return
	}

	mt, _, err := mime.ParseMediaType(c)

	if err != nil {
		base.Err = err
		return
	}

	switch {
	case mt == "application/x-www-form-urlencoded":
		return
	case mt == "multipart/form-data":
		return
	case mt == "application/vnd.api+json":
		base.isJson = true
		return
	case mt == "application/json":
		base.isJson = true
		return
	default:
		base.Err = &problem.UnsupportedMediaType
	}
}

func (base *Base) GetCoreAccount(name string, coreQ core.QInterface) *core.Account {
	if base.Err != nil {
		return nil
	}

	accountID := base.GetString(name)
	if accountID == "" {
		return nil
	}

	return base.GetCoreAccountByAddress(accountID, coreQ)
}

func (base *Base) GetCoreAccountByAddress(accountID string, coreQ core.QInterface) *core.Account {
	if base.Err != nil {
		return nil
	}

	account, err := coreQ.Accounts().ByAddress(accountID)
	if err != nil {
		base.Err = &problem.ServerError
		return nil
	}

	if account == nil {
		base.Err = &problem.NotFound
		return nil
	}

	return account
}

func (base *Base) UnmarshalBody(dest interface{}) {
	if !base.isJson {
		base.Err = &problem.UnsupportedMediaType
		return
	}
	decoder := json.NewDecoder(base.R.Body)
	err := decoder.Decode(&dest)
	if err != nil {
		base.Err = &problem.BadRequest
		return
	}
	base.ValidateToProblem(utils.ValidateStruct("", dest))
}

func (base *Base) ValidateToProblem(ok bool, result *utils.ValidateError) {
	if !ok {
		if result != nil {
			base.SetInvalidField(result.Name, result.Reason)
			return
		}
		base.Err = &problem.BadRequest
		return
	}
}

func (base *Base) CalculatePercentFee(percentFee, am, minimalAmount int64) (int64, bool) {
	if percentFee == 0 || am == 0 {
		return 0, false
	}

	return amount.BigDivide(am, percentFee, amount.One*100, amount.ROUND_UP, minimalAmount)
}

func (base *Base) ParseResponse(response *http.Response) (p *problem.P) {
	switch response.StatusCode {
	case http.StatusOK:
		{
			p = &problem.Success
			break
		}
	case http.StatusNotFound:
		{
			p = &problem.NotFound
			break
		}
	case http.StatusUnauthorized:
		{
			p = &problem.NotAllowed
			break
		}
	case http.StatusBadRequest:
		{
			p = &problem.BadRequest
			break
		}
	case http.StatusInternalServerError:
		{
			p = &problem.ServerError
			break
		}
	default:
		p = &problem.P{
			Type:   response.Status,
			Title:  response.Status,
			Status: response.StatusCode,
		}
	}
	return p
}
