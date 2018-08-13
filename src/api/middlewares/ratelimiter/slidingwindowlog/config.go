package slidingwindowlog

type Config struct {
	Prefix     string //Prefix in redis key
	MaxRequest uint   //Number of Request in a period
	Period     string //Period can be second, minute, hour, or day
}

func NewConfig(prefix string, maxRequest uint, period string) *Config {
	config := &Config{
		Prefix:     prefix,
		MaxRequest: maxRequest,
		Period:     period,
	}
	return config
}

func GetPeriodInt(period string) int64 {
	switch period {
	case "day":
		return 86400e9
	case "hour":
		return 3600e9
	case "minute":
		return 60e9
	case "second":
		return 1e9
	}
	return 0
}
