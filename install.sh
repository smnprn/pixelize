#!/bin/bash

APP_NAME="pixelize"
ARCHIVE_NAME="pixelize-linux-amd64.tar.gz"
BINARY_NAME="pixelize-linux-amd64"
DOWNLOAD_URL="https://github.com/smnprn/pixelize/releases/download/v2.0.0/${ARCHIVE_NAME}"

# Download the binary
echo "Downloading ${ARCHIVE_NAME} from ${DOWNLOAD_URL}..."
wget -O "${ARCHIVE_NAME}" "${DOWNLOAD_URL}"
if [ $? -ne 0 ]; then
    echo "Failed to download ${ARCHIVE_NAME}."
    exit 1
fi

# Extract the tarball
echo "Extracting ${ARCHIVE_NAME}..."
tar -xzvf "${ARCHIVE_NAME}"
if [ $? -ne 0 ]; then
    echo "Failed to extract ${ARCHIVE_NAME}."
    exit 1
fi

# Move the binary to /usr/local/bin
echo "Moving ${BINARY_NAME} to /usr/local/bin/..."
sudo mv "${BINARY_NAME}" /usr/local/bin/pixelize
if [ $? -ne 0 ]; then
    echo "Failed to move ${BINARY_NAME} to /usr/local/bin/."
    exit 1
fi

# Make the binary executable
echo "Making ${BINARY_NAME} executable..."
sudo chmod +x /usr/local/bin/"${BINARY_NAME}"
if [ $? -ne 0 ]; then
    echo "Failed to make ${BINARY_NAME} executable."
    exit 1
fi

# Cleanup
echo "Cleaning up..."
rm "${ARCHIVE_NAME}"
if [ $? -ne 0 ]; then
    echo "Failed to remove ${ARCHIVE_NAME}."
    exit 1
fi

echo "Installation complete. You can start the app by typing '${APP_NAME}' in the terminal."
