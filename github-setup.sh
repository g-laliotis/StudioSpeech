#!/bin/bash

# GitHub repository setup script
# Run this after creating the repository on GitHub

echo "Setting up GitHub remote..."

# Replace 'yourusername' with your actual GitHub username
read -p "Enter your GitHub username: " username

# Add the remote origin
git remote add origin https://github.com/$username/StudioSpeech.git

# Set the default branch name
git branch -M main

# Push to GitHub
git push -u origin main

echo "Repository pushed to GitHub successfully!"
echo "Visit: https://github.com/$username/StudioSpeech"