#!/bin/sh

# Set up git configuration if credentials are provided
if [ -n "$GIT_USERNAME" ] && [ -n "$GIT_EMAIL" ]; then
    git config --global user.name "$GIT_USERNAME"
    git config --global user.email "$GIT_EMAIL"

    if [ -n "$GIT_ACCESS_TOKEN" ]; then
        git config --global credential.helper store
        echo "https://$GIT_USERNAME:$GIT_ACCESS_TOKEN@github.com" >/root/.git-credentials
        git config --global credential.helper "store --file=/root/.git-credentials"
        chmod 600 /root/.git-credentials
    fi

    git config --global --add safe.directory /app/data/Main
fi

# If remote URL is provided, add it
if [ -n "$GIT_REMOTE_URL" ]; then
    git -C /app/data/Main remote add origin "$GIT_REMOTE_URL"
fi

# Fix ownership issues
chown -R root:root /app/data/Main/.git

# Execute the main application
exec ./main
