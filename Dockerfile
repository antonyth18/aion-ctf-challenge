FROM golang:1.22-alpine

WORKDIR /app

# 1. Copy the code
COPY main.go .

# 2. Build the application (As Root)
RUN go mod init aion_ctf && go build -o aion_server .

# 3. Create the Flag (As Root - before switching users!)
RUN echo "CTF{r4c3_c0nd1t10ns_1n_g0_4r3_d34dly}" > /flag.txt

# 4. Create the non-root user
RUN adduser -D ctf

# 5. Switch to the non-root user for security
USER ctf

EXPOSE 8080
CMD ["./aion_server"]