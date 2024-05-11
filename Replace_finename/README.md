This script searches for a PDF and its file in a given directory, removes unkind words in the file name, and removes unnecessary characters. 
The script also limits the length of the file name to 60 characters.
As a result of the work, a registry.txt file is created. 

The script runs in a container  
For running need build docker image from Dockerfile  *docker build -t name:tag .*


Command for running *docker run -ti --rm -v $(pwd):/files renamefiles:1 /files*
I recommend using the Go version. Python version is not very good and works slowly