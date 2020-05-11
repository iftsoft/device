package validator

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
)


type ValidatorProtocol struct {
	config  *config.LinkerConfig
	log     *core.LogAgent
	timeout uint16
}

func GetValidatorProtocol(cfg *config.LinkerConfig) *ValidatorProtocol {
	lbp := &ValidatorProtocol{
		config: cfg,
		log: core.GetLogAgent(core.LogLevelDump, "Linker"),
		timeout:    cfg.Timeout,
	}
	if lbp.timeout == 0 {
		lbp.timeout = 250
	}
	return lbp
}


////////////////////////////////////////////////////////////////

func (lbp *ValidatorProtocol) CheckLink() error {
	data := []byte{0xAA, 0x55, 0x00, 0xFF}
	back, err := lbp.exchange(data)
	if err == nil {
		err = lbp.checkReply(data, back)
	}
	lbp.logError("CheckLink", err)
	return err
}

////////////////////////////////////////////////////////////////

func (lbp *ValidatorProtocol) logError(cmd string, err error) {
	code, text := common.CheckError(err)
	lbp.log.Trace("ValidatorProtocol.%s return: %d - %s", cmd, code, text)
}

func (lbp *ValidatorProtocol) checkReply(data, back []byte) error {
	if len(data) != len(back) {
		return common.NewError(common.DevErrorLinkerFault, "length mismatch")
	}
	for i:=0; i<len(back); i++ {
		if data[i] != back[i] {
			return common.NewError(common.DevErrorLinkerFault, "data mismatch")
		}
	}
	return nil
}

func (lbp *ValidatorProtocol) exchange(data []byte) ([]byte, error) {
	return data, nil
}

