package storage

func (s *Storage) MarkStoresCorrupted() error {
	return nil
}

func (s *Storage) MarkStoresTainted() error {
	return nil
}

func (s *Storage) MarkStoresHealthy() error {
	return nil
}

func (s *Storage) AreStoresCorrupted() (bool, error) {
	return false, nil
}

func (s *Storage) AreStoresTainted() (bool, error) {
	return false, nil
}

func (s *Storage) CheckCorrectStoresVersion() (bool, error) {

	for _, h := range s.healthTrackers {
		correct, err := h.CheckCorrectStoreVersion()
		if err != nil {
			return false, err
		}
		if !correct {
			return false, nil
		}
	}

	return true, nil
}

// UpdateStoresVersion tries to migrate the existing data to the new store version.
func (s *Storage) UpdateStoresVersion() (bool, error) {

	allCorrect := true
	for _, h := range s.healthTrackers {
		_, err := h.UpdateStoreVersion()
		if err != nil {
			return false, err
		}

		correct, err := h.CheckCorrectStoreVersion()
		if err != nil {
			return false, err
		}
		if !correct {
			allCorrect = false
		}
	}

	return allCorrect, nil
}
