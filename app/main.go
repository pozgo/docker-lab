package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

func main() {
	printHeader()
	
	if len(os.Args) < 2 {
		printUsage()
		return
	}
	
	command := os.Args[1]
	
	switch command {
	case "start":
		startLab()
	case "stop":
		stopLab()
	case "clean":
		cleanLab()
	case "status":
		showStatus()
	case "inventory":
		generateInventory()
	case "test":
		testConnectivity()
	default:
		fmt.Printf("%s Unknown command: %s\n", red("❌"), command)
		printUsage()
	}
}

func printHeader() {
	fmt.Printf("\n%s\n", bold(cyan("🧪 LAB Operations Tool")))
	fmt.Printf("%s\n", blue("═══════════════════════════"))
}

func printUsage() {
	fmt.Printf("\n%s\n", bold("Usage: ./lab <command>"))
	fmt.Printf("\n%s\n", bold("Commands:"))
	fmt.Printf("  %s     - Start the lab environment\n", green("start"))
	fmt.Printf("  %s      - Stop the lab environment\n", yellow("stop"))
	fmt.Printf("  %s     - Clean up lab containers and images\n", red("clean"))
	fmt.Printf("  %s    - Show lab status and connection details\n", blue("status"))
	fmt.Printf("  %s - Generate Ansible inventory file\n", cyan("inventory"))
	fmt.Printf("  %s      - Test SSH and Ansible connectivity\n", blue("test"))
	fmt.Println()
}

func startLab() {
	fmt.Printf("\n%s %s\n", green("🚀"), bold("Starting LAB environment..."))
	fmt.Printf("%s\n", blue("═══════════════════════════════════"))
	
	// Check if docker-compose.yml exists
	if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
		fmt.Printf("%s docker-compose.yml not found in current directory\n", red("❌"))
		return
	}
	
	// Start the lab
	fmt.Printf("%s Building and starting containers...\n", cyan("📦"))
	cmd := exec.Command("docker", "compose", "up", "-d")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		fmt.Printf("%s Failed to start lab: %v\n", red("❌"), err)
		fmt.Printf("Output: %s\n", string(output))
		return
	}
	
	fmt.Printf("%s Lab started successfully!\n", green("✅"))
	
	// Wait a moment for containers to initialize
	time.Sleep(2 * time.Second)
	
	// Show connection details
	showConnectionDetails()
}

func stopLab() {
	fmt.Printf("\n%s %s\n", yellow("🛑"), bold("Stopping LAB environment..."))
	fmt.Printf("%s\n", blue("═══════════════════════════════════"))
	
	cmd := exec.Command("docker", "compose", "down")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		fmt.Printf("%s Failed to stop lab: %v\n", red("❌"), err)
		fmt.Printf("Output: %s\n", string(output))
		return
	}
	
	fmt.Printf("%s Lab stopped successfully!\n", green("✅"))
}

func cleanLab() {
	fmt.Printf("\n%s %s\n", red("🧹"), bold("Cleaning LAB environment..."))
	fmt.Printf("%s\n", blue("══════════════════════════════════"))
	
	// Stop containers first
	fmt.Printf("%s Stopping containers...\n", yellow("🛑"))
	stopCmd := exec.Command("docker", "compose", "down")
	stopCmd.Run()
	
	// Remove lab images
	fmt.Printf("%s Removing lab images...\n", red("🗑️"))
	imagesCmd := exec.Command("docker", "images", "-q", "lab/image")
	imageOutput, err := imagesCmd.Output()
	
	if err == nil && len(strings.TrimSpace(string(imageOutput))) > 0 {
		removeCmd := exec.Command("docker", "rmi", "lab/image:latest")
		removeCmd.Run()
	}
	
	// Clean up unused Docker resources
	fmt.Printf("%s Cleaning unused Docker resources...\n", cyan("🧽"))
	pruneCmd := exec.Command("docker", "system", "prune", "-f")
	pruneCmd.Run()
	
	fmt.Printf("%s Lab environment cleaned successfully!\n", green("✅"))
}

func showStatus() {
	fmt.Printf("\n%s %s\n", blue("📊"), bold("LAB Status"))
	fmt.Printf("%s\n", blue("═══════════════════"))
	
	// Check containers
	containers := getContainers()
	if len(containers) == 0 {
		fmt.Printf("%s No lab containers running\n", yellow("⚠️"))
		fmt.Printf("\nRun %s to start the lab\n", green("./lab start"))
		return
	}
	
	// Display container status
	displayContainerTable(containers)
	
	// Show connection details
	showConnectionDetails()
}

func getContainers() []Container {
	cmd := exec.Command("docker", "ps", "--filter", "name=lab-", "--format", "{{.Names}}:{{.Status}}:{{.Ports}}")
	output, err := cmd.Output()
	
	if err != nil {
		return []Container{}
	}
	
	lines := strings.Split(string(output), "\n")
	containers := []Container{}
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue // Skip empty lines
		}
		
		parts := strings.Split(line, ":")
		if len(parts) >= 3 {
			// Rejoin ports part (it may contain colons)
			ports := strings.Join(parts[2:], ":")
			container := Container{
				Name:   strings.TrimSpace(parts[0]),
				Status: strings.TrimSpace(parts[1]),
				Ports:  strings.TrimSpace(ports),
			}
			containers = append(containers, container)
		}
	}
	
	return containers
}

type Container struct {
	Name   string
	Status string
	Ports  string
}

func displayContainerTable(containers []Container) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Container", "Status", "SSH Port", "Hostname"})
	table.SetBorder(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	
	for _, container := range containers {
		status := container.Status
		if strings.Contains(status, "Up") {
			status = green("Running")
		} else {
			status = red("Stopped")
		}
		
		// Extract SSH port
		sshPort := extractSSHPort(container.Ports)
		hostname := extractHostname(container.Name)
		
		table.Append([]string{
			container.Name,
			status,
			sshPort,
			hostname,
		})
	}
	
	table.Render()
}

func extractSSHPort(ports string) string {
	re := regexp.MustCompile(`0\.0\.0\.0:(\d+)->22/tcp`)
	matches := re.FindStringSubmatch(ports)
	if len(matches) > 1 {
		return matches[1]
	}
	return "N/A"
}

func extractHostname(containerName string) string {
	if strings.Contains(containerName, "lab-01") {
		return "lab-01"
	} else if strings.Contains(containerName, "lab-02") {
		return "lab-02"
	}
	return "unknown"
}

func showConnectionDetails() {
	fmt.Printf("\n%s %s\n", cyan("🔗"), bold("Connection Details"))
	fmt.Printf("%s\n", blue("═══════════════════════"))
	
	containers := getContainers()
	if len(containers) == 0 {
		return
	}
	
	fmt.Printf("\n%s\n", bold("SSH Connections:"))
	
	for _, container := range containers {
		if strings.Contains(container.Status, "Up") {
			sshPort := extractSSHPort(container.Ports)
			hostname := extractHostname(container.Name)
			
			if sshPort != "N/A" {
				fmt.Printf("  %s %s:\n", green("→"), bold(hostname))
				fmt.Printf("    %s ssh labuser@localhost -p %s\n", cyan("$"), sshPort)
				fmt.Printf("    %s labpass123\n", yellow("Password:"))
				fmt.Println()
			}
		}
	}
	
	fmt.Printf("%s\n", bold("Environment Variables:"))
	fmt.Printf("  %s ROOT_PASSWORD: %s\n", blue("•"), yellow("labroot123"))
	fmt.Printf("  %s USER: %s\n", blue("•"), yellow("labuser"))
	fmt.Printf("  %s USER_PASSWORD: %s\n", blue("•"), yellow("labpass123"))
	fmt.Printf("  %s SUDO: %s\n", blue("•"), yellow("true"))
	fmt.Println()
}

func generateInventory() {
	fmt.Printf("\n%s %s\n", cyan("📋"), bold("Generating Ansible Inventory"))
	fmt.Printf("%s\n", blue("═══════════════════════════════════"))
	
	// Check if containers are running
	containers := getContainers()
	if len(containers) == 0 {
		fmt.Printf("%s No lab containers running\n", yellow("⚠️"))
		fmt.Printf("Run %s to start the lab first\n", green("./lab start"))
		return
	}
	
	// Generate dynamic inventory based on running containers
	inventoryContent := generateInventoryContent(containers)
	
	// Write to file
	err := os.WriteFile("inventory.yml", []byte(inventoryContent), 0644)
	if err != nil {
		fmt.Printf("%s Failed to write inventory file: %v\n", red("❌"), err)
		return
	}
	
	fmt.Printf("%s Ansible inventory generated: %s\n", green("✅"), bold("inventory.yml"))
	fmt.Printf("\n%s\n", bold("Usage with Ansible:"))
	fmt.Printf("  %s ansible -i inventory.yml lab_nodes -m ping\n", cyan("$"))
	fmt.Printf("  %s ansible-playbook -i inventory.yml playbook.yml\n", cyan("$"))
	fmt.Println()
}

func generateInventoryContent(containers []Container) string {
	content := `---
# LAB Ansible Inventory (Generated)
# This inventory file contains all running lab containers
# Use this with Ansible to manage lab containers via SSH

all:
  children:
    lab_environment:
      children:
        lab_nodes:
          hosts:
`
	
	// Add each running container to inventory
	for _, container := range containers {
		if strings.Contains(container.Status, "Up") {
			sshPort := extractSSHPort(container.Ports)
			hostname := extractHostname(container.Name)
			
			if sshPort != "N/A" {
				content += fmt.Sprintf(`            %s:
              ansible_host: localhost
              ansible_port: %s
              ansible_user: labuser
              ansible_ssh_pass: labpass123
              ansible_ssh_common_args: '-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null'
              container_name: %s
              hostname: %s
              ssh_port: %s
              
`, hostname, sshPort, container.Name, hostname, sshPort)
			}
		}
	}
	
	content += `          vars:
            # Common variables for all lab nodes
            ansible_python_interpreter: /usr/bin/python3
            lab_root_password: labroot123
            lab_user_password: labpass123
            lab_sudo_enabled: true
            lab_environment: lab
            
    # Group variables for specific roles
    ubuntu_nodes:
      children:
        lab_nodes:
      vars:
        ansible_os_family: Debian
        ansible_distribution: Ubuntu
        ansible_distribution_version: "22.04"
`
	
	return content
}

func testConnectivity() {
	fmt.Printf("\n%s %s\n", blue("🧪"), bold("Testing LAB Connectivity"))
	fmt.Printf("%s\n", blue("═══════════════════════════════════"))
	
	// Check if containers are running
	containers := getContainers()
	if len(containers) == 0 {
		fmt.Printf("%s No lab containers running\n", yellow("⚠️"))
		fmt.Printf("Run %s to start the lab first\n", green("./lab start"))
		return
	}
	
	fmt.Printf("\n%s\n", bold("SSH Connectivity Tests:"))
	
	allPassed := true
	for _, container := range containers {
		if strings.Contains(container.Status, "Up") {
			sshPort := extractSSHPort(container.Ports)
			hostname := extractHostname(container.Name)
			
			if sshPort != "N/A" {
				fmt.Printf("  %s Testing %s (port %s)... ", blue("→"), bold(hostname), sshPort)
				
				// Test SSH connection with timeout
				cmd := exec.Command("timeout", "10", "sshpass", "-p", "labpass123", "ssh", 
					"-o", "StrictHostKeyChecking=no", 
					"-o", "UserKnownHostsFile=/dev/null",
					"-o", "ConnectTimeout=5",
					"-p", sshPort, "labuser@localhost", "echo 'SSH_OK'")
				
				output, err := cmd.Output()
				if err != nil || !strings.Contains(string(output), "SSH_OK") {
					fmt.Printf("%s\n", red("FAILED"))
					allPassed = false
				} else {
					fmt.Printf("%s\n", green("PASSED"))
				}
			}
		}
	}
	
	// Test Ansible if available
	fmt.Printf("\n%s\n", bold("Ansible Connectivity Tests:"))
	
	// Check if ansible is installed
	_, err := exec.LookPath("ansible")
	if err != nil {
		fmt.Printf("  %s Ansible not installed - skipping Ansible tests\n", yellow("⚠️"))
		fmt.Printf("  %s Install with: sudo apt install ansible\n", cyan("💡"))
	} else {
		// Generate inventory for testing
		inventoryContent := generateInventoryContent(containers)
		err := os.WriteFile("inventory-test.yml", []byte(inventoryContent), 0644)
		if err != nil {
			fmt.Printf("  %s Failed to create test inventory\n", red("❌"))
		} else {
			fmt.Printf("  %s Testing Ansible ping... ", blue("→"))
			
			cmd := exec.Command("ansible", "-i", "inventory-test.yml", "lab_nodes", "-m", "ping")
			output, err := cmd.CombinedOutput()
			
			if err != nil {
				fmt.Printf("%s\n", red("FAILED"))
				fmt.Printf("    %s\n", string(output))
				allPassed = false
			} else if strings.Contains(string(output), "SUCCESS") {
				fmt.Printf("%s\n", green("PASSED"))
			} else {
				fmt.Printf("%s\n", yellow("PARTIAL"))
				allPassed = false
			}
			
			// Clean up test inventory
			os.Remove("inventory-test.yml")
		}
	}
	
	fmt.Printf("\n%s\n", bold("Test Summary:"))
	if allPassed {
		fmt.Printf("  %s All connectivity tests passed!\n", green("✅"))
	} else {
		fmt.Printf("  %s Some tests failed - check SSH configuration\n", red("❌"))
	}
	fmt.Println()
}