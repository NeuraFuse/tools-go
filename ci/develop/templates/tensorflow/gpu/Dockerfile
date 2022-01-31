FROM tensorflow/tensorflow:2.2.1-gpu-py3

WORKDIR /app

COPY neurakube-linux-amd64 neurakube-linux-amd64
RUN mkdir lightning

RUN pip install debugpy

CMD ["./neurakube-linux-amd64", "container", "tensorflow"]