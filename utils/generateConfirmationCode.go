package utils

import (
    "math/rand"
    "time"
    "fmt"
)

func GenerateConfirmationCode() string {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    code := r.Intn(1000000)
    return fmt.Sprintf("%06d", code)
}
