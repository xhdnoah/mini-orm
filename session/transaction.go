package session

import "mini-orm/log"

// 封装以统一打印日志
func (s *Session) Begin() (err error) {
	log.Info("transaction begins")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Info("transaction commits")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	log.Info("transaction rolls back")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}
