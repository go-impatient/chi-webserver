package sd

import (
	"fmt"
	"github.com/moocss/chi-webserver/src/pkg/render"
	"net/http"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

const (
	_  = iota             // 0
	KB = 1 << (10 * iota) // 1 << (10 * 1)
	MB                    // 1 << (10 * 2)
	GB                    // 1 << (10 * 3)
)

// HealthCheck shows `OK` as the ping-pong result.
func HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := "OK"
		render.JSON(w, message, http.StatusOK)
	}
}

// DiskCheck checks the disk usage.
func DiskCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _ := disk.Usage("/")

		usedMB := int(u.Used) / MB
		usedGB := int(u.Used) / GB
		totalMB := int(u.Total) / MB
		totalGB := int(u.Total) / GB
		usedPercent := int(u.UsedPercent)

		status := http.StatusOK
		text := "OK"

		if usedPercent >= 95 {
			status = http.StatusOK
			text = "CRITICAL"
		} else if usedPercent >= 90 {
			status = http.StatusTooManyRequests
			text = "WARNING"
		}

		message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)

		render.JSON(w, message, status)
	}
}

// CPUCheck checks the cpu usage.
func CPUCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cores, _ := cpu.Counts(false)

		a, _ := load.Avg()
		l1 := a.Load1
		l5 := a.Load5
		l15 := a.Load15

		status := http.StatusOK
		text := "OK"

		if l5 >= float64(cores-1) {
			status = http.StatusInternalServerError
			text = "CRITICAL"
		} else if l5 >= float64(cores-2) {
			status = http.StatusTooManyRequests
			text = "WARNING"
		}

		message := fmt.Sprintf("%s - Load average: %.2f, %.2f, %.2f | Cores: %d", text, l1, l5, l15, cores)
		render.JSON(w, message, status)
	}
}

// RAMCheck checks the disk usage.
func RAMCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _ := mem.VirtualMemory()

		usedMB := int(u.Used) / MB
		usedGB := int(u.Used) / GB
		totalMB := int(u.Total) / MB
		totalGB := int(u.Total) / GB
		usedPercent := int(u.UsedPercent)

		status := http.StatusOK
		text := "OK"

		if usedPercent >= 95 {
			status = http.StatusInternalServerError
			text = "CRITICAL"
		} else if usedPercent >= 90 {
			status = http.StatusTooManyRequests
			text = "WARNING"
		}

		message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
		render.JSON(w, message, status)
	}
}