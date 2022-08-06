package timeprovider

type TimeProvider interface {
	GetTimeStamp() int64
}
