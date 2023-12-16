// Code generated by goctl. DO NOT EDIT.
package types

type UploadFileReq struct {
	FileHash  string `json:"file_hash"`
	ChunkHash string `json:"chunk_hash"`
	ChunkNum  int    `json:"chunk_num"`
}

type UploadFileResp struct {
	ChunkSuccess bool `json:"chunk_success"`
}

type MergeChunksReq struct {
	FileHash string `json:"file_hash"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

type MergeChunksResp struct {
	FileID     int64  `json:"file_id"`
	FileURL    string `json:"file_url"`
	UploadTime int64  `json:"upload_time"` // seconds
}

type CheckFileExistReq struct {
	Hash string `path:"hash"`
}

type CheckFileExistResp struct {
	Exist      bool     `json:"exist"`
	FileID     int64    `json:"file_id"`
	FileURL    string   `json:"file_url"`
	ChunksHash []string `json:"chunks_hash"`
}

type GetFileListResp struct {
	FileID     int64  `json:"file_id"`
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
	FileURL    string `json:"file_url"`
	UploadTime int64  `json:"upload_time"` // seconds
}
