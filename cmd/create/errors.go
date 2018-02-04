package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stellar/go/clients/horizon"
)

func GetTxErrorResultCodes(err error, logger log.Logger) *horizon.TransactionResultCodes {
	level.Error(logger).Log("msg", err)
	switch e := err.(type) {
	case *horizon.Error:
		code, err := e.ResultCodes()
		if err != nil {
			level.Error(logger).Log("msg", "failed to extract result codes from horizon response")
			return nil
		}
		level.Error(logger).Log("code", code.TransactionCode)
		for i, opCode := range code.OperationCodes {
			level.Error(logger).Log("opcode_index", i, "opcode", opCode)
		}

		return code
	}
	return nil
}
