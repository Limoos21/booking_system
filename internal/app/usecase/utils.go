package usecase

import (
	"booking_system/internal/domain"
	"errors"
)

func (u UserService) checkGuestCapacity(tables []*domain.Table, reservationCapacity int) (int, []*domain.Table, error) {
	var maxCapacity int
	u.logger.Debug("Checking guest capacity %v", reservationCapacity)
	u.logger.Debug("len table %v", len(tables))
	tabelesDomain := make([]*domain.Table, 0, len(tables))
	for _, t := range tables {
		table, err := u.storage.GetTable(t.ID)
		if err != nil {
			u.logger.Error("Get table error %v", err)
			return 0, nil, err
		}
		u.logger.Debug("table %v", table)
		u.logger.Debug("Checking table capacity %v", table.Capacity)
		tabelesDomain = append(tabelesDomain, table)
		maxCapacity += table.Capacity
	}
	u.logger.Debug("Checking guest capacity max %v", maxCapacity)

	if reservationCapacity > maxCapacity {
		return 0, nil, errors.New("capacity exceeded")
	}

	return maxCapacity, tabelesDomain, nil
}

func (u UserService) checkGuestCapacityMax(tables []domain.Table, reservationCapacity int) (int, []*domain.Table, error) {
	var maxCapacity int
	u.logger.Debug("Checking guest capacity %v", reservationCapacity)
	u.logger.Debug("len table %v", len(tables))
	tabelesDomain := make([]*domain.Table, 0, len(tables))
	for _, t := range tables {
		maxCapacity += t.Capacity
		tabelesDomain = append(tabelesDomain, &t)
	}
	u.logger.Debug("Checking guest capacity max %v", maxCapacity)

	if reservationCapacity > maxCapacity {
		return 0, nil, errors.New("capacity exceeded")
	}

	return maxCapacity, tabelesDomain, nil
}
