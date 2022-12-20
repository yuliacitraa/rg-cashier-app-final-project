package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type SessionsRepository struct {
	db db.DB
}

func NewSessionsRepository(db db.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) ReadSessions() ([]model.Session, error) {
	records, err := u.db.Load("sessions")
	if err != nil {
		return nil, err
	}

	var listSessions []model.Session
	err = json.Unmarshal([]byte(records), &listSessions)
	if err != nil {
		return nil, err
	}

	return listSessions, nil
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return err
	}

	listSessions = []model.Session{ 
		{
			Token:    "",
			Username: "",
			Expiry:   time.Time{},
		},
    }

	jsonData, err := json.Marshal(listSessions)
	if err != nil {
		return err
	}

	err = u.db.Save("sessions", jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
    add := []model.Session{ 
		{
			Token: session.Token,
			Username: session.Username,
			Expiry:   session.Expiry,
		},
    }

	jsonData, err := json.Marshal(add)
    if err != nil {
        panic(err)
    }

	err = u.db.Save("sessions", jsonData)
	if err != nil {
		return err
	}
	
	return nil
}

func (u *SessionsRepository) CheckExpireToken(token string) (model.Session, error) {
	listOfToken, err := u.ReadSessions()
	if err != nil{
		return model.Session{}, err
	}

	session := model.Session{}
	for _, val := range listOfToken {
		if val.Token == token {
			session = val
		}
	}

	expired := session.Expiry
	now := time.Now()
	if now.After(expired){
		return model.Session{}, errors.New("Token is Expired!")
	}

	return session, nil
}

func (u *SessionsRepository) ResetSessions() error {
	err := u.db.Reset("sessions", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) TokenExist(req string) (model.Session, error) {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return model.Session{}, err
	}
	for _, element := range listSessions {
		if element.Token == req {
			return element, nil
		}
	}
	return model.Session{}, fmt.Errorf("Token Not Found!")
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
