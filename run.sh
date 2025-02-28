sudo docker run -d -p 7777:7777 \
    -v $(pwd)/data:/app/data itt \
    -e GIT_USERNAME="your-username" \
    -e GIT_EMAIL="your-email@example.com" \
    -e GIT_ACCESS_TOKEN="your-personal-access-token"
