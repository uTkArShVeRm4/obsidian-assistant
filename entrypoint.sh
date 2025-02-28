#!/bin/sh

# Ensure SSH agent is running
eval $(ssh-agent -s)

# Add SSH key (mounted from the host)
ssh-add /root/.ssh/id_ed25519

# Configure Git
git config --global user.name "$GIT_USERNAME"
git config --global user.email "$GIT_EMAIL"
git config --global --add safe.directory /app/data/Main

# Make sure the repo is using SSH instead of HTTPS
if [ -n "$GIT_REMOTE_URL" ]; then
    git -C /app/data/Main remote set-url origin "$GIT_REMOTE_URL"
fi

# Execute the main application
exec ./main
