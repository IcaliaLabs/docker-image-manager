# Docker Image Manager

A CLI tool used to download and push docker images within a CI
pipeline.

It can be used to download images from a "cache_from" key in a compose file,
generate image tags based on Git properties - such as commit id, and branch name.


## Extracting the binary from an image

```bash
# 1: Pull the image:
docker pull icalialabs/docker-image-manager:latest

# 1: create a container - will pull the image:
docker create --name image-manager icalialabs/docker-image-manager:latest

# 2: Extract the executable from the container:
docker cp image-manager:/docker-image-manager /usr/local/bin/

# 3: Remove the container:
docker rm image-manager
```