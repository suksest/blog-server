package tokenbucket

import "time"

//Bucket represent token bucket
type Bucket struct {
	Prefix          string //Prefix in redis key
	Capacity        uint   //Number of Capacity in a Bucket
	AvailableTokens uint   //Number of AvailableToken in the Bucket === deleted soon
	Period          string //Period can be second, minute, hour, or day
	StartTimestamp  int64  //Last refill time in Unix time format === change to StartTimestamp
}

//NewConfig return new Bucket configuration
func NewConfig(prefix string, capacity uint, period string) *Bucket {
	config := &Bucket{
		Prefix:          prefix,
		Capacity:        capacity,
		AvailableTokens: capacity,
		Period:          period,
		StartTimestamp:  time.Now().Unix(),
	}
	return config
}

//GetPeriodInt return duration of period in second
func GetPeriodInt(period string) int64 {
	switch period {
	case "day":
		return 86400
	case "hour":
		return 3600
	case "minute":
		return 60
	case "second":
		return 1
	}
	return 0
}
