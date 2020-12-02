/**
 * @Author chenjun
 * @Date 2020/12/2
 **/
package main

import (
	"database/sql"
	 "github.com/pkg/errors"
)
//我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码

func main() {
	service()
}
func dao() error  {
	return errors.Wrap(sql.ErrNoRows,"dao error")
}
func service() error  {
return dao()
}