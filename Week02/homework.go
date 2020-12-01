/*
作业：
我们在数据库操作的时候，比如dao层中当遇到一个sql.ErrorNoRows的时候，是否应该Wrap这个error,
抛给上层。为什么，应该怎么做请写出代码？
*/
package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

func main()  {
	user, err := biz()
	if err!=nil {
		log.Printf("biz: %+v\n", err)
		return
	}
	//查询成功，返回数据给用户，这里简单的打印
	fmt.Println("main user", user)
}

type user struct {
	ID uint64
	name string
}

func dao(id uint64) (*user ,error) {
	//get data from db.
	if id <= 100 {
		//查询数据库，这里假设查到数据
		user := &user{ID: id, name: "user_name"}
		return user, nil
	} else  {
		//ErrNoRows
		err := sql.ErrNoRows
		return nil, errors.Wrap(err,"dao error")
	}
}

func biz() (*user ,error) {
	id := uint64(123456789)
	user, err := dao(id)
	//other biz logic
	return user,errors.WithMessage(err, "biz error")
}
