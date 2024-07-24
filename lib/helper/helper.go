package helper

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	perm []int
	mu   sync.Mutex
)

func GenerateSnowflakeID() (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	if len(perm) == 0 {
		perm = rand.Perm(1023)

		for i := range perm {
			perm[i]++
		}
	}

	nodeNum := perm[len(perm)-1]
	perm = perm[:len(perm)-1]

	node, err := snowflake.NewNode(int64(nodeNum))
	if err != nil {
		return 0, err
	}

	// Sleep for 1 milliseond to prevent generating the same ID.
	time.Sleep(1 * time.Millisecond)

	// Generate a snowflake ID.
	id := node.Generate()

	return id.Int64(), nil
}

func GenerateRandomDigits(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = byte(rand.Intn(10) + 48)
	}

	return string(code)
}

func GenerateRandomString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456790")

	code := make([]rune, length)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}

	return string(code)
}

func GenerateSlug(title string) string {
	// Convert title to lowercase
	title = strings.ToLower(title)

	// Replace spaces with hyphens
	title = strings.ReplaceAll(title, " ", "-")

	// Remove special characters
	title = removeSpecialChars(title)

	return title
}

func removeSpecialChars(title string) string {
	// Define special characters to remove
	specialChars := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+", "=", "{", "}", "[", "]", "|", "\\", ":", ";", "\"", "'", "<", ">", ",", ".", "?", "/"}

	// Remove special characters from the title
	for _, char := range specialChars {
		title = strings.ReplaceAll(title, char, "")
	}

	return title
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ShortCode(name string, numLetters int) (string, error) {
	if len(name) == 0 {
		return "", errors.New("empty name provided")
	}

	if numLetters <= 0 || numLetters > len(name) {
		return "", errors.New("invalid number of letters")
	}

	return name[:numLetters], nil
}

func GenerateTimestamp() string {
	return time.Now().Format("20060102150405")
}

func StringToPointer(s string) *string {
	return &s
}

func GetLastNMonths(n int) []string {
	var months []string

	for i := 0; i < n; i++ {
		months = append(months, time.Now().AddDate(0, -i, 0).Format("2006-01"))
	}

	return months
}

func GetLastDayOfMonth(year int, month time.Month) time.Time {
	return time.
		Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).
		AddDate(0, 1, 0).
		Add(-time.Second)
}

func GetLastNDays(n int) []string {
	var days []string

	for i := 0; i < n; i++ {
		days = append(days, time.Now().AddDate(0, 0, -i).Format("2006-01-02"))
	}

	return days
}

func GetLastNWeeks(n int) []string {
	var weeks []string

	for i := 0; i < n; i++ {
		weeks = append(weeks, time.Now().AddDate(0, 0, -i*7).Format("2006-01-02"))
	}

	return weeks
}

func GetLastNYears(n int) []string {
	var years []string

	for i := 0; i < n; i++ {
		years = append(years, time.Now().AddDate(-i, 0, 0).Format("2006"))
	}

	return years
}
