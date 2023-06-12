FROM alpine:latest

RUN apk --update add --no-cache openssh bash
RUN sed -i s/#PermitRootLogin.*/PermitRootLogin\ yes/ /etc/ssh/sshd_config
RUN sed -i 's/#AllowTcpForwarding.*/AllowTcpForwarding yes/' /etc/ssh/sshd_config
RUN sed -ie 's/#Port 22/Port 22/g' /etc/ssh/sshd_config
RUN echo "root:password" | chpasswd
RUN rm -rf /var/cache/apk/*
RUN /usr/bin/ssh-keygen -A
RUN ssh-keygen -t rsa -b 4096 -f  /etc/ssh/ssh_host_key
ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]