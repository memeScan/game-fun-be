package util

import "time"

type TimeRange struct {
    Start time.Time
    End   time.Time
}

func GetTimeRangeFromString(timeStr string) TimeRange {
    now := time.Now()
    switch timeStr {
    case "m1":
        return TimeRange{Start: now.Add(-1 * time.Minute), End: now}
    case "m5":
        return TimeRange{Start: now.Add(-5 * time.Minute), End: now}
    case "h1":
        return TimeRange{Start: now.Add(-1 * time.Hour), End: now}
    case "h6":
        return TimeRange{Start: now.Add(-6 * time.Hour), End: now}
    case "h24":
        return TimeRange{Start: now.Add(-24 * time.Hour), End: now}
    default:
        return TimeRange{Start: now.Add(-24 * time.Hour), End: now} // 默认返回24小时范围
    }
}

func FormatTime(t time.Time) string {
    return t.Format("2006-01-02 15:04:05")
}
