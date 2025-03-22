package main

import (
	"encoding/json"
	"net/http"
)

type metadataResp struct {
	Entry fileMeta
}

func apiMetadata(path string, w http.ResponseWriter, r *http.Request) {
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
