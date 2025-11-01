package rule

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/EgorO1999/go_final_project/pkg/db"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart format: %w", err)
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid 'd' format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid day interval")
		}
		for {
			startDate = startDate.AddDate(0, 0, days)
			if afterNow(startDate, now) {
				return startDate.Format("20060102"), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid format: year doesn't need an interval")
		}
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if afterNow(startDate, now) {
				return startDate.Format("20060102"), nil
			}
		}
	case "w":
		if len(parts) != 2 {
			return "", errors.New("invalid 'w' format")
		}

		days := strings.Split(parts[1], ",")

		var weekdays [8]bool

		for _, day := range days {
			d, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}
			if d >= 1 && d <= 7 {
				weekdays[d] = true
			} else {
				return "", errors.New("invalid week days format for 'w'")
			}
		}

		nextDate := startDate

		for {
			weekday := int(nextDate.Weekday())
			if weekday == 0 {
				weekday = 7
			}

			if weekdays[weekday] && afterNow(nextDate, now) {
				break
			}
			nextDate = nextDate.AddDate(0, 0, 1)
		}
		return nextDate.Format("20060102"), nil
	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid 'm' format")
		}
		days := strings.Split(parts[1], ",")
		var months []string
		if len(parts) == 3 {
			months = strings.Split(parts[2], ",")
		}

		var lastDay bool = false
		var lastButOneDay bool = false

		var day [32]bool
		var month [13]bool

		for _, d := range days {
			dayNum, err := strconv.Atoi(d)
			if err != nil {
				return "", err
			}

			if dayNum < -2 || dayNum == 0 || dayNum > 31 {
				return "", errors.New("invalid day for 'm'")
			}

			if dayNum > 0 {
				day[dayNum] = true
			} else if dayNum == -1 {
				lastDay = true
			} else {
				lastButOneDay = true
			}

		}

		if len(months) > 0 {
			for _, m := range months {
				monthNum, err := strconv.Atoi(m)
				if err != nil {
					return "", err
				}
				if monthNum < 1 || monthNum > 12 {
					return "", errors.New("invalid month for 'm'")
				}
				month[monthNum] = true
			}
		}

		nextDate := startDate

		if len(parts) == 2 {
			for i, _ := range month {
				month[i] = true
			}
		}
		counter := 0

		for {
			d := int(nextDate.Day())
			m := int(nextDate.Month())
			y := int(nextDate.Year())

			lastDayOfMonth := time.Date(y, time.Month(m)+1, 0, 0, 0, 0, 0, nextDate.Location()).Day()

			if counter > 1462 {
				return "", errors.New("invalid 'm' format")
			}

			if lastDay && d == lastDayOfMonth {
				break
			}
			if lastButOneDay && d == lastDayOfMonth-1 {
				break
			}

			if day[d] && month[m] && afterNow(nextDate, now) {
				break
			}
			nextDate = nextDate.AddDate(0, 0, 1)
			counter++
		}
		return nextDate.Format("20060102"), nil
	default:
		return "", errors.New("unsupported repeat rule")
	}
}

func CheckDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("incorrect date format")
	}

	var nextDate string

	if task.Repeat != "" {
		nextDate, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("invalid repeat rule format")
		}
	}

	if now.Format("20060102") == t.Format("20060102") {
		return nil
	}

	if !afterNow(t, now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = nextDate
		}
	}

	return nil
}

func afterNow(t, now time.Time) bool {
	ty, tm, td := t.Date()
	ny, nm, nd := now.Date()

	if ty > ny {
		return true
	}
	if ty == ny {
		if tm > nm {
			return true
		}
		if tm == nm {
			return td > nd
		}
	}
	return false
}
