# Watermark Generator

A Go-based web service for applying watermarks to images, with a React frontend.

## Features

- Apply image watermarks to uploaded images
- Supports JPEG and PNG formats
- Resize images
- Apply text watermarks
- Adjust watermark opacity
- Tile watermarks across the image
- React-based frontend

## Prerequisites

- Go 1.21 or later
- Node.js and npm (for building the React app)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/watermark-generator.git
   cd watermark-generator
   ```

2. Install Go dependencies:
   ```
   go mod download
   ```

3. Install React dependencies and build the frontend:
   ```
   cd frontend
   npm install
   npm run build
   cd ..
   ```

## Building and Running the Server

1. Build the Go server:
   ```
   go build -o watermark-server
   ```

2. Run the server:
   ```
   ./watermark-server
   ```

The server will start on `http://localhost:8080`, serving both the API and the React frontend.

## Development

For development, you can run the Go server and React app separately:

1. Run the Go server:
   ```
   go run main.go
   ```

2. In a separate terminal, run the React development server:
   ```
   cd frontend
   npm start
   ```

The React dev server will run on `http://localhost:3000` and proxy API requests to the Go server.

## API Endpoints

### POST /api/watermark

Apply a watermark to an image.

**Request:**
- Method: POST
- Content-Type: multipart/form-data
- Body:
  - `image`: The image file to watermark
  - `text`: Text to use as watermark
  - `textColor`: Color of the text watermark (hex format, e.g., "#FFFFFF")
  - `opacity`: Opacity of the watermark (0.0 to 1.0)

**Response:**
- Content-Type: image/png
- Body: The watermarked image

## Code Structure

- `main.go`: Entry point of the application
- `watermark/service.go`: Core watermarking functionality
- `api/handler.go`: HTTP request handling
- `frontend/`: React frontend application

## Dependencies

### Go Dependencies
- github.com/nfnt/resize
- golang.org/x/image
- github.com/disintegration/imaging
- github.com/golang/freetype
- github.com/lucasb-eyer/go-colorful

### Frontend Dependencies
See `frontend/package.json` for a full list of React dependencies.

## License

[MIT License](LICENSE)