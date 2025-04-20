package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const alphabets = "abcdefghijklmnopqrstuvwxyx"

// RandomInt generates a random int b/w max and min
func RandomInt(max, min int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates random string of n characters
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// generates a random movie title
func RandomTitle() string {
	return fmt.Sprintf("Movie %s", RandomString(5))
}

// generates a random movie description
func RandomDescription() string {
	return fmt.Sprintf("This is a description about %s.", RandomString(10))
}

// generates a random poster_url
func RandomPosterURL() string {
	return fmt.Sprintf("https://image.kratos69.org/t/p/w500/%s.jpg", RandomString(10))
}

// generates a random genre (IDs 1–10)
func RandomGenreID() int32 {
	return int32(rand.Intn(10) + 1)
}

// generates a random role
func RandomRole() string {
	roles := []string{"admin", "customer"}
	return roles[rand.Intn(len(roles))]
}

func RandomFutureTime() time.Time {
	return time.Now().Add(time.Duration(rand.Intn(1000)) * time.Minute)
}

// generates a random row
func RandomRow() int32 {
	return int32(rand.Intn(5) + 1) // 1 to 5
}

// generates a random seat number
func RandomSeatNumber() int32 {
	return int32(rand.Intn(10) + 1) // 1 to 10
}

// Returns a random price as pgtype.Numeric
func RandomPrice() pgtype.Numeric {
	price := float64(rand.Intn(2000)+500) / 100.0 // e.g., 5.00 to 25.00
	priceStr := strconv.FormatFloat(price, 'f', 2, 64)

	var numeric pgtype.Numeric
	err := numeric.Scan(priceStr)
	if err != nil {
		panic("failed to scan random price into pgtype.Numeric")
	}

	return numeric
}
