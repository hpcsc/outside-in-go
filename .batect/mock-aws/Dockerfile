FROM motoserver/moto:3.1.3

HEALTHCHECK CMD curl --fail --connect-timeout 5 --max-time 60 http://127.0.0.1:5000/moto-api/data.json
RUN apt-get update && \
    apt-get install -y curl unzip jq && \
    rm -rf /var/lib/apt/lists/*

RUN curl https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip -o /tmp/awscliv2.zip && \
    unzip /tmp/awscliv2.zip -d /tmp && \
    /tmp/aws/install && \
    rm -rf /tmp/awscliv2.zip /tmp/aws

ENV AWS_PAGER=""
ENV INITIAL_NO_AUTH_ACTION_COUNT=3

COPY setup.sh /scripts/setup.sh
