package controllers

import "github.com/Oudwins/zog"

func formatErrors(errs zog.ZogIssueList) map[string]string {
	errorMap := make(map[string]string)
	for _, issue := range errs {
		if _, exists := errorMap[issue.Path[0]]; !exists {
			errorMap[issue.Path[0]] = issue.Message
		}
	}
	return errorMap
}
