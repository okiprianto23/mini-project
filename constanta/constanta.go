package constanta

import "time"

const PathMain = "/main"

const DefaultTimeFormat = "2006-01-02T15:04:05Z"
const DefaultDtoOutTimeFormat = "2006-01-02T15:04:05.999999999Z"
const DateOnlyTimeFormat = "2006-01-02"

const AuthorizationHeaderConstanta = "Authorization"
const ApplicationContextConstanta = "application_context"
const RequestIDConstanta = "X-REQUEST-ID"
const IPAddressConstanta = "X-Forwarded-For"
const SourceConstanta = "X-SOURCE"
const ClientRequestTimestamp = "X-Client-Timestamp"
const LocaleHitAPI = "en-US"
const TimestampSignatureHeaderNameConstanta = "X-Timestamp"
const SignatureHeaderNameConstanta = "X-Signature"
const DeviceHeaderConstanta = "X-Device"
const ResourceHeaderConstanta = "X-Resource"
const RedirectURINameConstanta = "X-Redirect-Uri"
const RedirectMethodNameConstanta = "X-Redirect-Method"

// ==================================LOGGER==================================
// untuk case handle get key logger
const (
	LoggerAccessToken    = "access_token"
	LoggerUserID         = "user_id"
	LoggerVersion        = "version"
	LoggerApplication    = "application"
	LoggerIP             = "ip"
	LoggerPID            = "pid"
	LoggerThread         = "thread"
	LoggerRequestID      = "request_id"
	LoggerSource         = "source"
	LoggerProcessingTime = "processing_time"
	LoggerByteIn         = "byte_in"
	LoggerByteOut        = "byte_out"
	LoggerStatus         = "status"
	LoggerCode           = "code"
	LoggerUrl            = "url"
	LoggerMethod         = "method"
)

// default 1 day
const Default1DayExpired = 24 * time.Hour

// Redis is logout
const INVALID_TOKEN_REDIS_VALUE = "INVALID"

// CONSTANTA
const ID = "ID"
const UserID = "USER_ID"
const Nik = "NIK"
const ConsumerID = "CONSUMER_ID"
