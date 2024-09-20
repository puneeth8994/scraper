package utils_test

import (
	"os"
	"scraper/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadURLsFromFile_Success(t *testing.T) {
	// Create a temporary file with some test URLs
	tempFile, err := os.CreateTemp("", "urls_test.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Clean up

	// Write test URLs to the temp file
	_, err = tempFile.WriteString("http://example.com\nhttp://test.com\n")
	assert.NoError(t, err)

	// Close the file so it can be opened by the function
	tempFile.Close()

	// Call the function
	urls, err := utils.ReadURLsFromFile(tempFile.Name())
	assert.NoError(t, err)

	// Verify the URLs were read correctly
	expectedUrls := []string{
		"http://example.com",
		"http://test.com",
	}
	assert.Equal(t, expectedUrls, urls)
}

func TestReadURLsFromFile_FileNotFound(t *testing.T) {
	// Call the function with a non-existent file
	_, err := utils.ReadURLsFromFile("non_existent_file.txt")

	// Verify the error is returned
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestReadURLsFromFile_EmptyFile(t *testing.T) {
	// Create an empty temporary file
	tempFile, err := os.CreateTemp("", "empty_test.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Clean up

	// Close the file
	tempFile.Close()

	// Call the function
	urls, err := utils.ReadURLsFromFile(tempFile.Name())
	assert.NoError(t, err)

	// Verify that the result is an empty slice
	assert.Empty(t, urls)
}
