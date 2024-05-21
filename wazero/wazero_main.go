package main

// #include <stdlib.h>
import "C"

import (
	"encoding/json"
	wescale_wasm_plugin_template "wescale-wasm-plugin-template"
	"wescale-wasm-plugin-template/tools"
	hostfunction "wescale-wasm-plugin-template/tools/host_functions"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

//export WazeroGuestFuncBeforeExecution
func WazeroGuestFuncBeforeExecution(hostInstancePtr, hostModulePtr uint64) {
	hostfunction.HostInstancePtr = hostInstancePtr
	hostfunction.HostModulePtr = hostModulePtr

	wescale_wasm_plugin_template.RunBeforeExecution()
	
}

//export wazeroGuestFuncAfterExecution
func wazeroGuestFuncAfterExecution(ptr, size uint32) (ptrSize uint64) {
	dataFromHost := tools.PtrToString(ptr, size)

	w := tools.WasmPluginRunAfterExecutionExchange{}
	// todo, how to handle errors?
	json.Unmarshal([]byte(dataFromHost), &w)

	wescale_wasm_plugin_template.RunAfterExecution(&w)

	// todo, how to handle errors?
	dataToHost, _ := json.Marshal(&w)
	dataToHostString := string(dataToHost)

	ptr, size = tools.StringToPtr(dataToHostString)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

//export proxy_on_memory_allocate
func proxyOnMemoryAllocate(size uint) *byte {
	buf := make([]byte, size)
	return &buf[0]
}
