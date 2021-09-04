package errors

const (
	// Error Common
	CodeValueInvalid = Code(iota + 100)
)

const (
	// Error On Cache
	CodeCacheMarshal = Code(iota + 300)
	CodeCacheUnmarshal
	CodeCacheGetSimpleKey
	CodeCacheSetSimpleKey
	CodeCacheDeleteSimpleKey
	CodeCacheSetExpiration
	CodeCacheDecode
	CodeCacheLockNotAcquired
	CodeCacheLockFailed
	CodeCacheInvalidCastType
	CodeCacheNotFound
)

const (
	// Error On SQL
	CodeSQLBuilder = Code(iota + 700)
	CodeSQLTxBegin
	CodeSQLCreate
	CodeSQLTxCommit
	CodeSQLRead
	CodeSQLRowScan
	CodeSQLCannotRetrieveLastInsertID
)

const (
	// Code HTTP Handler
	CodeHTTPBadRequest = Code(iota + 900)
	CodeHTTPNotFound
	CodeHTTPInternalServerError
	CodeHTTPUnmarshal
)

const (
	// Code Usecase
	CodeUsecase = Code(iota + 1000)
)

const (
	// Code Domain
	CodeDomain = Code(iota + 2000)
)
