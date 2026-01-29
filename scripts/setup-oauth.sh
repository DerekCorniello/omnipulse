#!/bin/bash
#
# OmniPulse OAuth Setup Script
# Helps set up OAuth credentials for YouTube, X, and LinkedIn APIs
#

set -e

echo "========================================"
echo "OmniPulse OAuth Setup"
echo "========================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env exists
if [ ! -f .env ]; then
    if [ -f .env.example ]; then
        echo -e "${YELLOW}Creating .env from .env.example...${NC}"
        cp .env.example .env
    else
        echo -e "${RED}Error: .env.example not found${NC}"
        exit 1
    fi
fi

echo "This script will guide you through setting up API credentials."
echo "Press Enter to continue or Ctrl+C to cancel."
read

# =============================================================================
# YouTube Setup
# =============================================================================
echo ""
echo -e "${GREEN}=== YouTube API Setup ===${NC}"
echo ""
echo "1. Go to: https://console.cloud.google.com/"
echo "2. Create a new project or select an existing one"
echo "3. Enable the following APIs:"
echo "   - YouTube Data API v3"
echo "   - YouTube Analytics API"
echo "4. Go to 'Credentials' and create:"
echo "   - An API Key (for public data)"
echo "   - An OAuth 2.0 Client ID (for private data/analytics)"
echo ""
echo "For OAuth 2.0:"
echo "   - Application type: Web application"
echo "   - Authorized redirect URIs: http://localhost:8080/oauth/youtube/callback"
echo ""
read -p "Have you completed the YouTube setup? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "Enter your YouTube API Key: " YOUTUBE_API_KEY
    read -p "Enter your YouTube Client ID: " YOUTUBE_CLIENT_ID
    read -p "Enter your YouTube Client Secret: " YOUTUBE_CLIENT_SECRET
    read -p "Enter your YouTube Channel ID: " YOUTUBE_CHANNEL_ID

    # Update .env file
    sed -i "s/YOUTUBE_API_KEY=.*/YOUTUBE_API_KEY=$YOUTUBE_API_KEY/" .env
    sed -i "s/YOUTUBE_CLIENT_ID=.*/YOUTUBE_CLIENT_ID=$YOUTUBE_CLIENT_ID/" .env
    sed -i "s/YOUTUBE_CLIENT_SECRET=.*/YOUTUBE_CLIENT_SECRET=$YOUTUBE_CLIENT_SECRET/" .env
    sed -i "s/YOUTUBE_CHANNEL_ID=.*/YOUTUBE_CHANNEL_ID=$YOUTUBE_CHANNEL_ID/" .env

    echo -e "${GREEN}YouTube credentials saved!${NC}"
fi

# =============================================================================
# X (Twitter) Setup
# =============================================================================
echo ""
echo -e "${GREEN}=== X (Twitter) API Setup ===${NC}"
echo ""
echo "1. Go to: https://developer.twitter.com/en/portal/dashboard"
echo "2. Create a project and app (Basic access is sufficient)"
echo "3. From the 'Keys and tokens' tab, get:"
echo "   - API Key and Secret"
echo "   - Bearer Token"
echo "   - Access Token and Secret (generate if needed)"
echo ""
echo "4. Get your User ID:"
echo "   - Go to: https://tweeterid.com/"
echo "   - Enter your X username to get your numeric user ID"
echo ""
read -p "Have you completed the X setup? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "Enter your X API Key: " X_API_KEY
    read -p "Enter your X API Secret: " X_API_SECRET
    read -p "Enter your X Bearer Token: " X_BEARER_TOKEN
    read -p "Enter your X Access Token: " X_ACCESS_TOKEN
    read -p "Enter your X Access Secret: " X_ACCESS_SECRET
    read -p "Enter your X User ID: " X_USER_ID

    sed -i "s/X_API_KEY=.*/X_API_KEY=$X_API_KEY/" .env
    sed -i "s/X_API_SECRET=.*/X_API_SECRET=$X_API_SECRET/" .env
    sed -i "s/X_BEARER_TOKEN=.*/X_BEARER_TOKEN=$X_BEARER_TOKEN/" .env
    sed -i "s/X_ACCESS_TOKEN=.*/X_ACCESS_TOKEN=$X_ACCESS_TOKEN/" .env
    sed -i "s/X_ACCESS_SECRET=.*/X_ACCESS_SECRET=$X_ACCESS_SECRET/" .env
    sed -i "s/X_USER_ID=.*/X_USER_ID=$X_USER_ID/" .env

    echo -e "${GREEN}X credentials saved!${NC}"
fi

# =============================================================================
# LinkedIn Setup
# =============================================================================
echo ""
echo -e "${GREEN}=== LinkedIn API Setup ===${NC}"
echo ""
echo "1. Go to: https://www.linkedin.com/developers/apps"
echo "2. Create a new app"
echo "3. Add products to your app:"
echo "   - 'Sign In with LinkedIn using OpenID Connect'"
echo "   - 'Share on LinkedIn'"
echo "   Note: Some analytics features require additional approval"
echo ""
echo "4. In the 'Auth' tab:"
echo "   - Copy Client ID and Client Secret"
echo "   - Add redirect URL: http://localhost:8080/oauth/linkedin/callback"
echo ""
echo "5. Required OAuth scopes:"
echo "   - openid, profile, email (from Sign In)"
echo "   - w_member_social (from Share on LinkedIn)"
echo ""
read -p "Have you completed the LinkedIn setup? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "Enter your LinkedIn Client ID: " LINKEDIN_CLIENT_ID
    read -p "Enter your LinkedIn Client Secret: " LINKEDIN_CLIENT_SECRET

    sed -i "s/LINKEDIN_CLIENT_ID=.*/LINKEDIN_CLIENT_ID=$LINKEDIN_CLIENT_ID/" .env
    sed -i "s/LINKEDIN_CLIENT_SECRET=.*/LINKEDIN_CLIENT_SECRET=$LINKEDIN_CLIENT_SECRET/" .env

    echo -e "${GREEN}LinkedIn credentials saved!${NC}"
    echo ""
    echo -e "${YELLOW}Note: To complete LinkedIn setup, you'll need to:${NC}"
    echo "1. Start the app: go run ./cmd/omnipulse serve"
    echo "2. Visit: http://localhost:8080/oauth/linkedin/authorize"
    echo "3. Complete the OAuth flow to get your access token"
fi

# =============================================================================
# Summary
# =============================================================================
echo ""
echo "========================================"
echo -e "${GREEN}Setup Complete!${NC}"
echo "========================================"
echo ""
echo "Next steps:"
echo "1. Review your .env file to ensure credentials are correct"
echo "2. Build the app: go build -o bin/omnipulse ./cmd/omnipulse"
echo "3. Run migrations: ./bin/omnipulse migrate"
echo "4. Start the server: ./bin/omnipulse serve"
echo ""
echo -e "${YELLOW}Remember: Never commit your .env file to version control!${NC}"
