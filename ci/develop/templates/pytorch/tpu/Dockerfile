FROM gcr.io/djw-ai/services/neurafuse/pytorch/tpu/base:latest

WORKDIR /app
COPY neurakube-linux-amd64 neurakube-linux-amd64
RUN mkdir lightning
CMD ["./neurakube-linux-amd64", "container", "lightning-py", "gpt"]