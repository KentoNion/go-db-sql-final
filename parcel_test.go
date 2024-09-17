package main

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db") //Подключение к БД
	require.NoError(t, err)                     //проверка на ошибку
	defer db.Close()                            //закрываем бд через дефер
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel) //добавляем посылку
	require.NotEqual(t, 0, id)   //Проверяем чтоб id не был равен 0
	require.NoError(t, err)      //проверяем отсутсвие ошибки
	parcel.Number = id           //Так как наша тестовая посылка не имеет параметра number, по умолчанию ей присваивается 0, мы это исправляем вручную присваивая параметру number id

	p, err := store.Get(id)     //получаем только что добавленную посылку
	require.NoError(t, err)     //проверяем отсутсвие ошибки
	require.Equal(t, parcel, p) //проверяем что то что мы добавили = тому что мы добавляли

	err = store.Delete(id)  //удаляем посылку
	require.NoError(t, err) //Проверяем отсутвие ошибки
	_, err = store.Get(id)  //получаем несуществующую посылку
	require.Error(t, err)   //Проверка наличия ошибки
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") //Подключение к БД
	require.NoError(t, err)                     //проверка на ошибку
	defer db.Close()                            //закрываем бд через дефер
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	id, err := store.Add(parcel) //добавляем посылку
	require.NoError(t, err)      //проверяем отсутсвие ошибки
	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)  //меняем адресс
	require.NoError(t, err)                 //проверяем отсутсвие ошибки
	a, err := store.Get(id)                 //Получаем "новую" посылку из дб с новым адрессом
	require.NoError(t, err)                 //Проверяем что нет ошибки
	require.Equal(t, newAddress, a.Address) //Проверяем что адресс изменился
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") //Подключение к БД
	require.NoError(t, err)                     //проверка на ошибку
	defer db.Close()                            //закрываем бд через дефер
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel) //добавляем посылку
	require.NoError(t, err)      //проверяем отсутсвие ошибки

	err = store.SetStatus(id, "kek") //присваиваем статус
	require.NoError(t, err)          //проверяем отсутсвие ошибки

	newParcel, err := store.Get(id)           //Присваиваем переменной значение посылки из дб
	require.Equal(t, "kek", newParcel.Status) //проверяем что новый статус равен тому что мы установили
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") //Подключение к БД
	require.NoError(t, err)                     //проверка на ошибку
	defer db.Close()                            //закрываем бд через дефер
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i]) //Добавляем новые посылки в бд
		require.NoError(t, err)          //Проверяем отсутсвие ошибки
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)    //Получаем посылки по индефикатору плиента
	require.NoError(t, err)                            //проверяем отсутсвие ошибок
	require.Equal(t, len(parcels), len(storedParcels)) //проверяем что кол-во полученных посылок = кол-во доставленных посылок

	// check
	for _, parcel := range storedParcels {
		assert.NotEmpty(t, parcelMap[parcel.Number])
		require.Equal(t, parcel, parcelMap[parcel.Number])
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
	}
}
