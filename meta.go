package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
)

type metadataReq struct {
	Path string
}

type metadataResp struct {
	Entry fileMeta
}

func apiMetadata(w http.ResponseWriter, r *http.Request) {
	var reqModel metadataReq

	d := json.NewDecoder(r.Body)
	err := d.Decode(&reqModel)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

    path := filepath.Join(dir, reqModel.Path)

	var respModel metadataResp
    entry, err := fileMetaOf(path)
	if err != nil {
		handleNotFoundOrInternalErr(w, err)
		return
	}
	respModel.Entry = entry
	e := json.NewEncoder(w)
	e.Encode(respModel)
}
