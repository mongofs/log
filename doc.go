// log 文件作为记录错误的
package log


type Logger interface {
	Info()
	Error()
}