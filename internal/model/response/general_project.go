package response

import (
	"TTPanel/internal/model"
)

type GeneralProject struct {
	*model.Project
	MainProcess    *ProcessInfoP               `json:"main_process_list"`
	ChildProcess   []*ProcessInfoP             `json:"child_process_list"`
	RunLogFilePath string                      `json:"run_log_file_path"`
	ProjectConfig  *model.GeneralProjectConfig `json:"project_config"`
}
