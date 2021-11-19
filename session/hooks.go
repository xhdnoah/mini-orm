package session

import (
	"mini-orm/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// calls the registered hooks
func (s *Session) CallMethod(method string, value interface{}) {
	// 使用反射获取当前会话正在操作对象 value/Model 的方法
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		fm = reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}
