package api

import "browser/handle"

func ListenBlock(){
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	fabsdk.ListenBlock()
}
