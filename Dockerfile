FROM ubuntu:22.04

# Prevent interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install systemd, SSH, and other essential packages
RUN apt-get update && apt-get install -y \
    systemd \
    systemd-sysv \
    openssh-server \
    sudo \
    curl \
    wget \
    vim \
    nano \
    net-tools \
    iputils-ping \
    python3 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Download and install docker-systemctl-replacement
RUN curl -fsSL https://raw.githubusercontent.com/gdraheim/docker-systemctl-replacement/master/files/docker/systemctl3.py -o /usr/local/bin/systemctl \
    && chmod +x /usr/local/bin/systemctl \
    && ln -sf /usr/local/bin/systemctl /usr/bin/systemctl \
    && ln -sf /usr/local/bin/systemctl /bin/systemctl

# Configure SSH
RUN mkdir -p /var/run/sshd \
    && mkdir -p /root/.ssh \
    && chmod 700 /root/.ssh

# Enable SSH service using systemctl replacement
RUN systemctl enable ssh

# Copy entrypoint script
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Environment variables for configuration
ENV ROOT_PASSWORD=""
ENV USER=""
ENV USER_PASSWORD=""
ENV SUDO=""

# Expose SSH port
EXPOSE 22

# Use custom entrypoint that processes env vars and starts systemctl replacement
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Start systemctl replacement as PID 1
CMD ["/usr/local/bin/systemctl"]