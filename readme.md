# USB Loader

A web-based tool to view disk information and mount USB drives using Go and HTML with Material-UI styling.

## Features

- Display disk and partition information from `fdisk -l`
- Mount disks to specified directories through a web interface
- Responsive Material-UI design
- RESTful API endpoints

## Requirements

- Go 1.16 or higher
- Linux system with `fdisk` and `mount` commands available
- Root/sudo privileges for mounting operations

## Installation

1. Clone or download this repository
2. Navigate to the project directory
3. Initialize Go modules:
   ```bash
   go mod init usb-loader
   go mod tidy
   ```

## Usage

1. Run the application with appropriate privileges:
   ```bash
   # For full functionality (mounting), run as root or with sudo
   sudo go run main.go

   # For testing without mounting capabilities
   go run main.go
   ```

2. Open your web browser and navigate to `http://localhost:8080`

3. The application will display available disks and partitions

4. Click "Mount" next to any device to specify a mount path and mount the device

## API Endpoints

### GET /api/disks
Returns disk and partition information from `fdisk -l`

**Response:**
```json
{
  "success": true,
  "message": "Disk information retrieved successfully",
  "data": [
    {
      "device": "/dev/sda",
      "size": "1000.2 GB",
      "type": "Disk"
    },
    {
      "device": "/dev/sda1",
      "size": "100M 7",
      "type": "Partition"
    }
  ]
}
```

### POST /api/mount
Mounts a device to a specified path

**Request:**
```json
{
  "device": "/dev/sda1",
  "path": "/mnt/usb"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Disk mounted successfully"
}
```

## Security Notes

- This application executes system commands and should only be run in trusted environments
- Mount operations require root privileges
- Consider running behind authentication in production environments
- The application listens on all interfaces by default; consider binding to localhost only for local use

## Frontend

The frontend is built with pure HTML, CSS, and JavaScript using Material Design principles:
- Responsive layout that works on desktop and mobile
- Material Icons for visual elements
- Clean, modern interface with card-based design
- Real-time feedback for user actions

## Troubleshooting

- **No disks shown**: Ensure `fdisk -l` works in your terminal and returns data
- **Mount fails**: Verify you have appropriate permissions and the mount directory exists
- **Server won't start**: Check that port 8080 is available and Go is properly installed

## License

This project is for educational and personal use only.