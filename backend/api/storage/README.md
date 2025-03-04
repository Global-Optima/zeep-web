# Storage System Documentation

## Overview

The application leverages an S3-compatible storage service for managing files like profile images, product images, and product videos. This system provides distinct storage paths, extensions, and maximum allowed sizes for each file type, ensuring organized and secure file handling.

## File Types and Specifications

Files' object keys are build without explicit format assigning, objects rely on content-type based on the original file extension

| Volume             | Storage Path       | Content Type        | Max Size |
|--------------------|--------------------|---------------------|----------|
| `converted-images` | `images/converted` | `image/webp`        | 5 MB     |
| `original-images`  | `images/original`  | `application/x-tar` | 5 MB     |
| `converted-videos` | `videos/products`  | `video/mp4`         | 20 MB    |

## Core Utilities

1. **`GetFileFromContext`**: Extracts the file from the request, identifies its type based on the `FileTypeMapping`, validates its size, and standardizes its name.
2. **`StandardizeFileName`**: Prevents duplicate extensions by applying the correct extension only once.
3. **`DetermineFileType`**: Maps a file’s extension to the corresponding `FileType`.

## Storage Handlers and Endpoints

### Endpoints

1. **Upload File**: `POST /api/v1/storage/`
   - Accepts a file (specified by volume, such as `profile`, `product-image`, or `product-video`) via a form field.
   - Checks for duplicate files, validates size, and uploads if validation passes.
   - **Response**: `{"filePath": "file/path/on/storage"}` or error message if validation fails.

   ```bash
   curl -X POST -F "profile=@path/to/profile.png" http://localhost:8080/api/v1/storage/
   ```

2. **Get File URL**: `GET /api/v1/storage/?filename=<file_name>`
   - Returns a publicly accessible URL for the requested file.
   - **Response**: `{"fileURL": "https://s3-endpoint/path/to/file"}`

   ```bash
   curl -X GET "http://localhost:8080/api/v1/storage/?filename=profile-picture.png"
   ```

3. **Delete File**: `DELETE /api/v1/storage/?filename=<file_name>`
   - Deletes the specified file if it exists in storage.
   - **Response**: `{"message": "File deleted successfully"}` or error message if the file is not found.

   ```bash
   curl -X DELETE "http://localhost:8080/api/v1/storage/?filename=profile-picture.png"
   ```

4. **Download and Save Locally** (Temp Endpoint): `GET /api/v1/storage/download?filename=<file_name>`
   - Downloads and saves a file to local storage (for testing).
   - **Response**: `{"message": "File downloaded and saved locally", "path": "<local_path>"}`

5. **List Buckets** (Temp Endpoint): `GET /api/v1/storage/list-buckets`
   - Lists available buckets in the S3-compatible storage service.
   - **Response**: `{"buckets": [{"Name": "bucket_name", "CreatedOn": "creation_date"}]}`

### Example Usage

1. **Uploading a Profile Picture**:
   ```bash
   curl -X POST -F "profile=@path/to/profile.png" http://localhost:8080/api/v1/storage/
   ```

2. **Retrieving a File URL**:
   ```bash
   curl -X GET "http://localhost:8080/api/v1/storage/?filename=profile.png"
   ```

3. **Deleting a File**:
   ```bash
   curl -X DELETE "http://localhost:8080/api/v1/storage/?filename=profile.png"
   ```

## Environment Configuration

To connect to the S3-compatible storage, ensure the following environment variables are correctly set:

- `S3_ACCESS_KEY`: Access key for storage(username for minio).
- `S3_SECRET_KEY`: Secret key for storage(password for minio).
- `S3_ENDPOINT`: Endpoint URL.
- `S3_BUCKET_NAME`: The bucket name for file storage.

### Important Notes

- **Size Validation**: Each file’s size is validated against its `MaxSize` to ensure compliance.
- **Error Handling**: Detailed error responses are provided for common issues like size validation failures, missing files, or duplicate uploads.

This storage configuration provides a structured, secure, and scalable file management solution within the application.
