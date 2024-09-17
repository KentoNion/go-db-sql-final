package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) { //Функция добавления новой записи
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:Client, :Status, :Address, :CreatedAt)",
		sql.Named("Client", p.Client),
		sql.Named("Status", p.Status),
		sql.Named("Address", p.Address),
		sql.Named("CreatedAt", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	// Считываем идефикатор последней изменённой строки
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) { //функция получения информации о послыке по number
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))
	// запонляем срез p данными
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) { //функция получения информации о всех посылках client
	var parcels []Parcel
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//заполняем срез p данными о посылках пока есть Next
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}
	err = rows.Err() //Проверяем курсор на наличие ошибки
	if err != nil {
		return nil, err
	}
	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error { //функция обновления статуса в посылке

	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("number", number),
		sql.Named("status", status))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error { //Функция обновления адресса в посылке, при условии status = registered
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("address", address),
		sql.Named("status", "registered"))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error { //Удаление строки из parcel если status = registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", "registered"))
	if err != nil {
		return err
	}

	return nil
}
