package credentials

// Fake is an in-memory Store adapter for tests. It mirrors the keyring
// adapter's semantics: keys are trimmed, blank keys are rejected, a missing
// key reads as ErrNoAPIKey, and every operation fails with
// ErrKeyringUnavailable when Unavailable is set.
type Fake struct {
	Key         string
	Unavailable bool
}

func (f *Fake) Get() (string, error) {
	if f.Unavailable {
		return "", ErrKeyringUnavailable
	}
	if f.Key == "" {
		return "", ErrNoAPIKey
	}
	return f.Key, nil
}

func (f *Fake) Set(key string) error {
	if f.Unavailable {
		return ErrKeyringUnavailable
	}
	key, err := normalizeKey(key)
	if err != nil {
		return err
	}
	f.Key = key
	return nil
}

func (f *Fake) Delete() error {
	if f.Unavailable {
		return ErrKeyringUnavailable
	}
	f.Key = ""
	return nil
}

func (f *Fake) Available() bool {
	return !f.Unavailable
}
