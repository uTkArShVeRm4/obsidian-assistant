#!/bin/sh

# Create a temporary askpass script
echo "#!/bin/sh\necho \$GIT_ACCESS_TOKEN" >/git-askpass.sh
chmod +x /git-askpass.sh

# Tell Git to use the askpass script
export GIT_ASKPASS=/git-askpass.sh

# Set up Git config
git config --global user.name "$GIT_USERNAME"
git config --global user.email "$GIT_EMAIL"

# Set remote URL without exposing credentials
git -C /app/data/Main remote set-url origin "https://github.com/$GIT_USERNAME/Main.git"

# Run the main application
exec ./main
