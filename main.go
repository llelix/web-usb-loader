package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// DiskInfo represents information about a disk device
type DiskInfo struct {
	Device     string `json:"device"`
	Size       string `json:"size"`
	Type       string `json:"type"`
	MountPoint string `json:"mountPoint"`
}

// Response structure for API responses
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// Set up routes
	http.HandleFunc("/api/disks", getDisksHandler)
	http.HandleFunc("/api/mount", mountDiskHandler)

	// Serve static files (HTML, CSS, JS)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("USB Loader server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// getDisksHandler handles requests to get disk information
func getDisksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	disks, err := getDiskInfo()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to get disk information: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Disk information retrieved successfully",
		Data:    disks,
	})
}

// mountDiskHandler handles requests to mount a disk
func mountDiskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Device string `json:"device"`
		Path   string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.Device == "" || req.Path == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Device and path are required",
		})
		return
	}

	// Ensure the mount directory exists
	if err := exec.Command("mkdir", "-p", req.Path).Run(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to create mount directory: " + err.Error(),
		})
		return
	}

	// Execute mount command
	cmd := exec.Command("mount", req.Device, req.Path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Mount failed: " + string(output),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Disk mounted successfully",
	})
}

// getDiskInfo executes fdisk -l and parses the output
func getDiskInfo() ([]DiskInfo, error) {
	cmd := exec.Command("fdisk", "-l")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute fdisk: %w", err)
	}

	return parseFdiskOutput(string(output)), nil
}

// parseFdiskOutput parses the output of fdisk -l command
func parseFdiskOutput(output string) []DiskInfo {
	var disks []DiskInfo
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Disk /dev/") && !strings.Contains(line, "Disklabel") {
			// Parse line like: "Disk /dev/sda: 1000.2 GB, 1000204886016 bytes, 1953525168 sectors"
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				device := strings.TrimSpace(strings.TrimPrefix(parts[0], "Disk"))
				sizeInfo := strings.TrimSpace(parts[1])
				// Extract size (e.g., "1000.2 GB")
				sizeParts := strings.Split(sizeInfo, ",")
				size := sizeParts[0]

				disks = append(disks, DiskInfo{
					Device: device,
					Size:   strings.TrimSpace(size),
					Type:   "Disk",
				})
			}
		} else if strings.HasPrefix(line, "/dev/") && strings.Contains(line, "*") {
			// Parse partition lines like: "/dev/sda1   *        2048   206847   204800   100M  7 HPFS/NTFS/exFAT"
			parts := strings.Fields(line)
			if len(parts) >= 6 {
				device := parts[0]
				size := parts[len(parts)-2] + " " + parts[len(parts)-1]

				disks = append(disks, DiskInfo{
					Device: device,
					Size:   size,
					Type:   "Partition",
				})
			}
		}
	}

	return disks
}