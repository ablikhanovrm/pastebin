package paste

import "github.com/google/uuid"

func parsePasteUuids(ids []string) ([]uuid.UUID, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	parsePasteUUIDs := make([]uuid.UUID, len(ids))

	for i, id := range ids {
		u, err := uuid.Parse(id)

		if err != nil {
			return nil, err
		}

		parsePasteUUIDs[i] = u
	}

	return parsePasteUUIDs, nil
}
