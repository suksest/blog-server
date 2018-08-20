package leakybucket

//Bucket represent token bucket
type Bucket struct {
	Prefix         string //Prefix in redis key
	Capacity       uint   //Number of Capacity in a Bucket
	Period         string //Period can be second, minute, hour, or day
	AllowedRequest uint   //AllowedRequest represent number of allowed request in certain period
}

//NewConfig return new Bucket configuration
func NewConfig(prefix string, capacity uint, period string, allowed uint) *Bucket {
	config := &Bucket{
		Prefix:         prefix,
		Capacity:       capacity,
		Period:         period,
		AllowedRequest: allowed,
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
