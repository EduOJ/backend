package request

// QueryParam: poll
// 1 for poll
// 0 for immediate response
type GetTaskRequest struct {
}

type UpdateRunRequest struct {
	/*
		PENDING / JUDGING / JUDGEMENT_FAILED / NO_COMMENT
		ACCEPTED / WRONG_ANSWER / RUNTIME_ERROR / TIME_LIMIT_EXCEEDED / MEMORY_LIMIT_EXCEEDED / DANGEROUS_SYSCALLS
	*/
	Status     string `json:"status"`
	MemoryUsed uint   `json:"memory_used"` // Byte
	TimeUsed   uint   `json:"time_used"`   // ms
	// 去掉空格回车tab后的sha256
	OutputStrippedHash string `json:"output_stripped_hash"`
	// OutputFile multipart:file
	// CompilerFile multipart:file
	// ComparerFile multipart:file
	Message string `json:"message"`
}
