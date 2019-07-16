// handlers_test.go
package handlers

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"bytes"
	"reflect"
	"testing"
)

func TestHandleConfig(t *testing.T) {
	SetupConfig()

	req, err := http.NewRequest("GET", "/config.json", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleConfig)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	responseConfig := Config{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseConfig)
	if err != nil {
		t.Fatalf("error unmarshalling config response: %v\n", err)
	}

	// Check the response body is what we expect.
	expectedConfig := Config{
		Local:    true,
		Protocol: "ws://",
	}
	if eq := reflect.DeepEqual(expectedConfig, responseConfig); !eq {
		t.Fatalf("expectedConfig not equal to responseConfig\nexpected:\n%v\ngot:\n%v\n",
			expectedConfig, responseConfig)
	}
}

func TestHandleCharacterCreate(t *testing.T) {
	SetupConfig()
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// set database to be our mock db
	Database = db
	defer func() {
		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}()

	testBadBody := func() {
		req, err := http.NewRequest("POST", "/characters/create", bytes.NewReader([]byte("{zz}")))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleCharacterCreate)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)
		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("wrong status code: got %v want %v", status, http.StatusOK)
		}
		expectedBody := "Could not process JSON body!"
		if rr.Body.String() != expectedBody {
			t.Fatalf("body not equal to expected body\nexpected:\n%v\ngot:\n%v\n",
				expectedBody, rr.Body.String())
		}

	}

	testSuccessfulCreation := func() {
		character := Character{
			CharacterId:   1,
			CharacterName: "elon",
			Attack:        420,
			Defense:       100,
			Health:        100,
			UserId:        420,
		}

		mock.ExpectExec("INSERT INTO Characters").WithArgs(character.CharacterName, character.Attack,
			character.Defense, character.Health, character.UserId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		marshalledCharacter, err := json.Marshal(character)
		if err != nil {
			t.Fatalf("error marshalling character: %v", err)
		}

		req, err := http.NewRequest("POST", "/characters/create", bytes.NewReader(marshalledCharacter))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleCharacterCreate)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)
		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
		}

		expectedBody := "Successfully created character"
		if rr.Body.String() != expectedBody {
			t.Fatalf("body not equal to expected body\nexpected:\n%v\ngot:\n%v\n",
				expectedBody, rr.Body.String())
		}

	}

	testBadBody()
	testSuccessfulCreation()

}

func TestHandleUserCreate(t *testing.T) {
	SetupConfig()
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// set database to be our mock db
	Database = db
	defer func() {
		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}()

	testBadBody := func() {
		req, err := http.NewRequest("POST", "/users/create", bytes.NewReader([]byte("{zz}")))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleUserCreate)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)
		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
		expectedBody := "Could not process JSON body!"
		if rr.Body.String() != expectedBody {
			t.Fatalf("body not equal to expected body\nexpected:\n%v\ngot:\n%v\n",
				expectedBody, rr.Body.String())
		}
	}

	testUserExists := func() {
		user := User{
			Username: "ilon",
			Password: "mask",
			FullName: "ilonmask",
		}
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		marshalledUser, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("error marshalling character: %v", err)
		}

		req, err := http.NewRequest("POST", "/users/create", bytes.NewReader(marshalledUser))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleUserCreate)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)
		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusConflict {
			t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
		}

		expectedBody := "user already exists. error: <nil>"
		if rr.Body.String() != expectedBody {
			t.Fatalf("body not equal to expected body\nexpected:\n%v\ngot:\n%v\n",
				expectedBody, rr.Body.String())
		}

	}

	testSuccessfulCreation := func() {
		user := User{
			Username: "ilon",
			Password: "mask",
			FullName: "ilonmask",
		}
		rows := sqlmock.NewRows([]string{"id"})

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		mock.ExpectExec("INSERT INTO Users").
			WillReturnResult(sqlmock.NewResult(1, 1))

		marshalledUser, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("error marshalling character: %v", err)
		}

		req, err := http.NewRequest("POST", "/users/create", bytes.NewReader(marshalledUser))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleUserCreate)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)
		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
		}

		expectedBody := "Successfully created user"
		if rr.Body.String() != expectedBody {
			t.Fatalf("body not equal to expected body\nexpected:\n%v\ngot:\n%v\n",
				expectedBody, rr.Body.String())
		}

	}

	testBadBody()
	testUserExists()
	testSuccessfulCreation()

}

func TestHandleUserExists(t *testing.T) {
	SetupConfig()
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// set database to be our mock db
	Database = db
	defer func() {
		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}()

	testUserExists := func() {
		rows := sqlmock.NewRows([]string{"id", "username", "fullname"}).AddRow(420, "ilon", "ilonmask")

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		req, err := http.NewRequest("GET", "/users/ilon", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleUserExists)
		vars := map[string]string{
			"username": "ilon",
		}

		// Hack to try to fake gorilla/mux vars
		req = mux.SetURLVars(req, vars)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
		}

		responseUser := User{}
		err = json.Unmarshal(rr.Body.Bytes(), &responseUser)
		if err != nil {
			t.Fatalf("error unmarshalling user response: %v\n", err)
		}

		expectedUser := User{
			Id:       420,
			Username: "ilon",
			FullName: "ilonmask",
		}
		if eq := reflect.DeepEqual(expectedUser, responseUser); !eq {
			t.Fatalf("expectedUser not equal to responseUser\nexpected:\n%v\ngot:\n%v\n",
				expectedUser, responseUser)
		}

	}

	testUserDoesNotExist := func() {
		rows := sqlmock.NewRows([]string{"id", "username", "fullname"})

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		req, err := http.NewRequest("GET", "/users/ilon", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleUserExists)
		vars := map[string]string{
			"username": "ilon",
		}

		//Hack to try to fake gorilla/mux vars
		req = mux.SetURLVars(req, vars)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	}

	testUserExists()
	testUserDoesNotExist()

}

func TestHandleUserLogin(t *testing.T) {
	// Todo: how to handle bcrypt?
	//bcrypt.CompareHashAndPassword = func(a, b []byte) error {
	//	return nil
	//}
}