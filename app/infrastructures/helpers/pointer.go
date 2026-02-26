package helpers

// Helper function untuk pointer uint
func IntPtr(u int) *int {
	return &u
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
