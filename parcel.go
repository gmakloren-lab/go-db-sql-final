package main

import (
	"database/sql"
	"errors"
)

type ParcelStore struct {
	db *sql.DB
}

// NewParcelStore создаёт хранилище посылок с подключением к БД
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Add добавляет новую посылку в БД и возвращает её идентификатор
func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(
		`INSERT INTO parcel (client, status, address, created_at)
		 VALUES (?, ?, ?, ?)`,
		p.Client,
		p.Status,
		p.Address,
		p.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get возвращает посылку по её идентификатору
func (s ParcelStore) Get(number int) (Parcel, error) {
	var p Parcel

	err := s.db.QueryRow(
		`SELECT number, client, status, address, created_at
		 FROM parcel
		 WHERE number = ?`,
		number,
	).Scan(
		&p.Number,
		&p.Client,
		&p.Status,
		&p.Address,
		&p.CreatedAt,
	)

	if err != nil {
		return p, err
	}

	return p, nil
}

// GetByClient возвращает список всех посылок указанного клиента
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query(
		`SELECT number, client, status, address, created_at
		 FROM parcel
		 WHERE client = ?
		 ORDER BY number`,
		client,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(
			&p.Number,
			&p.Client,
			&p.Status,
			&p.Address,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}

// SetStatus обновляет статус посылки
func (s ParcelStore) SetStatus(number int, status string) error {
	res, err := s.db.Exec(
		`UPDATE parcel
		 SET status = ?
		 WHERE number = ?`,
		status,
		number,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("parcel not found")
	}

	return nil
}

// SetAddress изменяет адрес доставки,
// разрешено только для посылок со статусом "registered"
func (s ParcelStore) SetAddress(number int, address string) error {
	res, err := s.db.Exec(
		`UPDATE parcel
		 SET address = ?
		 WHERE number = ?
		   AND status = ?`,
		address,
		number,
		ParcelStatusRegistered,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("address change not allowed")
	}

	return nil
}

// Delete удаляет посылку,
// разрешено только для посылок со статусом "registered"
func (s ParcelStore) Delete(number int) error {
	res, err := s.db.Exec(
		`DELETE FROM parcel
		 WHERE number = ?
		   AND status = ?`,
		number,
		ParcelStatusRegistered,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("delete not allowed")
	}

	return nil
}
