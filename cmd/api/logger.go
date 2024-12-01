package main

type logger interface {
	WriteInfo(msg string, fields map[string]interface{})
	WriteError(msg string, err error, fields map[string]interface{})
	WriteFatal(msg string, err error, fields map[string]interface{})
}
