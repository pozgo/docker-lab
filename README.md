# LAB 🧪

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS-lightgrey.svg)
![Docker](https://img.shields.io/badge/docker-required-blue.svg)
![Go](https://img.shields.io/badge/go-1.21%2B-00ADD8.svg)

**LAB** is a containerized laboratory environment built with Docker and Docker Compose, designed for educational purposes, testing, and development. It provides isolated Ubuntu 22.04 containers with SystemD support, SSH access, and Ansible integration.

## ✨ Features

- 🐧 **Ubuntu 22.04 LTS** containers with full SystemD support
- 🔧 **SystemD Services** - Enable and manage custom services with `systemctl`
- 🔐 **SSH Access** - Direct SSH connectivity to each container
- 🎯 **Multi-Container** - Easy scaling with Docker Compose
- 📋 **Ansible Ready** - Built-in inventory generation and connectivity testing
- 🛠️ **Cross-Platform** - Binaries for Linux (AMD64/ARM64) and macOS (Intel/M1/M2)
- 🎨 **Beautiful CLI** - Colorized output with intuitive commands
- 🧹 **Clean Management** - Safe cleanup without affecting other Docker resources

## 🚀 Quick Start

### Prerequisites

- **Docker** and **Docker Compose** installed
- **Go 1.21+** (for building from source)
- **Ansible** (optional, for automation features)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/pozgo/docker-lab.git
   cd LAB
   ```

2. **Download pre-built binary** (recommended):
   ```bash
   # For Linux AMD64
   wget https://github.com/pozgo/docker-lab/releases/latest/download/lab-linux-amd64
   chmod +x lab-linux-amd64
   mv lab-linux-amd64 lab
   
   # For macOS ARM64 (M1/M2)
   wget https://github.com/pozgo/docker-lab/releases/latest/download/lab-darwin-arm64
   chmod +x lab-darwin-arm64
   mv lab-darwin-arm64 lab
   ```

3. **Or build from source:**
   ```bash
   cd app
   go build -o ../lab .
   ```

### Basic Usage

```bash
# Start the lab environment (default: 2 containers)
./lab start

# Start with custom number of containers
./lab start --containers 5
./lab start -c 3

# Check status and connection details
./lab status

# Connect via SSH (ports start from 2222)
ssh labuser@localhost -p 2222  # lab-01
ssh labuser@localhost -p 2223  # lab-02
ssh labuser@localhost -p 2224  # lab-03 (if created)

# Stop the lab
./lab stop

# Clean up everything
./lab clean
```

## 📖 Documentation

### Commands

| Command | Description |
|---------|-------------|
| `start [--containers N]` | Start the lab environment with N containers (default: 2) |
| `stop` | Stop the lab environment |
| `clean` | Clean up lab containers and images |
| `status` | Show lab status and connection details |
| `inventory` | Generate Ansible inventory file |
| `test` | Test SSH and Ansible connectivity |

### Container Architecture

```
LAB Environment
├── lab-01 (lab-01)
│   ├── SSH: localhost:2222
│   ├── Hostname: lab-01
│   └── SystemD: ✅ Active
└── lab-02 (lab-02)
    ├── SSH: localhost:2223
    ├── Hostname: lab-02
    └── SystemD: ✅ Active
```

### Default Credentials

| Component | Username | Password |
|-----------|----------|----------|
| SSH User | `labuser` | `labpass123` |
| Root User | `root` | `labroot123` |
| Sudo Access | `labuser` | ✅ Enabled |

> ⚠️ **Security Note**: These are default credentials for lab use only. Change them for production environments.

## 🔧 Advanced Usage

### SystemD Services

Create and manage custom SystemD services:

```bash
# Connect to a container
ssh labuser@localhost -p 2222

# Create a custom service
sudo tee /etc/systemd/system/myapp.service > /dev/null <<EOF
[Unit]
Description=My Application
After=network.target

[Service]
Type=simple
ExecStart=/usr/bin/python3 -m http.server 8080
Restart=always
User=labuser

[Install]
WantedBy=multi-user.target
EOF

# Enable and start the service
sudo systemctl enable myapp.service --now
sudo systemctl status myapp.service
```

### Ansible Integration

#### Generate Dynamic Inventory

```bash
# Generate inventory based on running containers
./lab inventory

# Test connectivity
./lab test
```

#### Use with Ansible

```bash
# Ping all lab nodes
ansible -i inventory.yml lab_nodes -m ping

# Run ad-hoc commands
ansible -i inventory.yml lab_nodes -m shell -a "uptime"

# Execute playbooks
ansible-playbook -i inventory.yml your-playbook.yml
```

#### Example Playbook

```yaml
---
- name: Configure Lab Environment
  hosts: lab_nodes
  become: yes
  tasks:
    - name: Install packages
      apt:
        name:
          - htop
          - curl
          - vim
        state: present
        update_cache: yes
    
    - name: Create application directory
      file:
        path: /opt/myapp
        state: directory
        owner: labuser
        group: labuser
    
    - name: Deploy custom service
      template:
        src: myapp.service.j2
        dest: /etc/systemd/system/myapp.service
      notify: restart myapp
    
  handlers:
    - name: restart myapp
      systemd:
        name: myapp
        state: restarted
        enabled: yes
        daemon_reload: yes
```

### Persistent Storage

Data persistence is handled through Docker volumes:

- **Home Directories**: `/home` (persistent across container restarts)
- **SystemD Services**: `/etc/systemd/system` (custom services persist)

```bash
# View volumes
docker volume ls | grep lab

# Backup home directory
docker run --rm -v lab-01-home:/data -v $(pwd):/backup alpine tar czf /backup/lab-01-home.tar.gz -C /data .
```

## 🏗️ Development

### Building All Platforms

```bash
# Use the provided Makefile
make build-all

# Or manually
cd app
GOOS=linux GOARCH=amd64 go build -o ../lab-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o ../lab-linux-arm64 .
GOOS=darwin GOARCH=amd64 go build -o ../lab-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o ../lab-darwin-arm64 .
```

### Project Structure

```
LAB/
├── app/                    # Go application source
│   ├── main.go            # Main application logic
│   ├── go.mod             # Go module definition
│   └── go.sum             # Go dependencies
├── Dockerfile             # Container image definition
├── docker-compose.yml     # Multi-container orchestration
├── entrypoint.sh          # Container startup script
├── inventory.yml          # Ansible inventory
├── Makefile              # Build automation
└── README.md             # This file
```

### Environment Variables

Customize the environment by modifying `docker-compose.yml`:

```yaml
environment:
  - ROOT_PASSWORD=your_root_password
  - USER=your_username
  - USER_PASSWORD=your_password
  - SUDO=true
```

## 🧪 Use Cases

### Educational & Training

- **Linux Administration**: Practice SystemD, networking, and package management
- **Container Orchestration**: Learn Docker Compose and multi-container deployments
- **Configuration Management**: Experiment with Ansible playbooks and automation
- **SSH & Networking**: Understand port forwarding and remote access

### Development & Testing

- **Service Development**: Test SystemD services in isolated environments
- **CI/CD Pipelines**: Use as testing infrastructure for automation workflows
- **Network Simulation**: Create complex multi-node scenarios
- **Security Testing**: Practice in safe, isolated containers

### DevOps & SRE

- **Infrastructure as Code**: Develop and test Ansible automation
- **Monitoring Setup**: Deploy monitoring stacks across multiple nodes
- **Load Balancing**: Test service discovery and load balancing configurations
- **Disaster Recovery**: Practice backup and restore procedures

## 🔍 Troubleshooting

### Common Issues

#### Containers Won't Start

```bash
# Check Docker daemon status
sudo systemctl status docker

# View container logs
docker logs lab-01
docker logs lab-02

# Check available resources
docker system df
```

#### SSH Connection Failed

```bash
# Verify containers are running
./lab status

# Test connectivity
./lab test

# Manual SSH test
ssh -o StrictHostKeyChecking=no -p 2222 labuser@localhost
```

#### Port Conflicts

If ports 2222 or 2223 are in use, modify `docker-compose.yml`:

```yaml
ports:
  - "2224:22"  # Change to available port
```

#### Ansible Issues

```bash
# Install Ansible
sudo apt install ansible  # Ubuntu/Debian
brew install ansible      # macOS

# Verify inventory
ansible-inventory -i inventory.yml --list

# Test with verbose output
ansible -i inventory.yml lab_nodes -m ping -vvv
```

### Performance Tips

- **Resource Allocation**: Ensure sufficient RAM (minimum 2GB recommended)
- **Storage**: Use SSD storage for better I/O performance
- **Networking**: Avoid port conflicts with other services
- **Cleanup**: Regular cleanup prevents resource exhaustion

```bash
# Monitor resource usage
docker stats

# Clean up unused resources
docker system prune -a
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Coding Standards

- **Go**: Follow standard Go conventions and use `gofmt`
- **Docker**: Use multi-stage builds and minimize image size
- **Documentation**: Update README.md for new features
- **Testing**: Ensure all functionality works across platforms

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Ubuntu Team** for the excellent base image
- **Docker Community** for containerization technology
- **Ansible Project** for automation capabilities
- **Go Team** for the robust programming language

## 📞 Support

- 📖 **Documentation**: Check this README and inline help (`./lab`)
- 🐛 **Issues**: Report bugs via [GitHub Issues](https://github.com/pozgo/docker-lab/issues)
- 💬 **Discussions**: Join our [GitHub Discussions](https://github.com/pozgo/docker-lab/discussions)
- 📧 **Security**: Report security issues via email to security@yourorg.com

---

<div align="center">

**Made with ❤️ for the DevOps and Education Community**

[⭐ Star this repository](https://github.com/pozgo/docker-lab) if you find it helpful!

</div>