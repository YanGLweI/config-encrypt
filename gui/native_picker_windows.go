//go:build windows

package gui

/*
#cgo LDFLAGS: -lole32 -lcomdlg32 -lshell32 -luuid
#include <stdlib.h>
#include <stdio.h>
#include <windows.h>
#include <shobjidl.h>
#include <objbase.h>
#include <shlwapi.h>

const char* openFilePanel(const char* title, const char* defaultDir, const char** extensions, int extCount) {
    HRESULT hr = CoInitializeEx(NULL, COINIT_APARTMENTTHREADED | COINIT_DISABLE_OLE1DDE);
    if (FAILED(hr)) return NULL;

    IFileOpenDialog* pfd = NULL;
    hr = CoCreateInstance(&CLSID_FileOpenDialog, NULL, CLSCTX_INPROC_SERVER, &IID_IFileOpenDialog, (void**)&pfd);
    if (FAILED(hr)) { CoUninitialize(); return NULL; }

    // 设置标题
    if (title && strlen(title) > 0) {
        wchar_t wTitle[256];
        MultiByteToWideChar(CP_UTF8, 0, title, -1, wTitle, 256);
        pfd->lpVtbl->SetTitle(pfd, wTitle);
    }

    // 设置初始目录
    if (defaultDir && strlen(defaultDir) > 0) {
        wchar_t wDir[MAX_PATH];
        MultiByteToWideChar(CP_UTF8, 0, defaultDir, -1, wDir, MAX_PATH);
        IShellItem* psi = NULL;
        hr = SHCreateItemFromParsingName(wDir, NULL, &IID_IShellItem, (void**)&psi);
        if (SUCCEEDED(hr)) {
            pfd->lpVtbl->SetFolder(pfd, psi);
            psi->lpVtbl->Release(psi);
        }
    }

    // 设置文件类型过滤
    if (extensions && extCount > 0) {
        COMDLG_FILTERSPEC* filters = (COMDLG_FILTERSPEC*)CoTaskMemAlloc(sizeof(COMDLG_FILTERSPEC) * extCount);
        for (int i = 0; i < extCount; i++) {
            char spec[64];
            snprintf(spec, sizeof(spec), "*.%s", extensions[i]);
            wchar_t wSpec[64], wName[64];
            MultiByteToWideChar(CP_UTF8, 0, spec, -1, wSpec, 64);
            MultiByteToWideChar(CP_UTF8, 0, extensions[i], -1, wName, 64);
            filters[i].pszName = (LPCWSTR)CoTaskMemAlloc((wcslen(wName)+1)*sizeof(wchar_t));
            wcscpy((wchar_t*)filters[i].pszName, wName);
            filters[i].pszSpec = (LPCWSTR)CoTaskMemAlloc((wcslen(wSpec)+1)*sizeof(wchar_t));
            wcscpy((wchar_t*)filters[i].pszSpec, wSpec);
        }
        pfd->lpVtbl->SetFileTypes(pfd, extCount, filters);
        // 注意: filters 内存由 COM 管理，不需要手动释放
    }

    hr = pfd->lpVtbl->Show(pfd, NULL);
    if (SUCCEEDED(hr)) {
        IShellItem* psiResult = NULL;
        hr = pfd->lpVtbl->GetResult(pfd, &psiResult);
        if (SUCCEEDED(hr)) {
            LPWSTR pszPath = NULL;
            hr = psiResult->lpVtbl->GetDisplayName(psiResult, SIGDN_FILESYSPATH, &pszPath);
            if (SUCCEEDED(hr) && pszPath) {
                int len = WideCharToMultiByte(CP_UTF8, 0, pszPath, -1, NULL, 0, NULL, NULL);
                char* path = (char*)malloc(len);
                WideCharToMultiByte(CP_UTF8, 0, pszPath, -1, path, len, NULL, NULL);
                CoTaskMemFree(pszPath);
                psiResult->lpVtbl->Release(psiResult);
                pfd->lpVtbl->Release(pfd);
                CoUninitialize();
                return path;
            }
            if (psiResult) psiResult->lpVtbl->Release(psiResult);
        }
    }

    pfd->lpVtbl->Release(pfd);
    CoUninitialize();
    return NULL;
}

const char* openFolderPanel(const char* title, const char* defaultDir) {
    HRESULT hr = CoInitializeEx(NULL, COINIT_APARTMENTTHREADED | COINIT_DISABLE_OLE1DDE);
    if (FAILED(hr)) return NULL;

    IFileOpenDialog* pfd = NULL;
    hr = CoCreateInstance(&CLSID_FileOpenDialog, NULL, CLSCTX_INPROC_SERVER, &IID_IFileOpenDialog, (void**)&pfd);
    if (FAILED(hr)) { CoUninitialize(); return NULL; }

    if (title && strlen(title) > 0) {
        wchar_t wTitle[256];
        MultiByteToWideChar(CP_UTF8, 0, title, -1, wTitle, 256);
        pfd->lpVtbl->SetTitle(pfd, wTitle);
    }

    if (defaultDir && strlen(defaultDir) > 0) {
        wchar_t wDir[MAX_PATH];
        MultiByteToWideChar(CP_UTF8, 0, defaultDir, -1, wDir, MAX_PATH);
        IShellItem* psi = NULL;
        hr = SHCreateItemFromParsingName(wDir, NULL, &IID_IShellItem, (void**)&psi);
        if (SUCCEEDED(hr)) {
            pfd->lpVtbl->SetFolder(pfd, psi);
            psi->lpVtbl->Release(psi);
        }
    }

    DWORD dwFlags;
    pfd->lpVtbl->GetOptions(pfd, &dwFlags);
    pfd->lpVtbl->SetOptions(pfd, dwFlags | FOS_PICKFOLDERS);

    hr = pfd->lpVtbl->Show(pfd, NULL);
    if (SUCCEEDED(hr)) {
        IShellItem* psiResult = NULL;
        hr = pfd->lpVtbl->GetResult(pfd, &psiResult);
        if (SUCCEEDED(hr)) {
            LPWSTR pszPath = NULL;
            hr = psiResult->lpVtbl->GetDisplayName(psiResult, SIGDN_FILESYSPATH, &pszPath);
            if (SUCCEEDED(hr) && pszPath) {
                int len = WideCharToMultiByte(CP_UTF8, 0, pszPath, -1, NULL, 0, NULL, NULL);
                char* path = (char*)malloc(len);
                WideCharToMultiByte(CP_UTF8, 0, pszPath, -1, path, len, NULL, NULL);
                CoTaskMemFree(pszPath);
                psiResult->lpVtbl->Release(psiResult);
                pfd->lpVtbl->Release(pfd);
                CoUninitialize();
                return path;
            }
            if (psiResult) psiResult->lpVtbl->Release(psiResult);
        }
    }

    pfd->lpVtbl->Release(pfd);
    CoUninitialize();
    return NULL;
}
*/
import "C"

import (
	"unsafe"
)

// nativeOpenFile 调用 Windows 原生 IFileOpenDialog 选择文件
// 返回 (路径, 是否使用了原生选择器)。用户取消时返回 ("", true)
func nativeOpenFile(title string, defaultDir string, extensions []string) (string, bool) {
	cTitle := C.CString(title)
	cDir := C.CString(defaultDir)
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cDir))

	var cExts []*C.char
	for _, ext := range extensions {
		cExts = append(cExts, C.CString(ext))
	}
	var cExtPtr **C.char
	if len(cExts) > 0 {
		cExtPtr = (**C.char)(unsafe.Pointer(&cExts[0]))
	}

	result := C.openFilePanel(cTitle, cDir, cExtPtr, C.int(len(extensions)))
	if result == nil {
		return "", true // 用户取消
	}
	defer C.free(unsafe.Pointer(result))
	return C.GoString(result), true
}

// nativeOpenFolder 调用 Windows 原生 IFileOpenDialog 选择文件夹
// 返回 (路径, 是否使用了原生选择器)。用户取消时返回 ("", true)
func nativeOpenFolder(title string, defaultDir string) (string, bool) {
	cTitle := C.CString(title)
	cDir := C.CString(defaultDir)
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cDir))

	result := C.openFolderPanel(cTitle, cDir)
	if result == nil {
		return "", true // 用户取消
	}
	defer C.free(unsafe.Pointer(result))
	return C.GoString(result), true
}
