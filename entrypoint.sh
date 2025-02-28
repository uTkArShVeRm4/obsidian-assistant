#!/bin/sh

# Set up git configuration if credentials are provided
if [ -n "$GIT_USERNAME" ] && [ -n "$GIT_EMAIL" ]; then
    git config --global user.name "$GIT_USERNAME"
    git config --global user.email "$GIT_EMAIL"

    # If using a personal access token
    if [ -n "$GIT_ACCESS_TOKEN" ]; then
        # Configure git to store credentials
        git config --global credential.helper store

        # Create or update .git-credentials file
        echo "https://$GIT_USERNAME:$GIT_ACCESS_TOKEN@github.com" >/app/data/.git-credentials
        git config --global credential.helper "store --file=/app/data/.git-credentials"

        # Set file permissions
        chmod 600 /app/data/.git-credentials
    fi

    # Additional configurations
    git config --global --add safe.directory /app/data
fi

# Execute the main application
exec ./main
