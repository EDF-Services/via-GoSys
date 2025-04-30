package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/host"
    "github.com/shirou/gopsutil/net"
    "github.com/shirou/gopsutil/process"
    "github.com/spf13/cobra"
)

// Available package repositories
var packageRepos = map[string]string{
    "GoMii":      "https://github.com/example/GoMii.git",
    "via-GoSys":  "https://github.com/example/via-GoSys.git",
    "Goml":       "https://github.com/example/Goml.git",
    "vanilla-go": "https://github.com/example/vanilla-go.git",
}

func installPackage(pkgName string) {
    repoURL, exists := packageRepos[pkgName]
    if !exists {
        fmt.Printf("❌ Error: Package '%s' not found.\n", pkgName)
        return
    }

    fmt.Printf("Installing %s from GitHub...\n", pkgName)
    cloneCmd := exec.Command("git", "clone", repoURL)
    cloneCmd.Stdout = os.Stdout
    cloneCmd.Stderr = os.Stderr
    err := cloneCmd.Run()

    if err != nil {
        fmt.Printf("❌ Error installing %s: %v\n", pkgName, err)
    } else {
        fmt.Println("✅ Installation complete!")
    }
}

func systemInfo() {
    info, _ := host.Info()
    fmt.Printf("OS: %s %s\n", info.Platform, info.PlatformVersion)
    fmt.Printf("Hostname: %s\n", info.Hostname)
}

func diskUsage() {
    usage, _ := disk.Usage("/")
    fmt.Printf("Disk Usage: %v%% used, %v free\n", usage.UsedPercent, usage.Free)
}

func networkInfo() {
    interfaces, _ := net.Interfaces()
    for _, iface := range interfaces {
        fmt.Printf("Interface: %s, MAC: %s\n", iface.Name, iface.HardwareAddr)
    }
}

func uptime() {
    info, _ := host.Uptime()
    fmt.Printf("Uptime: %v seconds\n", info)
}

func listProcesses() {
    procs, _ := process.Processes()
    fmt.Println("Running processes:")
    for _, proc := range procs {
        name, _ := proc.Name()
        fmt.Printf("PID: %d, Name: %s\n", proc.Pid, name)
    }
}

func killProcess(pid int32) {
    proc, err := process.NewProcess(pid)
    if err == nil {
        proc.Kill()
        fmt.Println("✅ Process killed successfully!")
    } else {
        fmt.Println("❌ Error: Process not found.")
    }
}

func main() {
    var rootCmd = &cobra.Command{Use: "gosys"}

    rootCmd.AddCommand(&cobra.Command{
        Use:   "install [package]",
        Short: "Installs a GoSys package from GitHub",
        Args:  cobra.MinimumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            installPackage(strings.TrimSpace(args[0]))
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "info",
        Short: "Displays system information",
        Run: func(cmd *cobra.Command, args []string) {
            systemInfo()
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "disk",
        Short: "Displays disk usage",
        Run: func(cmd *cobra.Command, args []string) {
            diskUsage()
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "net",
        Short: "Displays network information",
        Run: func(cmd *cobra.Command, args []string) {
            networkInfo()
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "uptime",
        Short: "Shows system uptime",
        Run: func(cmd *cobra.Command, args []string) {
            uptime()
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "ps",
        Short: "Lists running processes",
        Run: func(cmd *cobra.Command, args []string) {
            listProcesses()
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "kill [pid]",
        Short: "Kills a process by PID",
        Args:  cobra.MinimumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            killProcess(parsePID(args[0]))
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "uuid",
        Short: "Generates a UUID",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Generated UUID:", uuid.New())
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "date",
        Short: "Shows the current date and time",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Current Date & Time:", time.Now().Format(time.RFC1123))
        },
    })

    err := rootCmd.Execute()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func parsePID(pidStr string) int32 {
    var pid int32
    fmt.Sscanf(pidStr, "%d", &pid)
    return pid
}
