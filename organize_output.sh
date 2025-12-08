#!/bin/bash

# Create directories if they don't exist
mkdir -p visual-briefs
mkdir -p sketchnotes

# Move Visual Brief markdown files
if ls Visual_Brief_*.md 1> /dev/null 2>&1; then
    mv Visual_Brief_*.md visual-briefs/
    echo "Moved Visual Briefs to visual-briefs/"
else
    echo "No Visual Briefs found in root."
fi

# Move PNG files
# Moves all .png files from root to sketchnotes/
if ls *.png 1> /dev/null 2>&1; then
    mv *.png sketchnotes/
    echo "Moved PNGs to sketchnotes/"
else
    echo "No PNGs found in root."
fi
