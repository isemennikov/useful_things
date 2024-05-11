# Use an official Python runtime as a parent image
FROM python:3.12-slim

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Make script executable
RUN chmod +x renamefile.py


# Run the script when the container launches
CMD ["python", "./renamefile.py"]

