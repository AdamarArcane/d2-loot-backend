# Use a minimal base image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the pre-compiled binary into the container
COPY d2-loot-backend .

# Ensure the binary is executable
RUN chmod +x d2-loot-backend

# Expose the port your application listens on (adjust if necessary)
EXPOSE 8080

# Command to run the application
CMD ["./d2-loot-backend"]