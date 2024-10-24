package utils

import "regexp"

func ValidateCredential(credentialType, credential string) (bool, string) {
	if credentialType == "email" {
		emailRegex := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
		re := regexp.MustCompile(emailRegex)

		if re.MatchString(credential) {
			return true, "Valid email."
		} else {
			return false, "Invalid email format."
		}
	} else if credentialType == "phone" {
		phoneRegex := `^\+?[1-9]\d{1,14}$`
		re := regexp.MustCompile(phoneRegex)

		if re.MatchString(credential) {
			return true, "Valid phone number."
		} else {
			return false, "Invalid phone number format."
		}
	} else {
		return false, "Unknown type. Please use 'email' or 'phone'."
	}
}
