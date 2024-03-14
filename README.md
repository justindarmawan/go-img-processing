# Image Processing

Welcome to my project! This repository contains the source code and documentation for [Image Processing].

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine.

### Prerequisites

- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Postman: [Install Postman](https://www.postman.com/downloads/)

### Installing

1. Clone the repository:

   ```bash
   git clone git@github.com:justindarmawan/go-img-processing.git
   ```

2. Open a terminal and navigate to the directory where the repository is cloned.

3. Build the Docker Image:

   ```bash
   docker build -t go-img-processing-justin-image .
   ```

4. Run the Docker Container:

   ```bash
   docker run -p 8080:8080 go-img-processing-justin-image
   ```

## Using the Application

### Postman Setup and Sending Request Using Postman

1. Open Postman.

2. Import the Postman collection from the 'postman' directory in this repository.

3. There will be 3 requests: 'convert', 'resize', 'compress'.

4. For each request, navigate to the Body tab and choose form-data.

5. Fill in the required form fields as described in API spec below.

6. After filling in the required form fields:
   - Hit the Send button to send the request.
   - To download the image, use the Send and Download button available in Postman. The image will be downloaded to your local machine.

### API Spec

#### Convert

- **Method:** POST
- **Endpoint:** http://localhost:8080/convert/
- **Form-data:**
  - Key: image
    - Type: File
    - Description: The image file to convert (".png").

#### Resize

- **Method:** POST
- **Endpoint:** http://localhost:8080/resize/
- **Form-data:**
  - Key: image
    - Type: File
    - Description: The image file to resize (".png", ".jpg").
  - Key: width
    - Type: Number
    - Description: The desired width of the resized image, must be greater than 0.
    - Example: 300
  - Key: height
    - Type: Number
    - Description: The desired height of the resized image.
    - Example: 300
  - Key: lockAspectRatio
    - Type: Boolean
    - Description: The flag for lockAspectRatio.
    - Example: true
  - Key: bgColorRGB
    - Type: String
    - Description: Background color to fill if 'lockAspectRatio' is enabled, value from 0-255.
    - Example: 255,255,255

#### Compress

- **Method:** POST
- **Endpoint:** http://localhost:8080/compress/
- **Form-data:**
  - Key: image
    - Type: File
    - Description: The image file to resize (".png", ".jpg").
  - Key: quality
    - Type: Number
    - Description: The desired quality of the compressed image, value from 0-100.
    - Example: 80

## Test Application

### Run Test Cases

1. Open a terminal and navigate to the directory where the repository is cloned.

2. Build the Docker Image:

   ```bash
   docker build -t go-img-processing-justin-image .
   ```

3. Run the Docker container with the command to execute the tests:

   ```bash
   docker run go-img-processing-justin-image ./app test
   ```
