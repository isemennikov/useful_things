# syntax=docker/dockerfile:1.7.1
FROM python:3.11-slim as builder

# set vars
ARG USERNAME=hermes
ARG GNAME=gods
ARG APP_DIR=/home/${USERNAME}/app/
ARG REQUIREMENTS_FILE=requirements.txt
ARG REPO="https://github.com/EnableSecurity/wafw00f.git"
# create user 
RUN groupadd -r ${GNAME} && useradd -mr -g ${GNAME} ${USERNAME}

# updates install 
RUN set -eux \
    && apt update \
    && apt install -y --no-install-recommends \
        build-essential \
        git \
#    && pip install --upgrade pip setuptools \
    && rm -rf /var/lib/apt/lists/*

# app install
RUN git clone ${REPO} ${APP_DIR}
WORKDIR ${APP_DIR}
RUN python setup.py install
#RUN pip install --no-cache-dir -r ${APP_DIR}/${REQUIREMENTS_FILE}

#set credts
RUN chown -R ${USERNAME}:${GNAME} ${APP_DIR}

# create cleanup image 
FROM python:3.11-slim

# set vars 
ARG USERNAME=hermes
ARG GNAME=gods
ARG APP_DIR=/home/${USERNAME}/app/

# rootless 
RUN groupadd -r ${GNAME} && useradd -mr -g ${GNAME} ${USERNAME}
RUN set -eux \
    && apt-get update \
    && apt-get install -y --no-install-recommends openssl \
    && rm -rf /var/lib/apt/lists/*

# copy layrs from  first stage and set depents for second stage
COPY --from=builder ${APP_DIR} ${APP_DIR}
COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

# set creds and switch user
RUN chown -R ${USERNAME}:${GNAME} ${APP_DIR}
USER ${USERNAME}

# run app
WORKDIR ${APP_DIR}
ENTRYPOINT ["wafw00f"]
CMD ["-l"]
