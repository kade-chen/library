package token

type RevolkToken struct {
	Action       string `json:"action"`        // revoke_token
	Success      bool   `json:"success"`       // true
	AffectedRows int64  `json:"affected_rows"` // 1
}
