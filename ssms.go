//go:build windows

package main

import (
    "bufio"
    "fmt"
    "syscall"
    "unsafe"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

func main() {

    // Try current working directory first
    cwd, err := os.Getwd()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
        os.Exit(1)
    }
    confPath := filepath.Join(cwd, "ssms.conf")
    ssmsPath, err := readConfig(confPath)
    if err != nil {
        // Fallback to executable directory
        exePath, err2 := os.Executable()
        if err2 != nil {
            fmt.Fprintf(os.Stderr, "Error getting executable path: %v\n", err2)
            os.Exit(1)
        }
        exeDir := filepath.Dir(exePath)
        confPath = filepath.Join(exeDir, "ssms.conf")
        ssmsPath, err = readConfig(confPath)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
            os.Exit(1)
        }
    }

    args, password := filterArgs(os.Args[1:])

    if password != "" {
        err := setClipboard(password)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to set clipboard: %v\n", err)
            os.Exit(1)
        }
//?        fmt.Println("Password copied to clipboard.")
    }

    cmd := exec.Command(ssmsPath, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    err = cmd.Start()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to start SSMS: %v\n", err)
        os.Exit(1)
    }
    // Do not wait for SSMS to exit; just terminate this shim.
}

// Reads the config file and returns the SSMS path.
func readConfig(confPath string) (string, error) {
    f, err := os.Open(confPath)
    if err != nil {
        return "", err
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" && !strings.HasPrefix(line, "#") {
            return line, nil
        }
    }
    return "", fmt.Errorf("no valid SSMS path found in config")
}

// Filters out -P and its value, returns remaining args and the password (if any).
func filterArgs(args []string) ([]string, string) {
    var out []string
    var password string
    skip := false
    for i := 0; i < len(args); i++ {
        if skip {
            skip = false
            continue
        }
        if args[i] == "-P" && i+1 < len(args) {
            password = args[i+1]
            skip = true
            continue
        }
        out = append(out, args[i])
    }
    return out, password
}

// Sets the clipboard text (Windows only).
func setClipboard(text string) error {
    user32 := syscall.NewLazyDLL("user32.dll")
    kernel32 := syscall.NewLazyDLL("kernel32.dll")

    openClipboard := user32.NewProc("OpenClipboard")
    closeClipboard := user32.NewProc("CloseClipboard")
    emptyClipboard := user32.NewProc("EmptyClipboard")
    setClipboardData := user32.NewProc("SetClipboardData")
    globalAlloc := kernel32.NewProc("GlobalAlloc")
    globalLock := kernel32.NewProc("GlobalLock")
    globalUnlock := kernel32.NewProc("GlobalUnlock")
    globalFree := kernel32.NewProc("GlobalFree")

    const GMEM_MOVEABLE = 0x0002
    const CF_TEXT = 1

    // Open clipboard
    r, _, err := openClipboard.Call(0)
    if r == 0 {
        return fmt.Errorf("OpenClipboard failed: %v", err)
    }
    defer closeClipboard.Call()

    // Empty clipboard
    r, _, err = emptyClipboard.Call()
    if r == 0 {
        return fmt.Errorf("EmptyClipboard failed: %v", err)
    }

    // Allocate global memory for the text (+1 for null terminator)
    hMem, _, err := globalAlloc.Call(GMEM_MOVEABLE, uintptr(len(text)+1))
    if hMem == 0 {
        return fmt.Errorf("GlobalAlloc failed: %v", err)
    }
    defer globalFree.Call(hMem)

    // Lock the memory and copy the text
    lpMem, _, err := globalLock.Call(hMem)
    if lpMem == 0 {
        return fmt.Errorf("GlobalLock failed: %v", err)
    }
    // Copy text to memory
    src := append([]byte(text), 0) // null-terminated
    mem := (*[1 << 20]byte)(unsafe.Pointer(lpMem))
    copy(mem[:], src)
    globalUnlock.Call(hMem)

    // Set clipboard data
    r, _, err = setClipboardData.Call(CF_TEXT, hMem)
    if r == 0 {
        return fmt.Errorf("SetClipboardData failed: %v", err)
    }
    // Do not free hMem after SetClipboardData succeeds

    return nil
}

