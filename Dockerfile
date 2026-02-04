FROM golang:1.22-alpine

WORKDIR /app
COPY main.go .

# Initialize module
RUN go mod init aion_ctf

# Build
RUN go build -o aion_server .

# Create the flag
RUN echo "CTF{r4c3_c0nd1t10ns_1n_g0_4r3_d34dly}" > /root/flag.txt
RUN chmod 400 /root/flag.txt

# Run as root (required to access /root, but typically CTFs run as non-root. 
# For this challenge, we assume the service has privileges or the flag is readable by the user).
# Let's make it standard:
RUN adduser -D ctf
USER ctf
COPY --chown=ctf:ctf main.go .
# Move flag to where ctf user can read it for this specific challenge mechanics
RUN echo "CTF{r4c3_c0nd1t10ns_1n_g0_4r3_d34dly}" > /app/flag.txt

EXPOSE 8080
CMD ["./aion_server"]