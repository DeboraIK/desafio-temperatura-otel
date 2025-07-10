package validators

import "regexp"

func IsValidZipCode(cep string) bool {
	if cep == "" {
		return false
	}
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(cep) && len(cep) == 8
}
