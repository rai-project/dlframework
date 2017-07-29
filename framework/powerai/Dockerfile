FROM nvidia/cuda-ppc64le:8.0-cudnn6-devel
MAINTAINER Carl Pearson <pearson@illinois.edu>

RUN apt-get update \
    && apt-get install -y \
    curl \
    apt-transport-https \
    && apt-get upgrade -y \
    && rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://public.dhe.ibm.com/software/server/POWER/Linux/mldl/ubuntu/mldl-repo-network_3.4.0_ppc64el.deb -O

RUN dpkg -i mldl-repo-network_3.4.0_ppc64el.deb
RUN apt-get update
