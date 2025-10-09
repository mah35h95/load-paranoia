package auth

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetIdentityToken - returns Googles Identity token
func GetIdentityToken() string {
	out, err := exec.Command("gcloud", "auth", "print-identity-token").Output()
	if err != nil {
		fmt.Println(err)
	}
	return fmt.Sprintf("Bearer %s", strings.ReplaceAll(strings.ReplaceAll(string(out), "\r", ""), "\n", ""))
}

// GetAccessToken - returns Googles Access token
func GetAccessToken() string {
	out, err := exec.Command("gcloud", "auth", "print-access-token").Output()
	if err != nil {
		fmt.Println(err)
	}
	return fmt.Sprintf("Bearer %s", strings.ReplaceAll(strings.ReplaceAll(string(out), "\r", ""), "\n", ""))
}
