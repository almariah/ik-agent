FROM golang:1.8
ADD ik-agent.yaml /opt/ik-agent/etc/
ADD .build/ik-mysql.so /opt/ik-agent/lib/
ADD .build/ik-agent /opt/ik-agent/bin/
CMD ["/opt/ik-agent/bin/ik-agent"]
