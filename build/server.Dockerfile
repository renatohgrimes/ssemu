FROM debian:bookworm-slim
RUN apt upgrade && apt update -y
ARG HOST_UID
RUN groupadd -r ssemu && useradd --system --gid ssemu --uid $HOST_UID ssemu
USER ssemu
WORKDIR /ssemu/bin
CMD [ "./ssemu" ]