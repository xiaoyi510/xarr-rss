FROM --platform=$TARGETPLATFORM  alpine:3.15
MAINTAINER Baozi<admin@52nyg.com>
# 使用变量必须申明
ARG TARGETOS
ARG TARGETARCH

# 时区
ENV TZ=Asia/Shanghai


#5.更新Alpine的软件源为阿里云，因为从默认官源拉取实在太慢了
RUN echo https://mirrors.aliyun.com/alpine/v3.15/main/ > /etc/apk/repositories && \
    echo https://mirrors.aliyun.com/alpine/v3.15/community/ >> /etc/apk/repositories
RUN apk update && apk upgrade
RUN apk add  wget bash && \
    apk add --update tzdata
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo "$TZ" > /etc/timezone \
    && rm -rf /var/cache/apk/*


###############################################################################
#                                INSTALLATION
###############################################################################
# 定义APP运行目录
ENV WORKDIR  /app

# 导入文件
ADD ./release/xarr-rss-${TARGETOS}-${TARGETARCH}/xarr-rss $WORKDIR/XArr-Rss
# 导入web
ADD ./web $WORKDIR/web

# 创建默认路径
RUN mkdir $WORKDIR/conf/ &&  mkdir $WORKDIR/conf/cache &&  mkdir $WORKDIR/conf/trans && chmod +x $WORKDIR/XArr-Rss
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2


###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./XArr-Rss

EXPOSE 8086