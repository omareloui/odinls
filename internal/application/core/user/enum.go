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
	NoAuthority RoleEnum = iota
	Moderator
	Admin
	SuperAdmin
)

func (r RoleEnum) String() string {
	return [...]string{
		"NO_AUTHORITY", "MODERATOR",
		"ADMIN", "SUPER_ADMIN",
	}[r]
}

func (r RoleEnum) View() string {
	return [...]string{
		"No Authority", "Moderator",
		"Admin", "Super Admin",
	}[r]
}

func RoleFromString(role string) RoleEnum {
	switch role {
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
