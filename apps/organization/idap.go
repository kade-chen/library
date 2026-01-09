package domain

// NewDefaultConfig represents the default LDAP config.
func NewDefaultLDAPConfig() *LdapConfig {
	return &LdapConfig{
		MailAttribute:        "mail",
		DisplayNameAttribute: "displayName",
		GroupnameAttribute:   "cn",
		UsernameAttribute:    "uid",
		UserFilter:           "(uid={input})",
		GroupFilter:          "(|(member={dn})(uid={username})(uid={input}))",
	}
}
