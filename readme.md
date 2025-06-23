# SSMS Shim

This project provides a Windows-only shim executable for launching SQL Server Management Studio (SSMS) with special handling for the `-P` (password) argument. When a password is provided, it is copied to the clipboard instead of being passed on the command line, improving security.

---

## Prerequisites

- **Operating System:** Windows 10 or later
- **Go:** Version 1.18 or newer ([Download Go](https://go.dev/dl/))
- **Git:** (optional, for cloning the repository)

---

## Building the Project

### 1. Install Go

If you do not have Go installed, download and install it from [https://go.dev/dl/](https://go.dev/dl/).

After installation, verify Go is available in your terminal:

```sh
go version
```

### 2. Clone or Download the Repository

If you have Git:

```sh
git clone https://your-repo-url/ssms-shim.git
cd ssms-shim
```

Or, download the source code as a ZIP and extract it.

### 3. Build the Executable

Open a **Command Prompt** or **PowerShell** in the project directory and run:

```sh
go build -o ssms.exe ssms.go
```

This will produce `ssms.exe` in the current directory.

---

## Configuration

Create a file named `ssms.conf` in the same directory as the executable or in your current working directory.  
This file should contain a single line: the full path to your `ssms.exe` (SQL Server Management Studio executable).

**Example `ssms.conf`:**
```
C:\Program Files (x86)\Microsoft SQL Server Management Studio 18\Common7\IDE\Ssms.exe
```

---

## Usage

Run the shim just like you would run SSMS, passing any arguments you need.

**Example:**
```sh
ssms.exe -S myserver -U myuser -P mypassword
```

- The `-P` argument (password) will be copied to the clipboard instead of being passed to SSMS.
- All other arguments are forwarded to SSMS.

---

## Notes

- This project only works on Windows.
- The clipboard is set using Windows API calls; no external dependencies are required.
- If `ssms.conf` is not found in the current directory, the shim will look in the executable's directory.

---

## Troubleshooting

- **"Error reading config"**: Make sure `ssms.conf` exists and contains the correct path to `ssms.exe`.
- **"Failed to set clipboard"**: Ensure you have permission to access the clipboard (try running as administrator).

---

## License

MIT License

---
