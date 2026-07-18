//go:build darwin

package gui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit
#import <AppKit/AppKit.h>
#include <stdlib.h>

const char* openFilePanel(const char* title, const char* defaultDir, const char** extensions, int extCount) {
    @autoreleasepool {
        NSOpenPanel* panel = [NSOpenPanel openPanel];
        [panel setCanChooseFiles:YES];
        [panel setCanChooseDirectories:NO];
        [panel setAllowsMultipleSelection:NO];
        [panel setTitle:[NSString stringWithUTF8String:title]];

        if (defaultDir && strlen(defaultDir) > 0) {
            [panel setDirectoryURL:[NSURL fileURLWithPath:[NSString stringWithUTF8String:defaultDir]]];
        }

        if (extensions && extCount > 0) {
            NSMutableArray* exts = [NSMutableArray arrayWithCapacity:extCount];
            for (int i = 0; i < extCount; i++) {
                [exts addObject:[NSString stringWithUTF8String:extensions[i]]];
            }
            [panel setAllowedFileTypes:exts];
        }

        if ([panel runModal] == NSModalResponseOK) {
            NSURL* url = [[panel URLs] objectAtIndex:0];
            const char* path = [[url path] UTF8String];
            return strdup(path);
        }
        return NULL;
    }
}

const char* openFolderPanel(const char* title, const char* defaultDir) {
    @autoreleasepool {
        NSOpenPanel* panel = [NSOpenPanel openPanel];
        [panel setCanChooseFiles:NO];
        [panel setCanChooseDirectories:YES];
        [panel setAllowsMultipleSelection:NO];
        [panel setTitle:[NSString stringWithUTF8String:title]];

        if (defaultDir && strlen(defaultDir) > 0) {
            [panel setDirectoryURL:[NSURL fileURLWithPath:[NSString stringWithUTF8String:defaultDir]]];
        }

        if ([panel runModal] == NSModalResponseOK) {
            NSURL* url = [[panel URLs] objectAtIndex:0];
            const char* path = [[url path] UTF8String];
            return strdup(path);
        }
        return NULL;
    }
}
*/
import "C"

import (
	"unsafe"
)

// nativeOpenFile 调用 macOS 原生 NSOpenPanel 选择文件
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

// nativeOpenFolder 调用 macOS 原生 NSOpenPanel 选择文件夹
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
