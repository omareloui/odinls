package user

type OAuthProvider string

const (
	NilOAuthProvider OAuthProvider = "NO_AUTH_PROVIDER"
	Google           OAuthProvider = "GOOGLE"
	Facebook         OAuthProvider = "FACEBOOK"
	Microsoft        OAuthProvider = "MICROSOFT"
	Github           OAuthProvider = "GITHUB"
)

func (p OAuthProvider) View() string {
	return map[OAuthProvider]string{
		NilOAuthProvider: "No OAuth Provider",
		Google:           "Google",
		Facebook:         "Facebook",
		Microsoft:        "Microsoft",
		Github:           "Github",
	}[p]
}

type RoleEnum uint8

const (
	UnknownAuthority RoleEnum = iota
	NoAuthority
	Moderator
	Admin
	SuperAdmin
)

func (r RoleEnum) String() string {
	arr := [...]string{
		"UNKNOWN_AUTHORITY", "NO_AUTHORITY", "MODERATOR",
		"ADMIN", "SUPER_ADMIN",
	}
	if int(r) >= len(arr) {
		return arr[0]
	}
	return arr[r]
}

func (r RoleEnum) View() string {
	arr := [...]string{
		"UnknownAuthority", "No Authority", "Moderator",
		"Admin", "Super Admin",
	}
	if int(r) >= len(arr) {
		return arr[0]
	}
	return arr[r]
}

func RoleFromString(role string) RoleEnum {
	switch role {
	case "UNKNOWN_AUTHORITY":
		return UnknownAuthority
	case "NO_AUTHORITY":
		return NoAuthority
	case "MODERATOR":
		return Moderator
	case "ADMIN":
		return Admin
	case "SUPER_ADMIN":
		return SuperAdmin
	default:
		return NoAuthority
	}
}

func (r RoleEnum) IsSuperAdmin() bool {
	return r >= SuperAdmin
}

func (r RoleEnum) IsAdmin() bool {
	return r >= Admin
}

func (r RoleEnum) IsModerator() bool {
	return r >= Moderator
}
