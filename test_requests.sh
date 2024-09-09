#!/bin/bash

# Test URLs for POST and GET requests
BASE_URL="http://localhost:8082"

# Array of destinations to shorten with more variety and longer paths
declare -a destinations=(
    "https://github.com/some/repo/path/with/multiple/levels"
    "https://golang.org/doc/install/source#installing-from-source"
    "https://www.example.com/product/category/item?id=12345"
    "https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Functions#defining_functions"
    "https://www.wikipedia.org/wiki/Special:Random"
    "https://www.longurl.com/this/is/a/very/long/url/example/that/i/want/to/shorten"
    "https://stackoverflow.com/questions/12345678/how-to-use-bash-in-linux"
    "https://www.google.com/search?q=bash+script+examples&hl=en"
    "https://www.reddit.com/r/programming/comments/a1b2c3/interesting_article_about_bash_scripting/"
    "https://www.microsoft.com/en-us/software-download/windows10"
)

# Initialize counters
declare -i post_request_count=0
declare -i get_request_count=0
declare -i error_count=0

# Array to hold the short URLs
short_urls=()

# Loop to send 100 POST requests
for ((i=0; i<100; i++)); do
    dest=${destinations[$((i % ${#destinations[@]}))]}
    echo "Saving long URL: $dest"
    
    # Send POST request and capture the response
    response=$(curl -s -X POST -H "Content-Type: application/json" -d "{\"destination\": \"$dest\"}" "$BASE_URL/")
    
    # Extract the short URL from the response JSON
    short_url=$(echo "$response" | grep -o '"shortUrl":"[^"]*' | sed 's/"shortUrl":"//')

    if [ "$short_url" == "null" ]; then
        echo "Failed to generate short URL for: $dest"
        echo "Response: $response"
        ((error_count++))  # Increment error counter
        continue
    fi

    echo "Response: $response"
    echo "Short URL generated: $short_url"
    
    # Increment POST request counter
    ((post_request_count++))
    
    # Store the short URL for later use
    short_urls+=("$short_url")
done

# GET requests to retrieve long URLs using the short URLs generated in the previous step
for short_url in "${short_urls[@]}"; do
    echo "Retrieving long URL from short URL: $short_url"
    
    # Send GET request to the short URL
    long_url=$(curl -s -L -o /dev/null -w '%{url_effective}' "$short_url")

    if [ $? -ne 0 ]; then
        echo "Error retrieving long URL for short URL: $short_url"
        ((error_count++))  # Increment error counter
        continue
    fi

    # Display the long URL
    echo "Redirected to: $long_url"
    
    # Increment GET request counter
    ((get_request_count++))
done

# Summary of requests made
echo ""
echo "Summary of requests made:"
echo "Total POST requests: $post_request_count"
echo "Total GET requests: $get_request_count"
echo "Total errors: $error_count"
