// Code generated by goctl. DO NOT EDIT.
package types

type UploadFileResp struct {
	FileID  int64  `json:"file_id"`
	FileURL string `json:"file_url"`
}

type CheckFileExistReq struct {
	Hash string `json:"hash"`
}

type CheckFileExistResp struct {
	Exist   bool   `json:"exist"`
	FileID  int64  `json:"file_id"`
	FileURL string `json:"file_url"`
}