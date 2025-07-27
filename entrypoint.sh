#!/bin/bash

# LAB Container Entrypoint Script
# Processes environment variables and configures the container

set -e

echo "=========================================="
echo "LAB Container Starting..."
echo "=========================================="

# Function to log environment variables
log_env_var() {
    local var_name="$1"
    local var_value="$2"
    local is_sensitive="$3"
    
    if [ "$is_sensitive" = "true" ]; then
        echo "🔧 $var_name: [REDACTED - LENGTH: ${#var_value}]"
    else
        echo "🔧 $var_name: $var_value"
    fi
}

# Log all environment variables for LAB transparency
echo "📋 Environment Variables:"
log_env_var "ROOT_PASSWORD" "$ROOT_PASSWORD" "true"
log_env_var "USER" "$USER" "false"
log_env_var "USER_PASSWORD" "$USER_PASSWORD" "true"
log_env_var "SUDO" "$SUDO" "false"

echo ""
echo "🔧 Container Configuration:"

# 1. Configure root password
if [ -n "$ROOT_PASSWORD" ]; then
    echo "root:$ROOT_PASSWORD" | chpasswd
    echo "✅ Root password configured"
else
    echo "⚠️  No root password set - root login disabled"
fi

# 2. Create user if specified
if [ -n "$USER" ]; then
    echo "👤 Creating user: $USER"
    
    # Check if user already exists
    if id "$USER" &>/dev/null; then
        echo "ℹ️  User $USER already exists"
    else
        # Create user with home directory
        useradd -m -s /bin/bash "$USER"
        echo "✅ User $USER created"
    fi
    
    # Set user password if provided
    if [ -n "$USER_PASSWORD" ]; then
        echo "$USER:$USER_PASSWORD" | chpasswd
        echo "✅ Password set for user $USER"
    else
        echo "⚠️  No password set for user $USER"
    fi
    
    # Add to sudo group if SUDO is true
    if [ "$SUDO" = "true" ] || [ "$SUDO" = "TRUE" ] || [ "$SUDO" = "1" ]; then
        usermod -aG sudo "$USER"
        echo "✅ User $USER added to sudo group"
        
        # Allow passwordless sudo for lab convenience
        echo "$USER ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/lab-$USER
        chmod 440 /etc/sudoers.d/lab-$USER
        echo "✅ Passwordless sudo configured for $USER"
    else
        echo "ℹ️  User $USER not granted sudo privileges"
    fi
else
    echo "ℹ️  No additional user specified"
fi

# 3. Generate SSH host keys if they don't exist
if [ ! -f /etc/ssh/ssh_host_rsa_key ]; then
    echo "🔑 Generating SSH host keys..."
    ssh-keygen -A
    echo "✅ SSH host keys generated"
else
    echo "ℹ️  SSH host keys already exist"
fi

# 4. Configure SSH daemon
echo "🌐 Configuring SSH daemon..."

# Allow root login if root password is set
if [ -n "$ROOT_PASSWORD" ]; then
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
    sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
    echo "✅ Root SSH login enabled"
else
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config
    sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
    echo "✅ Root SSH login disabled"
fi

# Enable password authentication for lab use
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config
sed -i 's/PasswordAuthentication no/PasswordAuthentication yes/' /etc/ssh/sshd_config

echo ""
echo "🚀 Starting SystemCtl Replacement..."
echo "=========================================="

# If no command provided, start systemctl replacement
if [ $# -eq 0 ]; then
    exec /usr/local/bin/systemctl
else
    exec "$@"
fi