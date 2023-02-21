package requests

import (
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/kit/pgdb"

	"gitlab.com/tokend/go/amount"

	"gitlab.com/distributed_lab/figure"

	"fmt"

	"math"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

const (
	pageParamLimit  = "page[limit]"
	pageParamNumber = "page[number]"
	pageParamCursor = "page[cursor]"
	pageParamOrder  = "page[order]"
)

const defaultLimit uint64 = 15
const maxLimit uint64 = 100

//base - base struct describing params provided by user for the request handler
type base struct {
	include           map[string]struct{}
	supportedIncludes map[string]struct{}
	filter            map[string]string

	queryValues *url.Values
	request     *http.Request
}

type baseOpts struct {
	supportedIncludes map[string]struct{}
	supportedFilters  map[string]struct{}
}

// newRequest - creates new instance of request
func newBase(r *http.Request, opts baseOpts) (*base, error) {
	request := base{
		request:           r,
		supportedIncludes: opts.supportedIncludes,
	}

	err := request.unmarshalQuery(opts)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *base) unmarshalQuery(opts baseOpts) error {
	r.queryValues = new(url.Values)
	*r.queryValues = r.request.URL.Query()

	var err error
	r.filter, err = r.getFilters(opts.supportedFilters)
	if err != nil {
		return err
	}

	r.include, err = r.getIncludes(opts.supportedIncludes)
	if err != nil {
		return err
	}

	return nil
}

func (r *base) URL() *url.URL {
	return r.request.URL
}

func (r *base) marshalQuery() string {
	var builder strings.Builder

	for key, values := range *r.queryValues {
		if !strings.Contains(key, "page[") {
			builder.WriteString(key + "=" + strings.Join(values, ",") + "&")
		}
	}

	return strings.TrimSuffix(builder.String(), "&")
}

var hooks = figure.Hooks{
	"regources.Amount": func(value interface{}) (reflect.Value, error) {
		strVal, ok := value.(string)
		if !ok {
			return reflect.Value{}, errors.New("Failed to parse value as string")
		}

		intVal, err := amount.Parse(strVal)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to parse value as int64")
		}

		result := regources.Amount(intVal)

		return reflect.ValueOf(result), nil
	},
	"[]string": func(value interface{}) (reflect.Value, error) {
		strVal, ok := value.(string)
		if !ok {
			return reflect.Value{}, errors.New("Failed to parse value as string")
		}

		slice := strings.Split(strVal, ",")

		return reflect.ValueOf(slice), nil
	},
	"[]uint64": func(value interface{}) (reflect.Value, error) {
		strVal, ok := value.(string)
		if !ok {
			return reflect.Value{}, errors.New("Failed to parse value as string")
		}

		slice := strings.Split(strVal, ",")

		result := make([]uint64, 0, len(slice))
		for _, val := range slice {
			intVal, err := cast.ToUint64E(val)
			if err != nil {
				return reflect.Value{}, err
			}
			result = append(result, intVal)
		}

		return reflect.ValueOf(result), nil
	},
}

func mkJsonTag(fieldName string) string {
	return fmt.Sprintf("json:\"%s\"", fieldName)
}

func toSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}

func (r *base) populateFilters(target interface{}) error {
	//_ = r.validateFilters(target)
	filter := make(map[string]interface{})
	for k, v := range r.filter {
		filter[k] = v
	}

	err := figure.Out(target).With(figure.BaseHooks, hooks).From(filter).Please()
	if err != nil {
		f := errors.GetFields(err)
		return validation.Errors{
			toSnakeCase(cast.ToString(f["field"])): err,
		}
	}

	return nil
}

func (r *base) ShouldFilter(name string) bool {
	_, ok := r.filter[name]
	return ok
}

//ShouldInclude - returns true if user requested to include resource
// panics if tries to check not supported include
func (r *base) ShouldInclude(name string) bool {
	if _, ok := r.supportedIncludes[name]; !ok {
		panic(errors.From(errors.New("unexpected include check of the request"), logan.F{
			"supported_includes": getSliceOfSupportedIncludes(r.supportedIncludes),
			"checking_include":   name,
		}))
	}
	_, ok := r.include[name]
	return ok
}

//ShouldIncludeAny - returns true if user requested to include one of the provided resources
func (r *base) ShouldIncludeAny(names ...string) bool {
	for _, name := range names {
		if r.ShouldInclude(name) {
			return true
		}
	}

	return false
}

// getString - tries to get string from URL param, if empty gets from query values
func (r *base) getString(name string) string {
	result := chi.URLParam(r.request, name)
	if result != "" {
		return strings.TrimSpace(result)
	}

	return strings.TrimSpace(r.queryValues.Get(name))
}

func (r *base) getUint64(name string) (uint64, error) {
	strVal := r.getString(name)
	if strVal == "" {
		return 0, nil
	}

	return strconv.ParseUint(strVal, 0, 64)
}

func (r *base) getUint64ID() (uint64, error) {
	result, err := r.getUint64("id")
	if err != nil {
		return 0, validation.Errors{
			"id": err,
		}
	}

	return result, nil
}

func (r *base) getInt64(name string) (int64, error) {
	strVal := r.getString(name)
	if strVal == "" {
		return 0, nil
	}

	return strconv.ParseInt(strVal, 0, 64)
}

func (r *base) getInt32(name string) (int32, error) {
	strVal := r.getString(name)
	if strVal == "" {
		return 0, nil
	}

	raw, err := strconv.ParseInt(strVal, 0, 32)
	if err != nil {
		return 0, errors.Wrap(err, "overflow during int32 parsing")
	}

	return int32(raw), nil
}

func (r *base) getFilters(supportedFilters map[string]struct{}) (map[string]string, error) {
	filters := make(map[string]string)
	for queryParam, values := range *r.queryValues {
		if strings.Contains(queryParam, "filter") {
			filterKey := strings.TrimPrefix(queryParam, "filter[")
			filterKey = strings.TrimSuffix(filterKey, "]")
			if _, supported := supportedFilters[filterKey]; !supported {
				return nil, validation.Errors{
					queryParam: errors.New(
						fmt.Sprintf("filter is not supported; supported values: %v",
							getSliceOfSupportedIncludes(supportedFilters)),
					),
				}
			}

			if len(values) == 0 {
				continue
			}

			if len(values) > 1 {
				return nil, validation.Errors{
					queryParam: errors.New("multiple values per one filter are not supported"),
				}
			}

			if values[0] == "" {
				continue
			}

			filters[filterKey] = values[0]
		}
	}

	return filters, nil
}

func (r *base) getIncludes(supportedIncludes map[string]struct{}) (map[string]struct{}, error) {
	const fieldName = "include"
	rawIncludes := r.getString(fieldName)
	if rawIncludes == "" {
		return nil, nil
	}
	includes := strings.Split(rawIncludes, ",")
	requestIncludes := map[string]struct{}{}
	for _, include := range includes {
		if _, supported := supportedIncludes[include]; !supported {
			return nil, validation.Errors{
				fieldName: errors.New(fmt.Sprintf("`%s` is not supported; supported values`: %v", include,
					getSliceOfSupportedIncludes(supportedIncludes))),
			}
		}

		// Note: Because compound documents require full linkage (except when relationship linkage is excluded by sparse fieldsets),
		// intermediate resources in a multi-part path must be returned along with the leaf nodes. For example,
		// a response to a request for comments.author should include comments as well as the author of each of those comments.
		subIncludes := strings.Split(include, ".")
		parentInclude := ""
		for i := range subIncludes {
			if parentInclude == "" {
				parentInclude = subIncludes[i]
			} else {
				parentInclude += "." + subIncludes[i]
			}
			requestIncludes[parentInclude] = struct{}{}
		}
	}

	return requestIncludes, nil
}

func (r *base) getOffsetBasedPageParams() (*pgdb.OffsetPageParams, error) {
	limit, err := r.getLimit(defaultLimit, maxLimit)
	if err != nil {
		return nil, err
	}

	pageNumber, err := r.getPageNumber()
	if err != nil {
		return nil, err
	}

	order, err := r.getOrder()
	if err != nil {
		return nil, err
	}

	return &pgdb.OffsetPageParams{
		Order:      order,
		Limit:      limit,
		PageNumber: pageNumber,
	}, nil
}

func (r *base) getCursorBasedPageParams() (*pgdb.CursorPageParams, error) {
	limit, err := r.getLimit(defaultLimit, maxLimit)
	if err != nil {
		return nil, err
	}

	order, err := r.getOrder()
	if err != nil {
		return nil, err
	}

	cursor, err := r.getCursor()
	if err != nil {
		return nil, err
	}

	if order == pgdb.OrderTypeDesc && cursor == 0 {
		cursor = math.MaxInt64
	}

	return &pgdb.CursorPageParams{
		Limit:  limit,
		Order:  order,
		Cursor: cursor,
	}, nil
}
func Invert(o string) string {

	switch o {
	case pgdb.OrderTypeDesc:
		return pgdb.OrderTypeAsc
	case pgdb.OrderTypeAsc:
		return pgdb.OrderTypeDesc
	default:
		panic(errors.From(errors.New("unexpected order type"), logan.F{
			"order_type": o,
		}))
	}

}

//GetCursorLinks - returns links for cursor based page params
func (r *base) GetCursorLinks(p pgdb.CursorPageParams, last string) *regources.Links {
	result := regources.Links{
		Self: r.getCursorLink(p.Cursor, p.Limit, p.Order),
		Prev: r.getCursorLink(p.Cursor, p.Limit, Invert(p.Order)),
	}

	if last != "" {
		lastI, err := strconv.ParseUint(last, 10, 64)
		if err != nil {
			panic(errors.Wrap(err, "failed to parse cursor", logan.F{
				"last": last,
			}))
		}
		result.Next = r.getCursorLink(lastI, p.Limit, p.Order)
	}

	return &result
}

//GetOffsetLinks - returns links for offset based page params
func (r *base) GetOffsetLinks(p pgdb.OffsetPageParams) *regources.Links {
	result := regources.Links{
		Next: r.getOffsetLink(p.PageNumber+1, p.Limit, p.Order),
		Self: r.getOffsetLink(p.PageNumber, p.Limit, p.Order),
	}

	return &result
}

func (r *base) getCursorLink(cursor, limit uint64, order string) string {
	u := r.URL()
	query := u.Query()
	query.Set(pageParamCursor, strconv.FormatUint(cursor, 10))
	query.Set(pageParamLimit, strconv.FormatUint(limit, 10))
	query.Set(pageParamOrder, order)
	u.RawQuery = query.Encode()
	return u.String()
}

func (r *base) getOffsetLink(pageNumber, limit uint64, order string) string {
	u := r.URL()
	query := u.Query()
	query.Set(pageParamNumber, strconv.FormatUint(pageNumber, 10))
	query.Set(pageParamLimit, strconv.FormatUint(limit, 10))
	query.Set(pageParamOrder, order)
	u.RawQuery = query.Encode()
	return u.String()
}

func getSliceOfSupportedIncludes(includes map[string]struct{}) []string {
	result := make([]string, 0, len(includes))
	for include := range includes {
		result = append(result, include)
	}

	return result
}

func (r *base) getLimit(defaultLimit, maxLimit uint64) (uint64, error) {
	result, err := r.getUint64(pageParamLimit)
	if err != nil {
		return 0, validation.Errors{
			pageParamLimit: errors.New("Must be a valid uint64 value"),
		}
	}

	if result == 0 {
		return defaultLimit, nil
	}

	if result > maxLimit {
		return 0, validation.Errors{
			pageParamLimit: fmt.Errorf("limit must not exceed %d", maxLimit),
		}
	}

	return result, nil
}

func (r *base) getOrder() (string, error) {

	order := r.getString(pageParamOrder)
	switch order {
	case pgdb.OrderTypeAsc, pgdb.OrderTypeDesc:
		return order, nil
	case "":
		return pgdb.OrderTypeAsc, nil
	default:
		return pgdb.OrderTypeDesc, validation.Errors{
			pageParamOrder: fmt.Errorf("allowed order types: %s, %s", pgdb.OrderTypeAsc, pgdb.OrderTypeDesc),
		}
	}
}

func (r *base) getPageNumber() (uint64, error) {
	result, err := r.getUint64(pageParamNumber)
	if err != nil {
		return 0, validation.Errors{
			pageParamNumber: err,
		}
	}

	return result, nil
}

func (r *base) getCursor() (uint64, error) {
	result, err := r.getUint64(pageParamCursor)
	if err != nil {
		return 0, validation.Errors{
			pageParamCursor: err,
		}
	}

	if result > math.MaxInt64 {
		return 0, validation.Errors{
			pageParamCursor: fmt.Errorf("cursor %d exceed max allowed %d", result, math.MaxInt64),
		}
	}

	return result, nil
}
func SetDefaultCursorPageParams(pageParams *pgdb.CursorPageParams) error {
	switch pageParams.Order {
	case pgdb.OrderTypeAsc, pgdb.OrderTypeDesc:
	case "":
		pageParams.Order = pgdb.OrderTypeAsc
	default:
		pageParams.Order = pgdb.OrderTypeDesc
		return validation.Errors{
			pageParamOrder: fmt.Errorf("allowed order types: %s, %s", pgdb.OrderTypeAsc, pgdb.OrderTypeDesc),
		}
	}

	if pageParams.Limit == 0 {
		pageParams.Limit = defaultLimit
	}
	if pageParams.Limit > maxLimit {
		pageParams.Limit = 0
		return validation.Errors{
			pageParamLimit: fmt.Errorf("limit must not exceed %d", maxLimit),
		}

	}

	if pageParams.Order == pgdb.OrderTypeDesc && pageParams.Cursor == 0 {
		pageParams.Cursor = math.MaxInt64
	}

	return nil
}
func (r *base) SetDefaultOffsetPageParams(pageParams *pgdb.OffsetPageParams) error {
	switch pageParams.Order {
	case pgdb.OrderTypeDesc:
		var containsOrder = false
		for queryParam, _ := range *r.queryValues {
			if strings.Contains(queryParam, "page") {
				filterKey := strings.TrimPrefix(queryParam, "page[")
				filterKey = strings.TrimSuffix(filterKey, "]")
				if filterKey == "order" {
					containsOrder = true
					break
				}
			}
		}
		if !containsOrder {
			pageParams.Order = pgdb.OrderTypeAsc
		}
	case pgdb.OrderTypeAsc:
	default:
		pageParams.Order = pgdb.OrderTypeDesc
		return validation.Errors{
			pageParamOrder: fmt.Errorf("allowed order types: %s, %s", pgdb.OrderTypeAsc, pgdb.OrderTypeDesc),
		}
	}

	if pageParams.Limit > maxLimit {
		pageParams.Limit = 0
		return validation.Errors{
			pageParamLimit: fmt.Errorf("limit must not exceed %d", maxLimit),
		}
	}

	return nil
}
