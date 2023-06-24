package whats

// Bool stores v in a new bool value and returns a pointer to it.
func BoolP(v bool) *bool { return &v }

// Int32 stores v in a new int32 value and returns a pointer to it.
func Int32P(v int32) *int32 { return &v }

// Int64 stores v in a new int64 value and returns a pointer to it.
func Int64P(v int64) *int64 { return &v }

// Float32 stores v in a new float32 value and returns a pointer to it.
func Float32P(v float32) *float32 { return &v }

// Float64 stores v in a new float64 value and returns a pointer to it.
func Float64P(v float64) *float64 { return &v }

// Uint32 stores v in a new uint32 value and returns a pointer to it.
func Uint32P(v uint32) *uint32 { return &v }

// Uint64 stores v in a new uint64 value and returns a pointer to it.
func Uint64P(v uint64) *uint64 { return &v }

// String stores v in a new string value and returns a pointer to it.
func StringP(v string) *string { return &v }
