package utils

import "webauthn_api/internal/domain"

func ContainsChannel(user domain.UserModel, channel domain.Channel) int {

	for k, a := range user.Subscribtion {
		if isSub := a; isSub.Id == channel.Id && isSub.OwnerId == channel.OwnerId {
			return k
		}
	}
	return -1
}

func HasRole(user domain.UserModel, role domain.Role) int {

	for k, a := range user.Role {
		if isSub := a; isSub.Id == role.Id && role.Name == isSub.Name {
			return k
		}
	}
	return -1
}
