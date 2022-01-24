FROM ubuntu:latest

ENV TZ=America/New_York
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update && apt-get install --yes --no-install-recommends \
  curl \
  wget \
  git \
  vim \
  cmake \
  python3 \
  libssl-dev \
  libcurl4-openssl-dev \
  libxml2-dev \
  libpng-dev \
  build-essential \
  nfs-kernel-server \
  && apt-get autoclean \
  && curl -sL https://deb.nodesource.com/setup_14.x  | bash - \
  && apt-get install -y python-setuptools \
  default-jre \
  && apt-get install -y alien unzip wget libz-dev \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

RUN curl https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh \
  > Miniconda3-latest-Linux-x86_64.sh \
  && yes \
  | bash Miniconda3-latest-Linux-x86_64.sh -b
RUN rm Miniconda3-latest-Linux-x86_64.sh

# put system path first so that conda doesn't override python
ENV PATH=$PATH:/root/miniconda3/bin/

# install "report" environment's dependencies
RUN conda update -n base -c defaults conda

WORKDIR /home
RUN git clone https://github.com/edotau/goFish.git

RUN bash goFish/.github/workflows/install.sh

