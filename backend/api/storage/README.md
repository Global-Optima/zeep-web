# Storage System Documentation

## Overview

The application leverages an S3-compatible storage service for managing files like profile images, product images, and product videos. This system provides distinct storage paths, extensions, and maximum allowed sizes for each file type, ensuring organized and secure file handling.

## File Types and Specifications

Each file type is defined in the `FileType` struct, specifying the file path, extension, and maximum file size. The current configuration is as follows:

| Volume          | Storage Path       | Extension | Max Size   |
|-----------------|--------------------|-----------|------------|
| `profile`       | `images/profile`   | `.png`    | 2 MB       |
| `product-image` | `images/products`  | `.jpg`    | 5 MB       |
| `product-video` | `videos/products`  | `.mp4`    | 20 MB      |

### `FileType` Struct

- `Path`: Defines the directory path where files of a specific type are stored.
- `Extension`: Ensures consistency in file extensions for each file type.
- `MaxSize`: Limits file size for each type, validated during the upload process.

### `FileTypeMapping` Example

```go
var FileTypeMapping = map[string]FileType{
    "profile":       {Path: "images/profile", Extension: ".png", MaxSize: 2 * 1024 * 1024},
    "product-image": {Path: "images/products", Extension: ".jpg", MaxSize: 5 * 1024 * 1024},
    "product-video": {Path: "videos/products", Extension: ".mp4", MaxSize: 20 * 1024 * 1024},
}
```

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

- `PSKZ_ACCESS_KEY`: Access key for storage.
- `PSKZ_SECRET_KEY`: Secret key for storage.
- `PSKZ_ENDPOINT`: Endpoint URL.
- `PSKZ_BUCKETNAME`: The bucket name for file storage.

### Important Notes

- **Duplicate Checks**: Before uploading, the handler verifies if the file already exists using `FileExists`.
- **Size Validation**: Each file’s size is validated against its `MaxSize` to ensure compliance.
- **Error Handling**: Detailed error responses are provided for common issues like size validation failures, missing files, or duplicate uploads.

This storage configuration provides a structured, secure, and scalable file management solution within the application.
