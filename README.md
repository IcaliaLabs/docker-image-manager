# Docker Image Manager

A CLI tool used to download and push docker images within a CI
pipeline.

It can be used to download images from a "cache_from" key in a compose file,
generate image tags based on Git properties - such as commit id, and branch name.

## vs a Shell file

...I got tired of different environments, where sometimes /bin/sh is replaced
with dash, and nothing ever worked the same. THat's why I decided to turn my old
shell file into an executable!

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